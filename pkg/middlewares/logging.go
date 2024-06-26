package middlewares

import (
	"net/http"
	"log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: logged to a file in production
		log.Println("Before ",r.RemoteAddr,r.Method)
		defer log.Println("After",r.RemoteAddr,r.Method)
		next.ServeHTTP(w,r) //actual function that wrapped up
	})
}