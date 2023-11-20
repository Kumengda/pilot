package proxy

import (
	"github.com/Kumengda/pilot/filter"
	"github.com/Kumengda/pilot/options"
	"github.com/elazarl/goproxy"
	"io"
	"net/http"
	"net/url"
)

type Pilot struct {
	reqHandler      func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response, error)
	reqErrorHandler func(req *http.Request, ctx *goproxy.ProxyCtx, err error) (*http.Request, *http.Response)
	transport       *http.Transport
	proxy           *url.URL
	filter          []*options.Filter
}

func NewPilot(option ...options.Option) *Pilot {
	var pilot Pilot
	for _, o := range option {
		switch o.(type) {
		case *options.TransportOpt:
			pilot.transport = o.GetValue().(*http.Transport)
		case *options.Proxy:
			pilot.proxy = o.GetValue().(*url.URL)
		case *options.Filter:
			pilot.filter = append(pilot.filter, o.GetValue().(*options.Filter))
		}
	}
	return &pilot

}

func (p *Pilot) SetReqHandler(reqHandler func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response, error)) {
	p.reqHandler = reqHandler
}
func (p *Pilot) SetReqErrHandler(reqErrHandler func(req *http.Request, ctx *goproxy.ProxyCtx, err error) (*http.Request, *http.Response)) {
	p.reqErrorHandler = reqErrHandler
}

func (p *Pilot) Listen(addr string, verbose bool) error {
	proxy := goproxy.NewProxyHttpServer()
	if p.transport != nil {
		proxy.Tr = p.transport
	} else {
		proxy.Tr.Proxy = http.ProxyURL(p.proxy)
	}
	var conds []goproxy.ReqCondition
	for _, pf := range p.filter {
		for _, t := range pf.FilterType {
			switch t {
			case filter.Host:
				conds = append(conds, goproxy.ReqHostMatches(pf.Filters...))
			case filter.URL:
				for _, f := range pf.Filters {
					conds = append(conds, goproxy.UrlMatches(f))
				}
			}
		}
	}
	if conds == nil {
		conds = append(conds, goproxy.ReqHostMatches(filter.DropAll))
	}
	proxy.OnRequest(conds...).HandleConnect(goproxy.AlwaysMitm)

	proxy.OnRequest(conds...).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if !p.doFilter(req) {
			return req, nil
		}

		if p.reqHandler == nil {
			return req, nil
		}
		req, resp, err := p.reqHandler(req, ctx)
		if err != nil {
			if p.reqErrorHandler != nil {
				return p.reqErrorHandler(req, ctx, err)
			} else {
				return req, nil
			}
		}
		return req, resp
	})
	proxy.Verbose = verbose
	err := http.ListenAndServe(addr, proxy)
	if err != nil {
		return err
	}
	return nil
}
func (p *Pilot) doFilter(req *http.Request) bool {
	for _, fp := range p.filter {
		for _, t := range fp.FilterType {
			switch t {
			case filter.BODY:
				readRes, err := io.ReadAll(req.Body)
				if err != nil {
					return false
				}
				for _, v := range fp.Filters {
					match := v.MatchString(string(readRes))
					if match {
						return true
					}
				}
			case filter.HEADER:
				for k, v := range req.Header {
					for _, f := range fp.Filters {
						match := f.MatchString(k)
						if match {
							return true
						}
						var keyVal string
						for _, s := range v {
							keyVal = keyVal + s
						}
						match = f.MatchString(keyVal)
						if match {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
