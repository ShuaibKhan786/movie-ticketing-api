Advances in Neural Information Processing Systems. Vol. 35. Curran Associatespackage config

type ContextKeyId string
type ContextKeyCredentials string

const (
	APIversion string = "/api/v1"
	AuthHeader string = "Authorization"
	AuthSchema string = "Bearer "

	IdContextKey ContextKeyId = "id"
	CredentialsContextKey ContextKeyCredentials = "credentials"

	AdminRole string = "admin"
	UserRole string = "user"
	BcryptHashingCost int = 12 // it will do 2^12 rounds of hashing
	RegexEmail string = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
)

