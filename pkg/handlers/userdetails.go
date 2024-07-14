package handlers

import (
	"fmt"
	"net/http"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)



func UserDetails(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}

	id := claims.Id
	role := claims.Role


	userDetails, err := database.GetUserDetails(role, id)
	if err != nil {
		fmt.Println(err)
		utils.JSONResponse(&w, "error in executing db", http.StatusInternalServerError)
		return
	}
	
	jsonData, err := utils.EncodeJson(&userDetails)
	if err != nil {
		utils.JSONResponse(&w, "error in encoding json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}