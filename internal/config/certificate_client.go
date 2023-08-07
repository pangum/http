package config

type ClientCertificate struct {
	// 公钥文件路径
	Public string `json:"public" yaml:"public" validate:"required,file"`
	// 私钥文件路径
	Private string `json:"private" yaml:"private" validate:"required,file"`
}
