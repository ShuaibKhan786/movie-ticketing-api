package security

import (
	"errors"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
)


type Claims struct {
	Id int64
	Role string
	// Hall bool
	// Email bool
	Exp int64  
}


func GenerateJWTtoken(secretKey []byte, claims Claims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id": claims.Id,
			"exp": claims.Exp,
			"role": claims.Role,
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
		// Check if the error is due to token expiration
		if errors.Is(err, jwt.ErrTokenExpired) {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				parsedClaims, errClaims := parseTheClaims(claims)
				if errClaims != nil {
					return Claims{}, errClaims
				}
				return *parsedClaims, err
			}
		}
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		parsedClaims, err := parseTheClaims(claims)
		if err != nil {
			return Claims{}, err
		}
		return *parsedClaims, nil
	}

	return Claims{}, errors.New("invalid token")
}

func parseTheClaims(claims jwt.MapClaims) (*Claims, error) {
	id, idOk := claims["id"].(float64)
	exp, expOk := claims["exp"].(float64)
	role, roleOk := claims["role"].(string)

	if !idOk || !expOk || !roleOk {
		return nil, errors.New("invalid token claims")
	}

	return &Claims{
		Id:   int64(id),
		Exp:  int64(exp),
		Role: role,
	}, nil
}