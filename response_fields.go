package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type responseFields struct {
	response *resty.Response
}

func (rf *responseFields) Fields() (fields gox.Fields[any]) {
	if nil == rf.response {
		return
	}

	fields = gox.Fields[any]{
		field.New("code", rf.response.StatusCode()),
		field.New("body", string(rf.response.Body())),
	}

	return
}
