package config

type Certificate struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 是否跳过证书检查
	Skip *bool `default:"true" json:"skip" yaml:"skip" xml:"skip" toml:"skip"`
	// 根秘钥文件路径
	Root string `json:"root" yaml:"root" validate:"file"`
	// 客户端
	Clients []ClientCertificate `json:"clients" yaml:"clients" xml:"clients" toml:"clients" validate:"structonly"`
}

func (c *Certificate) Skipped() bool {
	return nil != c.Skip && *c.Skip
}
