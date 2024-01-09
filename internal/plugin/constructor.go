package plugin

import (
	"github.com/goexl/http"
	"github.com/pangum/pangu"
)

type Constructor struct {
	// 构造方法
}

func (c *Constructor) New(config *pangu.Config) (client *http.Client, err error) {
	wrapper := new(Wrapper)
	if ge := config.Build().Get(wrapper); nil != ge {
		err = ge
	} else {
		client = c.new(&wrapper.Http.Client)
	}

	return
}

func (c *Constructor) new(config *Config) *http.Client {
	builder := http.New().Payload(*config.Payload).Timeout(config.Timeout)
	builder = builder.Headers(config.Headers)
	builder = builder.Forms(config.Forms)
	builder = builder.Queries(config.Queries)

	if nil != config.Proxy {
		config.Proxies = append(config.Proxies, config.Proxy)
	}
	proxy := builder.Proxy()
	for _, conf := range config.Proxies {
		proxy.Host(conf.Host).Port(conf.Port).Target(conf.Target).Exclude(conf.Exclude).
			Scheme(conf.Scheme).
			Basic(conf.Username, conf.Password).
			Build()
	}

	if nil != config.Auth && config.Auth.Enable() {
		conf := config.Auth
		ab := builder.Auth()
		_ = ab.Scheme(conf.Scheme).Token(conf.Token).Basic(conf.Username, conf.Password).Build()
	}

	if config.Certificate.Enable() {
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
