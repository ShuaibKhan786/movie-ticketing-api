package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie(config.RefreshTokenCookieName)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			utils.JSONResponse(&w, "refresh token cookie not found", http.StatusBadRequest)
		default:
			utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	refreshToken := cookie.Value
	redisKey := fmt.Sprintf("%s:%d", claims.Role, claims.Id)
	fmt.Println(redisKey)
	redisCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	redisRefreshToken, err := redisdb.Get(redisCtx, redisKey)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	if redisRefreshToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if redisRefreshToken != refreshToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	claims.Exp = time.Now().Add(time.Hour * 1).Unix()
	newTokenString, err := security.GenerateJWTtoken(config.Env.JWTSECRETKEY, claims)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	cookieExp := time.Now().Add(time.Hour * 24 * 7)
	utils.SetCookie(&w, config.JWTAuthCookieName, newTokenString, cookieExp)
	w.WriteHeader(http.StatusNoContent)
}
