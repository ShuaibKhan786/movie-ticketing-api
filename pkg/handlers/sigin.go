package handlers

import (
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

type SignInResponse struct {
	AuthCodeURL string `json:"url"`
}


// payload:
// {
//		Role: "admin",
//		Provider: "google",
//		RedirectedURL: "http://localhost:5173/dashboard",
//		Origin: "http://localhost:5173"	
//}
// SignIn handles the OAuth sign-in process
// 	- Step 0: Read the request body
// 	- Step 2: Validate the JSON payload
// 	- Step 3: Decode the credentials
// 	- Validate the credentials
// 	- Step 4: Configure the OAuth2.0 provider
// 	- Step 5: Generate the authorization URL
// 	- Step 6: Encode the credentials to a base64 string
// 	- Step 7: Encode the response to JSON
// 	- Step 8: Set the encoded credentials as a cookie
// 	- Step 9: Send the response to the client
func SignIn(w http.ResponseWriter, r *http.Request) {
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, "Failed to read the request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	
	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	
	var credentials models.SignInCredentials
	if err := utils.DecodeJson(body, &credentials); err != nil {
		utils.JSONResponse(&w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	
	if err := utils.ValidateSignInCredentials(&credentials); err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	
	oauthConfig, err := security.ConfigOauth(credentials.Provider)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	
	authResponse := SignInResponse{
		AuthCodeURL: oauthConfig.AuthCodeURL(config.Env.OAUTH_STATE),
	}

	
	encodedCredentials := encodeToBase64(body)

	
	jsonResponse, err := utils.EncodeJson(authResponse)
	if err != nil {
		utils.JSONResponse(&w, "Failed to encode response JSON", http.StatusInternalServerError)
		return
	}
	
	cookieExp := time.Now().Add(15 * time.Minute)
	utils.SetCookie(&w, config.OAuthCookieName, encodedCredentials, cookieExp)

	
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonResponse)
}



func encodeToBase64(body []byte) string {
	return base64.StdEncoding.EncodeToString(body)
}