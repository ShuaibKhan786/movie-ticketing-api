package config

import "time"

type ContextKeyClaims string
type ContextKeyCredentials string

const (
	APIversion               string = "/api/v1"
	OAuthCookieName          string = "oauth_credentials"
	JWTAuthCookieName        string = "Authorization"
	RefreshTokenCookieName   string = "RefreshToken"
	HallRegisteredCookieName string = "IsHallRegistered"

	ClaimsContextKey      ContextKeyClaims      = "claims"
	CredentialsContextKey ContextKeyCredentials = "credentials"

	AdminRole         string = "admin"
	UserRole          string = "user"
	BcryptHashingCost int    = 12 // it will do 2^12 rounds of hashing
	RegexEmail        string = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`

	RedisZeroExpirationTime time.Duration = 0

	Cash = "cash"
	UPI = "upi"
)
