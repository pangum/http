package plugin

import (
	"github.com/goexl/http"
	"github.com/pangum/http/internal/core"
	"github.com/pangum/pangu"
)

type Creator struct {
	// 用于提供构造方法
}

func (c *Creator) New(config *pangu.Config) (client *http.Client, err error) {
	wrapper := new(core.Wrapper)
	if ge := config.Build().Get(wrapper); nil != ge {
		err = ge
	} else {
		client = c.new(&wrapper.Http.Client)
	}

	return
}

func (c *Creator) new(config *core.Config) *http.Client {
	builder := http.New().Payload(*config.Payload).Timeout(config.Timeout)
	builder = builder.Headers(config.Headers)
	builder = builder.Forms(config.Forms)
	builder = builder.Queries(config.Queries)

	if nil != config.Proxy {
		config.Proxies = append(config.Proxies, config.Proxy)
	}
	pb := builder.Proxy()
	for _, proxy := range config.Proxies {
		pb.Host(proxy.Host).Target(proxy.Target).Scheme(proxy.Scheme).Basic(proxy.Username, proxy.Password).Build()
	}

	if nil != config.Auth && config.Auth.Enable() {
		conf := config.Auth
		ab := builder.Auth()
		_ = ab.Scheme(conf.Scheme).Token(conf.Token).Basic(conf.Username, conf.Password).Build()
	}

	if nil != config.Certificate && config.Certificate.Enable() {
		conf := config.Certificate
		cb := builder.Certificate()
		cb = cb.Skip(*conf.Skip).Root(conf.Root)
		for _, client := range conf.Clients {
			cb = cb.Client(client.Public, client.Private)
		}
		_ = cb.Build()
	}

	return builder.Build()
}
