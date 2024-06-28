package utils

import (
	"encoding/json"
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
