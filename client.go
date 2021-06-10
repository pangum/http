package http

import (
	`net/http`
	`time`

	`github.com/storezhang/gox`
)

// ClientConfig 客户端配置
type ClientConfig struct {
	// 超时
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	// 代理
	Proxy gox.ProxyConfig `json:"proxy" yaml:"proxy" validate:"structonly"`
	// 授权配置
	Auth gox.AuthConfig `json:"auth" yaml:"auth" validate:"structonly"`
	// 是否允许Get方法传输数据
	AllowGetPayload bool `default:"true" json:"allowGetPayload" yaml:"allowGetPayload"`
	// 秘钥配置
	Certificate gox.CertificateConfig `json:"certificate" yaml:"certificate" validate:"structonly"`
	// 通用的查询参数
	Queries map[string]string `json:"queries" yaml:"queries"`
	// 表单参数，只对POST和PUT方法有效
	Forms map[string]string `json:"forms" yaml:"forms"`
	// 通用头信息
	Headers map[string]string `json:"headers" yaml:"headers"`
	// 通用Cookie
	Cookies []*http.Cookie `json:"cookies" yaml:"cookies"`
}
