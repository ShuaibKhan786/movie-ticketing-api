package utils

import (
	"errors"
	"time"
	"net/url"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"

	"strings"
)

// ValidateSignInCredentials validates the sign-in credentials
func ValidateSignInCredentials(credentials *models.SignInCredentials) error {
	if credentials.Role != config.AdminRole && credentials.Role != config.UserRole {
		return errors.New("role must be either admin or user")
	}

	if credentials.Provider != "google" {
		return errors.New("only google OAuth provider is supported")
	}

	if !isValidURL(credentials.RedirectedURL) {
		return errors.New("invalid redirect URL")
	}

	return nil
}

func isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}



func ValidateHall(hall *models.Hall) error {
	if strings.TrimSpace(hall.Name) == "" {
		return errors.New("hall name is required")
	}

	if strings.TrimSpace(hall.Manager) == "" {
		return errors.New("hall manager is required")
	}

	if len(strings.TrimSpace(hall.Contact)) < 10 {
		return errors.New("invalid contact information")
	}

	if err := ValidateLocation(&hall.Location); err != nil {
		return err
	}

	if err := ValidateOperationTime(&hall.OperationTime); err != nil {
		return err
	}

	return nil
}

func ValidateLocation(location *models.Location) error {
	if strings.TrimSpace(location.Address) == "" {
		return errors.New("address is required")
	}

	if strings.TrimSpace(location.City) == "" {
		return errors.New("city is required")
	}

	if strings.TrimSpace(location.State) == "" {
		return errors.New("state is required")
	}

	if strings.TrimSpace(location.PostalCode) == "" {
		return errors.New("postal code is required")
	}

	return nil
}


func ValidateOperationTime(operationTime *models.OperationTime) error {
	if !isValidTime(operationTime.OpenTime) {
		return errors.New("invalid open time")
	}

	if !isValidTime(operationTime.CloseTime) {
		return errors.New("invalid close time")
	}

	return nil
}


func isValidTime(timeStr string) bool {
	_, err := time.Parse("15:04:00", timeStr) // assuming time in HH:MM:SS format
	return err == nil
}