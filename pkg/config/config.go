// Advances in Neural Information Processing Systems. Vol. 35. Curran Associates
package config

type ContextKeyClaims string
type ContextKeyCredentials string

const (
	APIversion                string = "/api/v1"
	OAuthCookieName           string = "oauth_credentials"
	JWTAuthCookieName         string = "Authorization"
	JWTRefreshTokenCookieName string = "RefreshToken"

	ClaimsContextKey      ContextKeyClaims      = "claims"
	CredentialsContextKey ContextKeyCredentials = "credentials"

	AdminRole         string = "admin"
	UserRole          string = "user"
	BcryptHashingCost int    = 12 // it will do 2^12 rounds of hashing
	RegexEmail        string = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
)
