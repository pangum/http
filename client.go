package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/pangum/logging"
	"github.com/pangum/pangu"
)

// Client 客户端封装
type Client struct {
	*resty.Client

	logger  logging.Logger
	proxies map[string]*proxy
}

func newClient(config *pangu.Config, logger logging.Logger) (client *Client, err error) {
	_panguConfig := new(panguConfig)
	if err = config.Load(_panguConfig); nil != err {
		return
	}
	_config := _panguConfig.Http.Client

	client = new(Client)
	client.Client = resty.New()
	client.logger = logger
	client.proxies = make(map[string]*proxy)

	if 0 != _config.Timeout {
		client.SetTimeout(_config.Timeout)
	}
	if nil != _config.Payload {
		client.SetAllowGetMethodPayload(_config.Payload.Get)
	}
	if nil != _config.Certificate && *_config.Certificate.Enabled {
		if _config.Certificate.Skip {
			// nolint:gosec
			client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		} else {
			if "" != _config.Certificate.Root {
				client.SetRootCertificate(_config.Certificate.Root)
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
				client.SetCertificates(certificates...)
			}
		}
	}
	if 0 != len(_config.Headers) {
		client.SetHeaders(_config.Headers)
	}
	if 0 != len(_config.Queries) {
		client.SetQueryParams(_config.Queries)
	}
	if 0 != len(_config.Forms) {
		client.SetFormData(_config.Forms)
	}
	if 0 != len(_config.Cookies) {
		client.SetCookies(_config.Cookies)
	}
	if nil != _config.Auth && *_config.Auth.Enabled {
		switch _config.Auth.Type {
		case authTypeBasic:
			client.SetBasicAuth(_config.Auth.Username, _config.Auth.Password)
		case authTypeToken:
			client.SetAuthToken(_config.Auth.Token)
			if "" != _config.Auth.Scheme {
				client.SetAuthScheme(string(_config.Auth.Scheme))
			}
		}
	}

	// 设置动态代理
	if nil != _config.Proxy && *_config.Proxy.Enabled && "" == _config.Proxy.Target && 0 == len(_config.Proxies) {
		addr := _config.Proxy.addr()
		client.SetProxy(addr)
		logger.Debug("设置通用代理服务器", field.New("proxy", addr))
	} else {
		if nil != _config.Proxy && *_config.Proxy.Enabled {
			target := gox.Ift("" == _config.Proxy.Target, targetDefault, _config.Proxy.Target)
			client.proxies[target] = _config.Proxy
		}
		for _, _proxy := range _config.Proxies {
			if *_proxy.Enabled {
				client.proxies[_proxy.Target] = _proxy
			}
		}
	}
	// 动态代理
	if 0 != len(client.proxies) {
		client.OnBeforeRequest(client.setProxy)
		client.OnAfterResponse(client.unsetProxy)
	}
	// 记录日志
	client.SetPreRequestHook(client.log)

	return
}

func (c *Client) Curl(rsp *resty.Response) (string, error) {
	return c.curl(rsp.Request)
}

func (c *Client) Fields(rsp *resty.Response) (fields gox.Fields[any]) {
	if nil == rsp {
		return
	}

	fields = gox.Fields[any]{
		field.New("url", rsp.Request.URL),
		field.New("code", rsp.StatusCode()),
		field.New("body", string(rsp.Body())),
	}

	return
}

func (c *Client) log(_ *resty.Client, req *http.Request) (err error) {
	fields := gox.Fields[any]{
		field.New("url", req.URL),
	}
	if nil != req.Body {
		if body, re := io.ReadAll(req.Body); nil != re {
			err = re
		} else {
			req.Body = nopCloser{bytes.NewBuffer(body)}
			fields = append(fields, field.New[json.RawMessage]("body", body))
		}
	}
	if nil != err {
		return
	}

	for key, value := range req.Header {
		if 1 == len(value) {
			fields = append(fields, field.New(fmt.Sprintf("header.%s", key), value[0]))
		} else if 1 < len(value) {
			fields = append(fields, field.New(fmt.Sprintf("header.%s", key), value))
		}
	}
	c.logger.Debug("向服务器发送请求", fields...)

	return
}

func (c *Client) setProxy(client *resty.Client, req *resty.Request) (err error) {
	if host, he := c.host(req.URL); nil != he {
		err = he
	} else if hp, hostOk := c.proxies[host]; hostOk {
		addr := hp.addr()
		client.SetProxy(addr)
		c.logger.Debug("设置代理服务器", field.New("url", req.URL), field.New("proxy", addr))
	} else if dp, defaultOk := c.proxies[targetDefault]; defaultOk {
		addr := dp.addr()
		client.SetProxy(addr)
		c.logger.Debug("设置代理服务器", field.New("url", req.URL), field.New("proxy", addr))
	}

	return
}

func (c *Client) unsetProxy(client *resty.Client, _ *resty.Response) (err error) {
	client.RemoveProxy()

	return
}

func (c *Client) host(raw string) (host string, err error) {
	if _url, ue := url.Parse(raw); nil != ue {
		err = ue
	} else {
		host = _url.Host
	}

	return
}

func (c *Client) curl(req *resty.Request) (curl string, err error) {
	command := new(strings.Builder)
	command.WriteString("curl")
	command.WriteString("-X")
	command.WriteString(c.bashEscape(req.Method))

	if nil != req.Body {
		if body, re := io.ReadAll(req.RawRequest.Body); nil != re {
			err = re
		} else {
			req.Body = nopCloser{bytes.NewBuffer(body)}
			command.WriteString("-d")
			command.WriteString(c.bashEscape(string(body)))
		}
	}
	if nil != err {
		return
	}

	keys := make([]string, 0, len(req.Header))
	for key := range req.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		command.WriteString("-H")
		command.WriteString(c.bashEscape(fmt.Sprintf("%s: %s", key, strings.Join(req.Header[key], " "))))
	}
	command.WriteString(c.bashEscape(req.URL))

	return
}

func (c *Client) bashEscape(from string) string {
	return `'` + strings.Replace(from, `'`, `'\''`, -1) + `'`
}
