package config

const (
	AuthTypeBasic AuthType = "basic"
	AuthTypeToken AuthType = "token"
)

type AuthType string
