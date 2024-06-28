package handlers

import (
	"net/http"
	"time"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	security "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	utils "github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func CommingSoon(w http.ResponseWriter, r *http.Request) {

}

func OnGoing(w http.ResponseWriter, r *http.Request) {
	
}


// Signup handler for both user/admin
func Signup(w http.ResponseWriter, r *http.Request) {
	credentials := r.Context().Value(config.CredentialsContextKey).(models.UserAdminCredentials)

	//TODO: verification for email address

	hashPassword, err := security.GenerateBcryptPassword(credentials.Password)
	if err != nil {
		utils.JSONResponse(&w, "error encrypting password", http.StatusInternalServerError)
	}
	credentials.Password = hashPassword

	//TODO: save the user details to db according to roles

	id := 32 //TODO: take the id from the db
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


// Login handler for both user/admin
func Login(w http.ResponseWriter, r *http.Request) {
	// credentials := r.Context().Value(config.CredentialsContextKey).(models.UserAdminCredentials)

	//TODO: //validate the login from db by using bcrypt

	id := 32 //TODO: take the id from the db
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
