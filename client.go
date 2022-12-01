package http

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

// Client 客户端封装
type Client struct {
	*resty.Client

	proxies map[string]*proxy
}

func _newClient(client *resty.Client, proxies []*proxy) (_client *Client) {
	_proxies := make(map[string]*proxy)
	for _, _proxy := range proxies {
		_proxies[_proxy.Target] = _proxy
	}

	_client = &Client{
		Client: client,

		proxies: _proxies,
	}
	// 设置动态代理
	client.OnBeforeRequest(_client.beforeRequest)
	client.OnAfterResponse(_client.afterResponse)

	return
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

func (c *Client) beforeRequest(client *resty.Client, req *resty.Request) (err error) {
	if host, he := c.host(req.URL); nil != he {
		err = he
	} else if _proxy, hostOk := c.proxies[host]; hostOk {
		client.SetProxy(_proxy.addr())
	} else if _proxy, defaultOk := c.proxies[targetDefault]; defaultOk {
		client.SetProxy(_proxy.addr())
	}

	return
}

func (c *Client) afterResponse(client *resty.Client, _ *resty.Response) (err error) {
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
