package config

type ContextId string

const (
	APIversion = "/api/v1"
	AuthHeader = "Authorization"
	AuthSchema = "Bearer "
	IdContextKey = ContextId("id")
	AdminRole = string("admin")
	UserRole = string("user")
)