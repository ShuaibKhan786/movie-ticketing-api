package handlers

import (
	"io"
	"net/http"
	"time"

	config "github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	models "github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services"
	utils "github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

func CommingSoon(w http.ResponseWriter, r *http.Request) {

}

func OnGoing(w http.ResponseWriter, r *http.Request) {
	
}


// Signup handler for both user/admin
func Signup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, "failed to read the body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var credentials models.UserAdminCredentials
	if err := utils.DecodeJson(body, &credentials); err != nil {
		utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		return
	}

	if !utils.ValidateLoginOrSigin(&credentials) {
		utils.JSONResponse(&w, "missing credentials", http.StatusBadRequest)
		return
	}

	hashPassword, err := services.GenerateBcryptPassword(credentials.Password)
	if err != nil {
		utils.JSONResponse(&w, "error encrypting password", http.StatusInternalServerError)
	}
	credentials.Password = hashPassword

	//TODO: save the user details to db according to roles

	id := 32 //TODO: take the id from the db
	expirationTime := time.Now().Add(time.Hour * 24).Unix()
	claims := utils.Claims{
		Id:  id,
		Exp: expirationTime,
	}

	tokenString, err := utils.GenerateJWTtoken([]byte("secret-key"), claims)
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, "failed to read the body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var credentials models.UserAdminCredentials
	if err := utils.DecodeJson(body, &credentials); err != nil {
		utils.JSONResponse(&w, "invalid credentials", http.StatusBadRequest)
		return
	}

	if !utils.ValidateLoginOrSigin(&credentials) {
		utils.JSONResponse(&w, "missing credentials", http.StatusBadRequest)
		return
	}

	//TODO: //validate the login from db by using bcrypt

	id := 32 //TODO: take the id from the db
	expirationTime := time.Now().Add(time.Hour * 24).Unix()
	claims := utils.Claims{
		Id:  id,
		Exp: expirationTime,
	}

	tokenString, err := utils.GenerateJWTtoken([]byte("secret-key"), claims)
	if err != nil {
		utils.JSONResponse(&w, "error generating tokens", http.StatusInternalServerError)
		return
	}

	bearerSchema := config.AuthSchema + tokenString
	w.Header().Set(config.AuthHeader, bearerSchema)
	utils.JSONResponse(&w,"login successfully",http.StatusOK)
}
