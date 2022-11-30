package http

import (
	"net/http"
	"time"
)

type config struct {
	// 超时
	Timeout time.Duration `json:"timeout" yaml:"timeout" xml:"timeout" toml:"timeout"`
	// 代理
	Proxy *proxy `json:"proxy" yaml:"proxy" xml:"proxy" toml:"proxy"`
	// 授权配置
	Auth *auth `json:"auth" yaml:"auth" xml:"auth" toml:"auth"`
	// Body数据传输控制
	Payload *payload `json:"payload" yaml:"payload" xml:"payload" toml:"payload"`
	// 秘钥配置
	Certificate *certificate `json:"certificate" yaml:"certificate" xml:"certificate" toml:"certificate"`
	// 通用的查询参数
	Queries map[string]string `json:"queries" yaml:"queries" xml:"queries" toml:"queries"`
	// 表单参数
	// 只对POST和PUT方法有效
	Forms map[string]string `json:"forms" yaml:"forms" xml:"forms" toml:"forms"`
	// 通用头信息
	Headers map[string]string `json:"headers" yaml:"headers" xml:"headers" toml:"headers"`
	// 通用Cookie
	Cookies []*http.Cookie `json:"cookies" yaml:"cookies" xml:"cookies" toml:"cookies"`
}
