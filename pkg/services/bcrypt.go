package services

import "golang.org/x/crypto/bcrypt"

func GenerateBcryptPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.MaxCost) //cost can be changes according too
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func CompareBcryptPassword(hashPassword , password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}