package http

type certificate struct {
	// 是否跳过TLS检查
	Skip bool `default:"true" json:"skip" yaml:"skip"`
	// 根秘钥文件路径
	Root string `json:"root" yaml:"root" validate:"required,file"`
	// 客户端
	Clients []clientCertificate `json:"clients" yaml:"clients" validate:"structonly"`
}
