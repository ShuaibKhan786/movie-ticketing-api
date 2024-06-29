package middlewares

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
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
		credential, ok := r.Context().Value(config.CredentialsContextKey).(models.UserAdminCredentials) 
		if !ok {
			utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		}

		query := database.QueryIsExists
		exists, err := query.DBIsExists(credential.Role,"email",credential.Email)
		if err != nil {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// if it is for signup / register
		// exists must be false 
		// means thier email must not exists in the db
		// inorder to procced to the next handler
		if exists && strings.Contains(r.URL.RequestURI(),"signup") {
			utils.JSONResponse(&w, "email already exists", http.StatusConflict)
			return
		}

		// if it is for login
		// exists must be true 
		// means thier email must present first in db
		// inorder to procced to the next handler
		if !exists && strings.Contains(r.URL.RequestURI(),"login") {
			utils.JSONResponse(&w, "email does not exists", http.StatusConflict)
			return
		}

		next.ServeHTTP(w, r)
	})
}