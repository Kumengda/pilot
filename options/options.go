package options

import (
	"net/http"
	"net/url"
	"time"
)

type Option interface {
	GetValue() interface{}
}

type TransportOpt struct {
	transport *http.Transport
}

func (t *TransportOpt) GetValue() interface{} {
	if t.transport == nil {
		return &http.Transport{}
	}
	return t.transport
}

func SetTransportOpt(transport *http.Transport) *TransportOpt {
	return &TransportOpt{transport: transport}
}

type TimeoutOpt struct {
	timeout time.Duration
}

func (t *TimeoutOpt) GetValue() interface{} {
	return t.timeout
}

func SetTimeoutOpt(timeout int) *TimeoutOpt {
	return &TimeoutOpt{timeout: time.Duration(timeout) * time.Second}
}

type CheckRedirectOpt struct {
	checkFunc func(req *http.Request, via []*http.Request) error
}

func (c *CheckRedirectOpt) GetValue() interface{} {
	if c.checkFunc == nil {
		return func(req *http.Request, via []*http.Request) error {
			return nil
		}
	}
	return c.checkFunc
}
func SetCheckRedirectOpt(checkFunc func(req *http.Request, via []*http.Request) error) *CheckRedirectOpt {
	return &CheckRedirectOpt{checkFunc: checkFunc}
}

type Proxy struct {
	proxy string
}

func (p *Proxy) GetValue() interface{} {
	parse, err := url.Parse(p.proxy)
	if err != nil {
		return nil
	}
	return parse
}
func SetCheckProxy(proxy string) *Proxy {
	return &Proxy{proxy: proxy}
}
