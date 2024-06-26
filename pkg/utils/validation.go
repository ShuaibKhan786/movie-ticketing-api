package utils

import (
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
)

func  ValidateLoginOrSigin(loginCredentials *models.UserAdminCredentials) bool {
	var state bool 
	if loginCredentials.Role == config.AdminRole || loginCredentials.Role == config.UserRole {
		state = true
	}
	if loginCredentials.Email == "" || loginCredentials.Email == " " {
		state = false
	}
	if loginCredentials.Password == "" || loginCredentials.Password == " " {
		state = false
	}
	return state
}