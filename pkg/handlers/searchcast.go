package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

// URL schema:
// http://localhost:3090/api/v1/auth/admin/cast?search_role=actor&search_name=Christian Bale
func SearchCast(w http.ResponseWriter, r *http.Request) {
	role, name, err := getRoleAndNameFromTheQuery(r)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	casts, err := database.SearchCastByName(ctx, role, name)
	if err != nil {
		utils.JSONResponse(&w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(casts)

	if len(casts) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonMovies, err := utils.EncodeJson(&casts)
	if err != nil {
		utils.JSONResponse(&w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(jsonMovies)
}

func getRoleAndNameFromTheQuery(r *http.Request) (string, string, error) {
	role := r.URL.Query().Get("search_role")
	if role == "" {
		return "", "", errors.New("missing or empty 'search_role' query parameter")
	}

	name := r.URL.Query().Get("search_name")
	if name == "" {
		return "", "", errors.New("missing or empty 'search_name' query parameter")
	}

	if !isThisValidCastRole(role) {
		return "", "", errors.New("invalid search_role")
	}

	role = strings.Trim(role, `"`)
	name = strings.Trim(name, `"`)

	return role, name, nil
}

func isThisValidCastRole(role string) bool {
	switch role {
	case "actor",
		"actress",
		"director",
		"producer":
		return true
	default:
		return false
	}
}