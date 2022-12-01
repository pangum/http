package http

import (
	"crypto/tls"

	"github.com/go-resty/resty/v2"
	"github.com/pangum/pangu"
)

func newClient(config *pangu.Config) (client *Client, err error) {
	_panguConfig := new(panguConfig)
	if err = config.Load(_panguConfig); nil != err {
		return
	}
	_config := _panguConfig.Http.Client

	proxies := make([]*proxy, 0)
	if nil != _config.Proxy {
		_config.Proxy.Target = targetDefault
		proxies = append(proxies, _config.Proxy)
	}
	proxies = append(proxies, _config.Proxies...)

	_client := resty.New()
	if 0 != _config.Timeout {
		_client.SetTimeout(_config.Timeout)
	}
	if nil != _config.Payload {
		_client.SetAllowGetMethodPayload(_config.Payload.Get)
	}
	if nil != _config.Certificate {
		if _config.Certificate.Skip {
			// nolint:gosec
			_client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		} else {
			if "" != _config.Certificate.Root {
				_client.SetRootCertificate(_config.Certificate.Root)
			}
			if 0 != len(_config.Certificate.Clients) {
				certificates := make([]tls.Certificate, 0, len(_config.Certificate.Clients))
				for _, c := range _config.Certificate.Clients {
					certificate, certificateErr := tls.LoadX509KeyPair(c.Public, c.Private)
					if nil != certificateErr {
						continue
					}
					certificates = append(certificates, certificate)
				}
				_client.SetCertificates(certificates...)
			}
		}
	}
	if 0 != len(_config.Headers) {
		_client.SetHeaders(_config.Headers)
	}
	if 0 != len(_config.Queries) {
		_client.SetQueryParams(_config.Queries)
	}
	if 0 != len(_config.Forms) {
		_client.SetFormData(_config.Forms)
	}
	if 0 != len(_config.Cookies) {
		_client.SetCookies(_config.Cookies)
	}
	if nil != _config.Auth {
		switch _config.Auth.Type {
		case authTypeBasic:
			_client.SetBasicAuth(_config.Auth.Username, _config.Auth.Password)
		case authTypeToken:
			_client.SetAuthToken(_config.Auth.Token)
			if "" != _config.Auth.Scheme {
				_client.SetAuthScheme(string(_config.Auth.Scheme))
			}
		}
	}
	client = _newClient(_client, proxies)

	return
}

func newRequest(client *Client) *Request {
	return &Request{Request: client.R()}
}
