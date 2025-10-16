package dns

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/elect0/dns-blocker/internal/logging"
	"github.com/miekg/dns"
)

type Handler struct {
	logger         *logging.Logger
	blocklist      map[string]struct{}
	customRecords  map[string]net.IP
	upstreamServer string

	cache      map[string]*dns.Msg
	cacheMutex sync.RWMutex
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
		cache:          make(map[string]*dns.Msg),
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

	h.logger.Debug("forwarding query to upstream server", logProps)

	response, err := forwardDoH(r, h.upstreamServer)
	if err != nil {
		h.logger.Error("failed to forward query to doh upstream server", err, logProps)
		msg.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(msg)
		return
	}

	w.WriteMsg(response)
}
