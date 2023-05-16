package http

type certificate struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 是否跳过证书检查
	Skip *bool `default:"true" json:"skip" yaml:"skip" xml:"skip" toml:"skip"`
	// 根秘钥文件路径
	Root string `json:"root" yaml:"root" validate:"file"`
	// 客户端
	Clients []clientCertificate `json:"clients" yaml:"clients" xml:"clients" toml:"clients" validate:"structonly"`
}

func (c *certificate) skipped() bool {
	return nil != c.Skip && *c.Skip
}
