package http

import (
	"fmt"
	"net/url"

	"github.com/goexl/gox"
)

type proxy struct {
	// Host 主机（可以是Ip或者域名）
	Host string `json:"host" yaml:"host" xml:"host" toml:"host" validate:"required"`
	// Scheme 代理类型
	Scheme gox.UriScheme `default:"scheme" json:"scheme" yaml:"scheme" xml:"scheme" toml:"scheme" validate:"required,oneof=socks4 socks5 http https"`
	// Username 代理认证用户名
	Username string `json:"username" yaml:"username" xml:"username" toml:"username"`
	// Password 代理认证密码
	Password string `json:"password" yaml:"password" xml:"password" toml:"password"`
}

func (p *proxy) Addr() (addr string) {
	if "" != p.Username && "" != p.Password {
		addr = fmt.Sprintf(
			"%s://%s:%s@%s",
			p.Scheme,
			url.QueryEscape(p.Username), url.QueryEscape(p.Password),
			p.Host,
		)
	} else {
		addr = fmt.Sprintf("%s://%s", p.Scheme, p.Host)
	}

	return
}
