package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)


func Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	redisKey := fmt.Sprintf("%s:%d", claims.Role, claims.Id)
	redisCtx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond)
	defer cancel()
	if _, err := redisdb.Delete(redisCtx, redisKey); err != nil {
		utils.JSONResponse(&w, "failed to delete the refresh token", http.StatusInternalServerError)
		return
	}
	
	utils.DeleteCookie(&w, config.JWTAuthCookieName)
	utils.DeleteCookie(&w, config.RefreshTokenCookieName)

	utils.JSONResponse(&w, "successfully logout", http.StatusOK)
}