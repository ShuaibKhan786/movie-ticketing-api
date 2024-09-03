package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	security "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	utils "github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

//TODO: do some clean up and refactor the code

func EnsureTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//first it check the Authorization cookie is there or not
		cookie, err := r.Cookie(config.JWTAuthCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				utils.JSONResponse(&w, "cookie not found", http.StatusBadGateway)
			default:
				utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		tokenString := cookie.Value
		secretKey := config.Env.JWTSECRETKEY

		//then it parse the tokenString along with verification
		claims, err := security.ParseJWTtoken(secretKey, tokenString)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenExpired):
				if isURIContainsThatPattern(r.URL.Path, "/refresh/token") {
					goto refreshToken //jump too refreshToken
				}else {
					utils.JSONResponse(&w, "access token has expired", http.StatusUnauthorized)
					return
				}
			default:
				utils.JSONResponse(&w,"access token has tampered",http.StatusForbidden)
				return
			}	
		}

		if isURIContainsThatPattern(r.URL.Path, "/refresh/token") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//making sure only the admin role can access admin route 
		//		-can registered a hall
		//		-can get the hall metadata
		if isURIContainsThatPattern(r.URL.Path,"/admin/")  {
			if claims.Role != config.AdminRole {
				utils.JSONResponse(&w, "this route is for admin only", http.StatusBadRequest)
				return
			}
		}

refreshToken:
		ctx := context.WithValue(r.Context(), config.ClaimsContextKey, claims)

		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

func isURIContainsThatPattern(uri , pattern string) bool {
	return strings.Contains(uri, pattern)
}
