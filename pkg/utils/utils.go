package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
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

    return base64.StdEncoding.EncodeToString(token), nil
}

func ConvertToSeconds(date, timing string) (int64, error) {
	sqlLayout := "2006-01-02 15:04:05"
	dateTime := fmt.Sprintf("%s %s", date, timing)

	targetTime, err := time.Parse(sqlLayout, dateTime)
	if err != nil {
		return 0, err
	}

	currentTime := time.Now()

	// Calculate TTL in seconds
	ttl := targetTime.Sub(currentTime).Seconds()

	if ttl < 0 {
		return 0, fmt.Errorf("the target time is in the past")
	}

	return int64(ttl), nil
}