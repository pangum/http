package http

import (
	"github.com/goexl/gox"
)

type auth struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 授权类型
	Type authType `default:"type" json:"type" yaml:"type" xml:"type" toml:"type" validate:"oneof=basic token"`
	// 用户名
	Username string `json:"username" yaml:"username" xml:"username" toml:"username"`
	// 密码
	Password string `json:"password" yaml:"password" xml:"password" toml:"password"`
	// 授权码
	Token string `json:"token" yaml:"token" xml:"token" toml:"token" validate:"required_if=Type token"`
	// 身份验证方案类型
	Scheme gox.UriScheme `json:"scheme" yaml:"scheme" xml:"scheme" toml:"scheme"`
}
