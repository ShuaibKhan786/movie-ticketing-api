package middlewares

import (
	"context"
	"io"
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// middleware for validating specific JSON 
func IsValidJSONCred(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.JSONResponse(&w, "failed to read the body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		//validating the request body is json or not
		if !utils.IsValidJson(body) {
			utils.JSONResponse(&w, "invalid request payload", http.StatusBadRequest)
			return
		}

		//decoding the json body into struct 
		var credentials models.UserAdminCredentials
		if err := utils.DecodeJson(body, &credentials); err != nil {
			utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(),config.CredentialsContextKey,credentials)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//middleware for validating specific credentials
func IsValidCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		credentials := r.Context().Value(config.CredentialsContextKey).(models.UserAdminCredentials)
		if err := utils.ValidateLoginOrSigninCredentials(&credentials); err != nil {
			utils.JSONResponse(&w,err.Error(),http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//middleware for checking wether the email exits or not
func IsEmailExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: look up a db and checking wether the email exits already or not
		//	1.Parse down the email address credentials.Email
		//	2.Do db transiction 
		//	3.According to db result
		//		-send utils.JSONResponse(&w,"email already exists",http.StatusFound)
		//	or
		next.ServeHTTP(w, r)
	})
}