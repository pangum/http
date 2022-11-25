package http

const (
	authTypeBasic authType = "basic"
	authTypeToken authType = "token"
)

type authType string
