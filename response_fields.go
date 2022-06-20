package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type responseFields struct {
	response *resty.Response
}

func (f *responseFields) Fields() (fields gox.Fields) {
	if nil == f.response {
		return
	}

	fields = []gox.Field{
		field.Int(`code`, f.response.StatusCode()),
		field.String(`body`, string(f.response.Body())),
	}

	return
}
