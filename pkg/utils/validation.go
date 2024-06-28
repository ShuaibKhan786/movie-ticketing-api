package utils

import (
	"errors"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"

	"net/mail"
	"strings"
)

func ValidateLoginOrSigninCredentials(credentials *models.UserAdminCredentials) error {
	if credentials.Role != config.AdminRole && credentials.Role != config.UserRole {
		return errors.New("role must be either admin / user")
	}

	if !isValidEmail(credentials.Email) {
		return errors.New("invalid email address")
	}

	password := strings.TrimSpace(credentials.Password)
	if len(password) < 8 || len(password) > 12 {
		return errors.New("invalid password")
	}

    return nil
}


func isValidEmail(email string) bool{
    _, err := mail.ParseAddress(email)
    return err == nil
}