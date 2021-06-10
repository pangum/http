package http

type config struct {
	// Http配置
	Http struct {
		// 客户端配置
		Client ClientConfig `json:"client" yaml:"client" validate:"required"`
	} `json:"http" yaml:"http"`
}
