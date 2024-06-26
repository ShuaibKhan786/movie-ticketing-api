package services

import (
	"errors"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	Id int
	Exp int64  
}


func GenerateJWTtoken(secretKey []byte, claims Claims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": claims.Id,
			"exp": claims.Exp,
		},
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func ParseJWTtoken(secretKey []byte, tokenString string) (Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id, idOk := claims["id"].(float64)
		exp, expOk := claims["exp"].(float64)

		if !idOk || !expOk {
			fmt.Println(idOk,id)
			fmt.Println(expOk,exp)
			return Claims{}, errors.New("invalid token claims")
		}

		return Claims{
			Id: int(id),
			Exp: int64(exp),
		}, nil

	} else {
		return Claims{}, errors.New("invalid token")
	}
}

