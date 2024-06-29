package handlers

import (
	"errors"
	"net/http"
	"time"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	database "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	security "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	utils "github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// Login handler for both user/admin
func Login(w http.ResponseWriter, r *http.Request) {
	credentials, ok := r.Context().Value(config.CredentialsContextKey).(models.UserAdminCredentials)
	if !ok {
		utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		return
	}

	var id int64
	query := database.QueryLogin
	var password string

	if err := query.DBLogin(credentials, &id, &password); err != nil {
		if errors.Is(err, database.ErrDBPrepareStmt) {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, database.ErrDBNoRows) {
			utils.JSONResponse(&w, "no record found", http.StatusNotFound)
			return
		}

		utils.JSONResponse(&w, "unexpected error", http.StatusInternalServerError)
		return
	}

	if !security.CompareBcryptPassword(password, credentials.Password) {
		utils.JSONResponse(&w, "wrong password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24).Unix()
	claims := security.Claims{
		Id:  id,
		Exp: expirationTime,
	}

	secretKey := config.Env.JWTSECRETKEY
	tokenString, err := security.GenerateJWTtoken(secretKey, claims)
	if err != nil {
		utils.JSONResponse(&w, "error generating tokens", http.StatusInternalServerError)
		return
	}

	bearerSchema := config.AuthSchema + tokenString
	w.Header().Set(config.AuthHeader, bearerSchema)
	utils.JSONResponse(&w,"login successfully",http.StatusOK)
}