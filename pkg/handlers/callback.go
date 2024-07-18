package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"golang.org/x/oauth2"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	credentials, err := getCredentialsFromCookie(r)
	if err != nil {
		redirectWithError(w, r, err.Error())
		return
	}

	if err := validateOAuthState(r); err != nil {
		redirectWithError(w, r, err.Error())
		return
	}

	token, err := exchangeOAuthCodeForToken(r, credentials)
	if err != nil {
		redirectWithError(w, r, err.Error()+"token exchange failed")
		return
	}

	userDetails, err := getUserDetailsFromProvider(token, credentials.Provider)
	if err != nil {
		redirectWithError(w, r, err.Error())
		return
	}

	if err := registerUserIfNotExists(&w, credentials, userDetails); err != nil {
		redirectWithError(w, r, err.Error())
		return
	}

	if err := generateAndSetTokens(w, credentials, userDetails); err != nil {
		redirectWithError(w, r, err.Error())
		return
	}

	utils.DeleteCookie(&w, config.OAuthCookieName)
	http.Redirect(w, r, credentials.RedirectedURL, http.StatusPermanentRedirect)
}




func getCredentialsFromCookie(r *http.Request) (models.SignInCredentials, error) {
	cookie, err := r.Cookie(config.OAuthCookieName)
	if err != nil {
		return models.SignInCredentials{}, err
	}

	return decodeFromBase64(cookie.Value)
}

func validateOAuthState(r *http.Request) error {
	state := r.URL.Query().Get("state")
	if state != config.Env.OAUTH_STATE {
		return fmt.Errorf("invalid state")
	}
	return nil
}

func exchangeOAuthCodeForToken(r *http.Request, credentials models.SignInCredentials) (*oauth2.Token, error) {
	code := r.URL.Query().Get("code")
	if code == "" {
		return nil, fmt.Errorf("missing authorization code")
	}

	oauthConfig, err := security.ConfigOauth(credentials.Provider)
	if err != nil {
		return nil, fmt.Errorf("invalid provider configuration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return oauthConfig.Exchange(ctx, code)
}

func getUserDetailsFromProvider(token *oauth2.Token, provider string) (map[string]interface{}, error) {
	oauthConfig, err := security.ConfigOauth(provider)
	if err != nil {
		return nil, err
	}

	client := oauthConfig.Client(context.Background(), token)
	userDetailsUrl, err := security.GetUserInfoUrl(provider)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(userDetailsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userDetails map[string]interface{}
	if err := utils.DecodeJson(body, &userDetails); err != nil {
		return nil, err
	}
	return userDetails, nil
}

func registerUserIfNotExists(w *http.ResponseWriter,credentials models.SignInCredentials, userDetails map[string]interface{}) error {
	email := userDetails["email"].(string)
	exists, err := database.IsValueExists(credentials.Role, "email", email)
	if err != nil {
		return fmt.Errorf("existance check failed: %w", err)
	}

	if !exists {
		modelUserDetails := structureUserDetails(userDetails)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := database.RegisterUserDetails(ctx, modelUserDetails, credentials.Role);
		err != nil {
			return err
		}
	}
	return nil
}

func generateAndSetTokens(w http.ResponseWriter, credentials models.SignInCredentials, userDetails map[string]interface{}) error {
	email := userDetails["email"].(string)
	id, err := database.GetId(credentials.Role, "email", email)
	if err != nil {
		return err
	}

	jwtExp := time.Now().Add(time.Minute * 60).Unix()
	claims := security.Claims{
		Id:   id,
		Role: credentials.Role,
		Exp:  jwtExp,
	}
	jwtTokenString, err := security.GenerateJWTtoken(config.Env.JWTSECRETKEY, claims)
	if err != nil {
		return err
	}

	cookieExp := time.Now().Add(time.Hour * 24 * 7)
	utils.SetCookie(&w, config.JWTAuthCookieName, jwtTokenString, cookieExp)

	refreshToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		return fmt.Errorf("error generating refresh token: %w", err)
	}

	redisKey := fmt.Sprintf("%s:%d", credentials.Role, id)
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer redisCancel()
	redisExp := time.Hour * 24 * 7
	if err := redisdb.Set(redisCtx, redisKey, refreshToken, redisExp); err != nil {
		return err
	}

	utils.SetCookie(&w, config.RefreshTokenCookieName, refreshToken, cookieExp)
	return nil
}

func redirectWithError(w http.ResponseWriter, r *http.Request, errorMsg string) {
	origin := r.URL.Query().Get("origin")
	if origin == "" {
		origin = config.Env.DEFAULT_ORIGIN
	}
	http.Redirect(w, r, origin+"?error="+errorMsg, http.StatusTemporaryRedirect)
}

func decodeFromBase64(encodedCredentials string) (models.SignInCredentials, error) {
	var credentials models.SignInCredentials
	jsonData, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return credentials, err
	}
	err = utils.DecodeJson(jsonData, &credentials)
	if err != nil {
		return credentials, err
	}
	return credentials, nil
}

func structureUserDetails(userDetails map[string]interface{}) models.UserDetails {
	return models.UserDetails{
		Email:         userDetails["email"].(string),
		EmailVerified: userDetails["verified_email"].(bool),
		Profile: models.Profile{
			Name:      userDetails["name"].(string),
			PosterUrl: userDetails["picture"].(string),
		},
	}
}
