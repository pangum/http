package plugin

import (
	`github.com/go-resty/resty/v2`
)

// Request 请求封装
type Request struct {
	*resty.Request
}

func NewRequest(client *Client) *Request {
	return &Request{Request: client.R()}
}
