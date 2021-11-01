package http

import (
	`github.com/go-resty/resty/v2`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
)

// Client 客户端封装
type Client struct {
	*resty.Client
}

func (c *Client) ResponseFields(rsp *resty.Response) (fields gox.Fields) {
	if nil == rsp {
		return
	}

	fields = []gox.Field{
		field.Int(`code`, rsp.StatusCode()),
		field.String(`body`, string(rsp.Body())),
	}

	return
}
