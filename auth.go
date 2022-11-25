package http

type auth struct {
	// Type 授权类型
	Type authType `default:"type" json:"type" yaml:"type" validate:"oneof=basic token"`
	// Username 用户名
	Username string `json:"username" yaml:"username"`
	// Password 密码
	Password string `json:"password" yaml:"password"`
	// Token 授权码
	Token string `json:"token" yaml:"token"`
	// Scheme 身份验证方案类型
	Scheme string `json:"scheme" yaml:"scheme"`
}
