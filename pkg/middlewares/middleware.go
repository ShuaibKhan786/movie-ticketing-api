package middlewares

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs)-1 ; i >= 0 ; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before")
		defer log.Println("After")
		next.ServeHTTP(w,r) //actual function that wrapped up
	})
}