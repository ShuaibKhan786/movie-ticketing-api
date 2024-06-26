package middlewares

import (
	"net/http"
)

func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		//checking for preflight request
		if r.Method == http.MethodOptions {
			//logic for preflight response
			setCorsHeaders(&w)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		//not all request do a preflight request by browser
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		next.ServeHTTP(w, r)
	})
}

func setCorsHeaders(w *http.ResponseWriter) {
	//TODO: according to the needs set the cors headers
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE") 
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Max-Age", "86400")
}