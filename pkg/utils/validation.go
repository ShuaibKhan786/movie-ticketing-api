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

func ValidateHallMDUpd(m map[string]map[string]interface{}) error {
	for category, fields := range m {
		switch category {
		case "hall":
			for field, value := range fields {
				if _, valid := config.ValidSetsField.Hall[field]; !valid {
					return errors.New("invalid field '" + field + "' in category '" + category + "'")
				}
				if strValue, ok := value.(string); ok && strValue == "" {
					return errors.New("field '" + field + "' in category '" + category + "' has an empty value")
				}
			}
			fields["tName"] = "hall"
			fields["idName"] = "id"
		case "location":
			for field, value := range fields {
				if _, valid := config.ValidSetsField.Location[field]; !valid {
					return errors.New("invalid field '" + field + "' in category '" + category + "'")
				}
				if strValue, ok := value.(string); ok && strValue == "" {
					return errors.New("field '" + field + "' in category '" + category + "' has an empty value")
				}
			}
			fields["tName"] = "hall_location"
			fields["idName"] = "hall_id"
		case "operation":
			for field, value := range fields {
				if _, valid := config.ValidSetsField.Operation[field]; !valid {
					return errors.New("invalid field '" + field + "' in category '" + category + "'")
				}
				if strValue, ok := value.(string); ok && strValue == "" {
					return errors.New("field '" + field + "' in category '" + category + "' has an empty value")
				}
			}
			fields["tName"] = "hall_operation_time"
			fields["idName"] = "hall_id"
		default:
			return errors.New("invalid category: " + category)
		}
	}
	return nil
}

func ValidateBookedRequestPayload(payload models.BookedRequestPayload) bool {
	if payload.ID == nil {
		return false
	}

	if payload.Counts == nil {
		return false
	}

	if payload.Seats == nil {
		return false
	}

	if payload.PayableAmount == nil {
		return false
	}

	if payload.PaymentMode == nil {
		return false
	}

	if payload.CustomerPhoneNo == nil {
		return false
	}

	return false
}