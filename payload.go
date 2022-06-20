package http

type payload struct {
	// 是否允许Get方法使用Bogy传输数据
	Get bool `default:"true" json:"get" yaml:"get" xml:"get" toml:"get"`
}
