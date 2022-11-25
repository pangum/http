package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

// Client 客户端封装
type Client struct {
	*resty.Client
}

func (c *Client) Fields(rsp *resty.Response) (fields gox.Fields[any]) {
	if nil == rsp {
		return
	}

	fields = gox.Fields[any]{
		field.New("code", rsp.StatusCode()),
		field.New("body", string(rsp.Body())),
	}

	return
}
