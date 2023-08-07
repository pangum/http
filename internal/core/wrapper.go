package core

type Wrapper struct {
	Http struct {
		// 客户端配置
		Client Config `json:"client" yaml:"client" xml:"client" toml:"client" validate:"required"`
	} `json:"http" yaml:"http" xml:"http" toml:"http" validate:"required"`
}
