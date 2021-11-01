package http

import (
	`crypto/tls`

	`github.com/go-resty/resty/v2`
	`github.com/pangum/pangu`
	`github.com/storezhang/gox`
)

func newClient(config *pangu.Config) (client *Client, err error) {
	_panguConfig := new(panguConfig)
	if err = config.Load(_panguConfig); nil != err {
		return
	}
	clientConfig := _panguConfig.Http.Client

	restyClient := resty.New()
	if `` != clientConfig.Proxy.Host {
		restyClient.SetProxy(clientConfig.Proxy.Addr())
	}
	if 0 != clientConfig.Timeout {
		restyClient.SetTimeout(clientConfig.Timeout)
	}
	if clientConfig.Payload.Get {
		restyClient.SetAllowGetMethodPayload(true)
	}
	if clientConfig.Certificate.Skip {
		// nolint:gosec
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	} else {
		if `` != clientConfig.Certificate.Root {
			restyClient.SetRootCertificate(clientConfig.Certificate.Root)
		}
		if 0 != len(clientConfig.Certificate.Clients) {
			certificates := make([]tls.Certificate, 0, len(clientConfig.Certificate.Clients))
			for _, c := range clientConfig.Certificate.Clients {
				certificate, certificateErr := tls.LoadX509KeyPair(c.Public, c.Private)
				if nil != certificateErr {
					continue
				}
				certificates = append(certificates, certificate)
			}
			restyClient.SetCertificates(certificates...)
		}
	}
	if 0 != len(clientConfig.Headers) {
		restyClient.SetHeaders(clientConfig.Headers)
	}
	if 0 != len(clientConfig.Queries) {
		restyClient.SetQueryParams(clientConfig.Queries)
	}
	if 0 != len(clientConfig.Forms) {
		restyClient.SetFormData(clientConfig.Forms)
	}
	if 0 != len(clientConfig.Cookies) {
		restyClient.SetCookies(clientConfig.Cookies)
	}
	if `` != clientConfig.Auth.Type {
		switch clientConfig.Auth.Type {
		case gox.AuthTypeBasic:
			restyClient.SetBasicAuth(clientConfig.Auth.Username, clientConfig.Auth.Password)
		case gox.AuthTypeToken:
			restyClient.SetAuthToken(clientConfig.Auth.Token)
			if `` != clientConfig.Auth.Scheme {
				restyClient.SetAuthScheme(clientConfig.Auth.Scheme)
			}
		}
	}
	client = &Client{Client: restyClient}

	return
}

func newRequest(client *Client) *Request {
	return &Request{Request: client.R()}
}
