package http

import (
	"fmt"
	"net/url"

	"github.com/goexl/gox"
)

type proxy struct {
	// Host 主机（可以是Ip或者域名）
	Host string `json:"ip" yaml:"ip" validate:"required"`
	// Port 端口
	Port int `default:"80" json:"port" yaml:"port" validate:"required"`
	// Scheme 代理类型
	Scheme gox.UriScheme `default:"scheme" json:"scheme" yaml:"type" validate:"required,oneof=socks4 socks5 http https"`
	// Username 代理认证用户名
	Username string `json:"username" yaml:"username"`
	// Password 代理认证密码
	Password string `json:"password" yaml:"password"`
}

func (p *proxy) Addr() (addr string) {
	if "" != p.Username && "" != p.Password {
		addr = fmt.Sprintf(
			"%s://%s:%s@%s:%d",
			p.Scheme,
			url.QueryEscape(p.Username), url.QueryEscape(p.Password),
			p.Host, p.Port,
		)
	} else {
		addr = fmt.Sprintf("%s://%s:%d", p.Scheme, p.Host, p.Port)
	}

	return
}
