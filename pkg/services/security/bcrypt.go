package security

import (
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

func GenerateBcryptPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password),config.BcryptHashingCost) 
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func CompareBcryptPassword(hashPassword , password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}