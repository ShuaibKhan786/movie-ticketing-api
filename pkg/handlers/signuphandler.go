package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	security "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	utils "github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)


// Signup handler for both user/admin
func Signup(w http.ResponseWriter, r *http.Request) {
	credentials, ok := r.Context().Value(config.CredentialsContextKey).(models.UserAdminCredentials)
	if !ok {
		utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		return
	}

	hashPassword, err := security.GenerateBcryptPassword(credentials.Password)
	if err != nil {
		utils.JSONResponse(&w, "error encrypting password", http.StatusInternalServerError)
	}
	credentials.Password = hashPassword

	
	var id int64

	query := database.QuerySignup

	if err := query.DBsignup(credentials, &id); err != nil {
		if errors.Is(err, database.ErrDBPrepareStmt) {
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, database.ErrDBExec) {
			if strings.Contains(err.Error(), "1062") { // MySQL specific error code for duplicate entry
                utils.JSONResponse(&w, "Email already exists", http.StatusConflict)
                return
            }
			utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, database.ErrDBNoLastInsertedId) {
			query := database.QueryGetId
			if err := query.DBGetId(credentials.Role, "email", credentials.Email, &id); err != nil {
				utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		utils.JSONResponse(&w, "unexpected error", http.StatusInternalServerError)
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
	utils.JSONResponse(&w,"registered successfully",http.StatusOK)
}