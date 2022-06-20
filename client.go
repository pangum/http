package http

import (
	"github.com/go-resty/resty/v2"
)

// Client 客户端封装
type Client struct {
	*resty.Client
}

func (c *Client) ResponseFields(rsp *resty.Response) *responseFields {
	return &responseFields{
		response: rsp,
	}
}
