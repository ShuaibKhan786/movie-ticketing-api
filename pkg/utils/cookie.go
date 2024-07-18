package utils

import (
	"net/http"
	"time"
)


func SetCookie(w *http.ResponseWriter, key string, value string, expiry time.Time) {
	cookie := http.Cookie{
		Name: key,
		Value: value,
		Path: "/",
		Expires: expiry,
		Secure: false, //set true in production
		HttpOnly: true,
	}

	http.SetCookie(*w, &cookie)
}

func SetCookieWithNoExpiry(w *http.ResponseWriter, key string, value string) {
	cookie := http.Cookie{
		Name: key,
		Value: value,
		HttpOnly: true,
		Path: "/",
		Secure: false,
		MaxAge: 0,
	}
	http.SetCookie(*w, &cookie)
}

func DeleteCookie(w *http.ResponseWriter, name string) {
    http.SetCookie(*w, &http.Cookie{
        Name:   name,
        Value:  "",
        MaxAge: -1, //since I want to delete cookie after it being set or overwrite
        Path:   "/",
		HttpOnly: true,
    })
}