package dns

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/elect0/dns-blocker/internal/logging"
	"github.com/miekg/dns"
)

type Handler struct {
	logger         *logging.Logger
	blocklist      map[string]struct{}
	customRecords  map[string]net.IP
	upstreamServer string

	cache      map[string]*cacheEntry
	cacheMutex sync.RWMutex
}

type cacheEntry struct {
	msg    *dns.Msg
	expiry time.Time
}

func NewHandler(logger *logging.Logger, blocklist map[string]struct{}, customRecordsConfig map[string]string, upstreamServer string) (*Handler, error) {
	parsedRecords := make(map[string]net.IP)
	for domain, ipStr := range customRecordsConfig {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("invalid ip address '%s' for custom record '%s'", ipStr, domain)
		}

		if !strings.HasSuffix(domain, ".") {
			domain += "."
		}

		parsedRecords[domain] = ip
	}
	return &Handler{
		logger:         logger,
		blocklist:      blocklist,
		customRecords:  parsedRecords,
		upstreamServer: upstreamServer,
		cache:          make(map[string]*cacheEntry),
	}, nil
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
	question := r.Question[0]

	cleanDomain := strings.TrimSuffix(question.Name, ".")
	logProps := map[string]string{"domain": cleanDomain, "type": dns.TypeToString[question.Qtype]}

	if ip, ok := h.customRecords[question.Name]; ok {
		h.logger.Info("domain matched a custom record", logProps)

		if question.Qtype == dns.TypeA {
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", question.Name, ip.String()))
			if err != nil {
				msg.Answer = append(msg.Answer, rr)
			}
		}

		w.WriteMsg(msg)
		return
	}

	if _, ok := h.blocklist[cleanDomain]; ok {
		h.logger.Info("domain is on the blocklist", logProps)
		msg.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(msg)
		return
	}

	cacheKey := question.Name + dns.TypeToString[question.Qtype]
	h.cacheMutex.RLock()
	entry, found := h.cache[cacheKey]
	h.cacheMutex.RUnlock()

	if found && time.Now().Before(entry.expiry) {
		h.logger.Info("cache hit: serving response from cache", logProps)
		response := *entry.msg
		response.SetReply(r)
		w.WriteMsg(&response)
		return
	}

	h.logger.Debug("forwarding query to upstream server", logProps)

	response, err := forwardDoH(r, h.upstreamServer)
	if err != nil {
		h.logger.Error("failed to forward query to doh upstream server", err, logProps)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	if len(response.Answer) > 0 {
		ttl := response.Answer[0].Header().Header().Ttl
		for _, rr := range response.Answer {
			if rr.Header().Header().Ttl < ttl {
				ttl = rr.Header().Header().Ttl
			}
		}

		h.cacheMutex.Lock()
		h.cache[cacheKey] = &cacheEntry{
			msg:    response.Copy(),
			expiry: time.Now().Add(time.Duration(ttl) * time.Second),
		}
		h.cacheMutex.Unlock()

		propsWithTTL := map[string]string{"domain": cleanDomain, "ttl": fmt.Sprintf("%d", ttl)}
		h.logger.Info("response cached", propsWithTTL)
	}

	w.WriteMsg(response)
}
