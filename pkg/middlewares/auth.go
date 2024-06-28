package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func EnsureTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerSchema := r.Header.Get(config.AuthHeader)

		if bearerSchema == "" {
			utils.JSONResponse(&w,"token missing",http.StatusUnauthorized)
			return
		}

		tokenString := bearerSchema[len(config.AuthSchema):]
		secretKey := config.Env.JWTSECRETKEY

		claims, err := services.ParseJWTtoken(secretKey, tokenString)
		if err != nil {
			utils.JSONResponse(&w,err.Error(),http.StatusUnauthorized)
			return
		}

		if checkExpiry(claims.Exp) {
			utils.JSONResponse(&w,"token expired",http.StatusUnauthorized)
			return
		}
		//TODO: verify the id from the db

		ctx := context.WithValue(r.Context(), config.IdContextKey, claims.Id)

		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

func checkExpiry(exp int64) bool {
	return time.Now().Unix() > exp
}