package http

import (
	`crypto/tls`

	`github.com/go-resty/resty/v2`
	`github.com/storezhang/gox`
	`github.com/storezhang/pangu`
)

func newClient(config *pangu.Config) (restyClient *resty.Client, err error) {
	panguConfig := new(panguConfig)
	if err = config.Load(panguConfig); nil != err {
		return
	}
	client := panguConfig.Http.Client

	restyClient = resty.New()
	if "" != client.Proxy.Host {
		restyClient.SetProxy(client.Proxy.Addr())
	}
	if 0 != client.Timeout {
		restyClient.SetTimeout(client.Timeout)
	}
	if client.Payload.Get {
		restyClient.SetAllowGetMethodPayload(true)
	}
	if client.Certificate.Skip {
		// nolint:gosec
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		if "" != client.Certificate.Root {
			restyClient.SetRootCertificate(client.Certificate.Root)
		}
		if 0 != len(client.Certificate.Clients) {
			certificates := make([]tls.Certificate, 0, len(client.Certificate.Clients))
			for _, c := range client.Certificate.Clients {
				certificate, err := tls.LoadX509KeyPair(c.Public, c.Private)
				if nil != err {
					continue
				}
				certificates = append(certificates, certificate)
			}
			restyClient.SetCertificates(certificates...)
		}
	}
	if 0 != len(client.Headers) {
		restyClient.SetHeaders(client.Headers)
	}
	if 0 != len(client.Queries) {
		restyClient.SetQueryParams(client.Queries)
	}
	if 0 != len(client.Forms) {
		restyClient.SetFormData(client.Forms)
	}
	if 0 != len(client.Cookies) {
		restyClient.SetCookies(client.Cookies)
	}
	if "" != client.Auth.Type {
		switch client.Auth.Type {
		case gox.AuthTypeBasic:
			restyClient.SetBasicAuth(client.Auth.Username, client.Auth.Password)
		case gox.AuthTypeToken:
			restyClient.SetAuthToken(client.Auth.Token)
			if "" != client.Auth.Scheme {
				restyClient.SetAuthScheme(client.Auth.Scheme)
			}
		}
	}

	return
}

func newRequest(client *resty.Client) *resty.Request {
	return client.R()
}
