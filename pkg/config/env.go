package config

import (
	"os"
	"sync"
)

type ENV struct {
	JWTSECRETKEY                  []byte
	DSN                           string
	PORT                          string
	GOOGLE_CLIENT_ID              string
	GOOGLE_CLIENT_SECRET          string
	GOOGLE_SCOPE_EMAIL_URL        string
	GOOGLE_SCOPE_PROFILE_URL      string
	REDIRECT_URL                  string
	OAUTH_STATE                   string
	DEFAULT_ORIGIN                string
	GOOGLE_USERINFO_URL           string
	REDIS_URL                     string
	GRPC_IMAGE_UPLOAD_SERVER_HOST string
}

var (
	Env            ENV
	ValidSetsField ValidSetFields
	once           sync.Once
)

type ValidSetFields struct {
	Hall      map[string]struct{}
	Operation map[string]struct{}
	Location  map[string]struct{}
}

func LoadConfig() bool {
	var success = true
	once.Do(func() {
		Env.JWTSECRETKEY = []byte(os.Getenv("JWT_SECRET_KEY"))
		Env.DSN = os.Getenv("DSN")
		Env.PORT = os.Getenv("PORT")
		Env.GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
		Env.GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")
		Env.REDIRECT_URL = os.Getenv("REDIRECT_URL")
		Env.GOOGLE_SCOPE_EMAIL_URL = os.Getenv("GOOGLE_SCOPE_EMAIL_URL")
		Env.GOOGLE_SCOPE_PROFILE_URL = os.Getenv("GOOGLE_SCOPE_PROFILE_URL")
		Env.OAUTH_STATE = os.Getenv("OAUTH_STATE")
		Env.DEFAULT_ORIGIN = os.Getenv("DEFAULT_ORIGIN")
		Env.GOOGLE_USERINFO_URL = os.Getenv("GOOGLE_USERINFO_URL")
		Env.REDIS_URL = os.Getenv("REDIS_URL")
		Env.GRPC_IMAGE_UPLOAD_SERVER_HOST = os.Getenv("GRPC_IMAGE_UPLOAD_SERVER_HOST")

		ValidSetsField.Hall = map[string]struct{}{
			"name":    {},
			"manager": {},
			"contact": {},
		}

		ValidSetsField.Location = map[string]struct{}{
			"address":     {},
			"city":        {},
			"state":       {},
			"postal_code": {},
			"latitude":    {},
			"longitude":   {},
		}

		ValidSetsField.Operation = map[string]struct{}{
			"open_time":   {},
			"closed_time": {},
		}
	})
	return success
}


	// if err := godotenv.Load(); err != nil {
	// 	// fmt.Println(err)
	// 	return false
	// }