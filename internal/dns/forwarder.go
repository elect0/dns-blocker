package dns

import (
	"bytes"
	"io"
	"net/http"

	"github.com/miekg/dns"
)

func forwardDoH(r *dns.Msg, upstreamUrl string) (*dns.Msg, error) {
	packedMsg, err := r.Pack()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, upstreamUrl, bytes.NewReader(packedMsg))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/dns-message")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseMsg := new(dns.Msg)
	err = responseMsg.Unpack(body)
	if err != nil {
		return nil, err
	}

	return responseMsg, nil
}
