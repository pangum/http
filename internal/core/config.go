package core

import (
	"net/http"
	"time"

	"github.com/pangum/http/internal/config"
)

type Config struct {
	// 超时
	Timeout time.Duration `json:"timeout" yaml:"timeout" xml:"timeout" toml:"timeout"`
	// 代理
	Proxy *config.Proxy `json:"Proxy" yaml:"Proxy" xml:"Proxy" toml:"Proxy"`
	// 代理列表
	Proxies []*config.Proxy `json:"proxies" yaml:"proxies" xml:"proxies" toml:"proxies"`
	// 授权配置
	Auth *config.Auth `json:"Auth" yaml:"Auth" xml:"Auth" toml:"Auth"`
	// Body数据传输控制
	Payload *config.Payload `json:"Payload" yaml:"Payload" xml:"Payload" toml:"Payload"`
	// 秘钥配置
	Certificate *config.Certificate `json:"Certificate" yaml:"Certificate" xml:"Certificate" toml:"Certificate"`
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
