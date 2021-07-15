package http

import (
	`github.com/go-resty/resty/v2`
)

// Request 请求封装
type Request struct {
	*resty.Request
}
