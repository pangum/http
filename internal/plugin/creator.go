package plugin

import (
	"github.com/goexl/http"
	"github.com/goexl/simaqian"
	"github.com/pangum/http/internal/core"
	"github.com/pangum/pangu"
)

type Creator struct {
	// 用于提供构造方法
}

func (c *Creator) New(loader *pangu.Config, logger simaqian.Logger) (client *http.Client, err error) {
	wrapper := new(core.Wrapper)
	if le := loader.Load(wrapper); nil != le {
		err = le
	} else {
		client = c.new(&wrapper.Http.Client, logger)
	}

	return
}

func (c *Creator) new(config *core.Config, logger simaqian.Logger) *http.Client {
	builder := http.New().Payload(*config.Payload).Timeout(config.Timeout).Logger(logger)
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
