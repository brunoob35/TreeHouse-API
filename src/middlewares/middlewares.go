package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/responses"
)

var errPermissionDenied = errors.New("user does not have permission to access this resource")

// Logger logs request information to the terminal.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n%s %s %s\n", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

// Authentication verifies whether the request contains a valid authenticated token.
//
// This middleware only validates the token integrity and authentication state.
// It does not check route permissions.
func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := authentication.ValidateToken(r); err != nil {
			responses.Err(w, http.StatusUnauthorized, err)
			return
		}

		next(w, r)
	}
}

// Authorize verifies whether the authenticated user has at least one of the
// permissions required to access the route.
//
// This middleware does not query the database.
// It relies entirely on the permission mask stored inside the JWT token.
func Authorize(required ...authentication.Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userPermissions, err := authentication.ExtractPermissions(r)
			if err != nil {
				responses.Err(w, http.StatusUnauthorized, err)
				return
			}

			if !authentication.HasAnyPermission(userPermissions, required...) {
				responses.Err(w, http.StatusForbidden, errPermissionDenied)
				return
			}

			next(w, r)
		}
	}
}

// AuthorizeAll verifies whether the authenticated user has all permissions
// required to access the route.
// This middleware also relies entirely on the JWT token and does not query
// the database.
func AuthorizeAll(required ...authentication.Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userPermissions, err := authentication.ExtractPermissions(r)
			if err != nil {
				responses.Err(w, http.StatusUnauthorized, err)
				return
			}

			if !authentication.HasAllPermissions(userPermissions, required...) {
				responses.Err(w, http.StatusForbidden, errPermissionDenied)
				return
			}

			next(w, r)
		}
	}
}
