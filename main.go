package main

import (
	"fmt"
	"github.com/Kumengda/pilot/filter"
	"github.com/Kumengda/pilot/options"
	"github.com/Kumengda/pilot/proxy"
	"github.com/elazarl/goproxy"
	"net/http"
)

func main() {

	pilot := proxy.NewPilot([]options.Option{
		//options.SetCheckProxy("http://127.0.0.1:8082"),
		options.SetFilter(filter.NewFilterWithPrefix("www.bilibili.com", "www.baidu.com"), filter.Host),
		options.SetFilter(filter.NewFilterAllowAll(), filter.BODY),
		//options.SetFilter(filter.NewFilterWithPrefix("www.bilibili.com"), filter.Host),
		//options.SetFilter(filter.NewFilterAllowAll(), filter.HEADER),
	}...)

	pilot.SetReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response, error) {

		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusForbidden,
			"hhhhhhhhhhhhhhhh"), nil
	})

	err := pilot.Listen(":8080", true)

	if err != nil {
		fmt.Println(err)
	}
}
