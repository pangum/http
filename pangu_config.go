package http

type panguConfig struct {
	// Http配置
	Http struct {
		// 客户端配置
		Client config `json:"client" yaml:"client" xml:"client" toml:"client" validate:"required"`
	} `json:"http" yaml:"http" xml:"http" toml:"http" validate:"required"`
}
