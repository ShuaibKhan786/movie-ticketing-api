package utils

import "net/http"

type Response struct {
	Message string `json:"message"`
}

func JSONResponse(w *http.ResponseWriter, message string, statusCode int) {
	data := &Response {
		Message: message,
	}

	jsonData, err := EncodeJson(data)
	if err != nil {
		http.Error((*w),"internal server error",http.StatusInternalServerError)
	}
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(statusCode)
	(*w).Write(jsonData)
}