package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

// Client 客户端封装
type Client struct {
	*resty.Client

	proxy *proxy
	auth  *auth
}

func _newClient(client *resty.Client, proxy *proxy, auth *auth) *Client {
	return &Client{
		Client: client,

		proxy: proxy,
		auth:  auth,
	}
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
	if c.IsProxySet() {
		fields = append(fields, field.New("proxy", c.proxy.Addr()))
	}

	return
}
