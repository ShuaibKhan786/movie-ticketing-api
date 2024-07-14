package utils

import (
	"encoding/base64"
	"encoding/json"
	"crypto/rand"
)

func IsValidJson(body []byte) bool{
	return json.Valid(body)
}

func DecodeJson(body []byte,v interface{}) error {
	if err := json.Unmarshal(body,v); err != nil {
		return err
	}
	return nil
}

func EncodeJson(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil , err
	}
	return data, nil
}

func GenerateRandomToken(size int) (string, error) {
    token := make([]byte, size, size)
    _, err := rand.Read(token)
    if err != nil {
        return "", err
    }

    return base64.URLEncoding.EncodeToString(token), nil
}