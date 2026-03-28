package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

// Login is responsible for validating user credentials.
//
// This flow performs the following steps:
//   - reads the request body
//   - parses the incoming credentials
//   - loads the user by email
//   - validates the provided password
//   - loads the user's permission IDs from the database
//   - aggregates those permissions into a numeric bitmask
//   - generates a JWT token containing the user ID and permission mask
func Login(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewUsersRepository(db)

	// FetchByEmail loads the user base data required for authentication.
	userFound, err := repo.FetchByEmail(user.Email)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// ValidatePassword compares the stored password hash with the plain password
	// received in the login request.
	if err = security.ValidatePassword(userFound.Senha, user.Senha); err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		log.Println("Unauthorized")
		return
	}

	// FetchPermissionMaskByUserID loads the user's permission IDs from the
	// relationship table and aggregates them into a single numeric bitmask.
	permissionMask, err := repo.FetchPermissionMaskByUser(userFound.ID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	token, err := authentication.GenerateToken(userFound.ID, permissionMask)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.Write([]byte(token))
}
