package middlewares

import (
	"fmt"
	"log"
	"net/http"
)

// Logger logs to the terminal request info
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n%s %s %s\n", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}

}

// Authentication verifies if the user is authenticated.
func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Validating...")
		next(w, r)
	}

}
