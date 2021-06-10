package http

import (
	`crypto/tls`

	`github.com/go-resty/resty/v2`
	`github.com/storezhang/gox`
	`github.com/storezhang/pangu`
)

func newClient(conf pangu.Config) (client *resty.Client, err error) {
	config := new(config)
	if err = conf.Struct(config); nil != err {
		return
	}
	client = newClientWithConfig(config.Http.Client)

	return
}

func newClientWithConfig(config ClientConfig) (client *resty.Client) {
	client = resty.New()
	if "" != config.Proxy.Host {
		client.SetProxy(config.Proxy.Addr())
	}
	if 0 != config.Timeout {
		client.SetTimeout(config.Timeout)
	}
	if config.AllowGetPayload {
		client.SetAllowGetMethodPayload(true)
	}
	if config.Certificate.Skip {
		// nolint:gosec
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		if "" != config.Certificate.Root {
			client.SetRootCertificate(config.Certificate.Root)
		}
		if 0 != len(config.Certificate.Clients) {
			certificates := make([]tls.Certificate, 0, len(config.Certificate.Clients))
			for _, c := range config.Certificate.Clients {
				certificate, err := tls.LoadX509KeyPair(c.Public, c.Private)
				if nil != err {
					continue
				}
				certificates = append(certificates, certificate)
			}
			client.SetCertificates(certificates...)
		}
	}
	if 0 != len(config.Headers) {
		client.SetHeaders(config.Headers)
	}
	if 0 != len(config.Queries) {
		client.SetQueryParams(config.Queries)
	}
	if 0 != len(config.Forms) {
		client.SetFormData(config.Forms)
	}
	if 0 != len(config.Cookies) {
		client.SetCookies(config.Cookies)
	}
	if "" != config.Auth.Type {
		switch config.Auth.Type {
		case gox.AuthTypeBasic:
			client.SetBasicAuth(config.Auth.Username, config.Auth.Password)
		case gox.AuthTypeToken:
			client.SetAuthToken(config.Auth.Token)
			if "" != config.Auth.Scheme {
				client.SetAuthScheme(config.Auth.Scheme)
			}
		}
	}

	return
}

func newRequest(client *resty.Client) *resty.Request {
	return client.R()
}
