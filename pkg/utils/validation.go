package utils

import (
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"

	"strings"
)

func ValidateLoginOrSiginCredentials(loginCredentials *models.UserAdminCredentials) bool {
    if loginCredentials.Role != config.AdminRole && loginCredentials.Role != config.UserRole {
        return false
    }

    if strings.TrimSpace(loginCredentials.Email) == "" {
        return false
    }

    password := strings.TrimSpace(loginCredentials.Password)
    if len(password) < 12 || len(password) > 24 {
        return false
    }

    return true
}
