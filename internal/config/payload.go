package config

type Payload struct {
	// 是否允许Get方法使用Body传输数据
	Get bool `default:"true" json:"get" yaml:"get" xml:"get" toml:"get"`
}
