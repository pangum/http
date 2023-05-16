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
	_       gox.CannotCopy
}

func newClient(config *pangu.Config, logger logging.Logger) (client *Client, err error) {
	wrap := new(wrapper)
	if err = config.Load(wrap); nil != err {
		return
	}

	conf := wrap.Http.Client
	client = new(Client)
	client.Client = resty.New()
	client.logger = logger
	client.proxies = make(map[string]*proxy)

	if 0 != conf.Timeout {
		client.SetTimeout(conf.Timeout)
	}
	if nil != conf.Payload {
		client.SetAllowGetMethodPayload(conf.Payload.Get)
	}
	if nil != conf.Certificate && *conf.Certificate.Enabled {
		if conf.Certificate.skipped() {
			// nolint:gosec
			client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		} else {
			if "" != conf.Certificate.Root {
				client.SetRootCertificate(conf.Certificate.Root)
			}
			if 0 != len(conf.Certificate.Clients) {
				certificates := make([]tls.Certificate, 0, len(conf.Certificate.Clients))
				for _, c := range conf.Certificate.Clients {
					cert, ce := tls.LoadX509KeyPair(c.Public, c.Private)
					if nil != ce {
						continue
					}
					certificates = append(certificates, cert)
				}
				client.SetCertificates(certificates...)
			}
		}
	}
	if 0 != len(conf.Headers) {
		client.SetHeaders(conf.Headers)
	}
	if 0 != len(conf.Queries) {
		client.SetQueryParams(conf.Queries)
	}
	if 0 != len(conf.Forms) {
		client.SetFormData(conf.Forms)
	}
	if 0 != len(conf.Cookies) {
		client.SetCookies(conf.Cookies)
	}
	if nil != conf.Auth && *conf.Auth.Enabled {
		switch conf.Auth.Type {
		case authTypeBasic:
			client.SetBasicAuth(conf.Auth.Username, conf.Auth.Password)
		case authTypeToken:
			client.SetAuthToken(conf.Auth.Token)
			if "" != conf.Auth.Scheme {
				client.SetAuthScheme(string(conf.Auth.Scheme))
			}
		}
	}

	// 设置动态代理
	if nil != conf.Proxy && *conf.Proxy.Enabled && "" == conf.Proxy.Target && 0 == len(conf.Proxies) {
		addr := conf.Proxy.addr()
		client.SetProxy(addr)
		logger.Debug("设置通用代理服务器", field.New("proxy", addr))
	} else {
		if nil != conf.Proxy && *conf.Proxy.Enabled {
			target := gox.Ift("" == conf.Proxy.Target, targetDefault, conf.Proxy.Target)
			client.proxies[target] = conf.Proxy
		}
		for _, _proxy := range conf.Proxies {
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
