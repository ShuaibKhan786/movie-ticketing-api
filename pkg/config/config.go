package config

type ContextKeyType string

const (
	APIversion string = "/api/v1"
	AuthHeader string = "Authorization"
	AuthSchema string = "Bearer "
	IdContextKey ContextKeyType = "id"
	AdminRole string = "admin"
	UserRole string = "user"
	BcryptHashingCost int = 12 // it will do 2^12 rounds of hashing
)