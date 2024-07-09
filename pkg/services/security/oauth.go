package security

import (
	"errors"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	oauth2 "golang.org/x/oauth2"
	google "golang.org/x/oauth2/google"
)

func ConfigOauth(provider string) (*oauth2.Config, error) {
	switch provider {
	case "google":
		return &oauth2.Config{
			ClientID: config.Env.GOOGLE_CLIENT_ID,
			ClientSecret: config.Env.GOOGLE_CLIENT_SECRET,
			Endpoint: google.Endpoint,
			RedirectURL: config.Env.REDIRECT_URL,
			Scopes: []string{
				config.Env.GOOGLE_SCOPE_EMAIL_URL,
				config.Env.GOOGLE_SCOPE_PROFILE_URL,
			},
		}, nil
	default:
		return nil, errors.New("unsupported provider")
	}
}

func GetUserInfoUrl(provider string) (string, error) {
	switch provider {
	case "google":
		return config.Env.GOOGLE_USERINFO_URL, nil
	default:
		return "", errors.New("unsupported provider")
	}
}