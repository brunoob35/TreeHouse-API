package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/gorilla/mux"
)

// CreateUser is responsible for creating a new user.
//
// This flow:
//   - reads the request body
//   - parses the incoming JSON into a user struct
//   - prepares and validates the user data
//   - inserts the base user record into the database
//
// Permission assignments are not handled here and must be managed separately
// through the relation table "usuarios_permissoes".
func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var newUser models.User
	if err = json.Unmarshal(bodyRequest, &newUser); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = newUser.Prepare("create"); err != nil {
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

	newUser.ID, err = repo.Insert(newUser)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	newUser.Senha = ""

	responses.JSON(w, http.StatusCreated, newUser)
}

// FetchUsers returns all users filtered by the optional "nome" query parameter.
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	nome := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("nome")))

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewUsersRepository(db)

	users, err := repo.FetchAllUsers(nome)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// FetchUser returns a single user by its ID.
func FetchUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
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

	user, err := repo.FetchByID(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	user.Senha = ""

	responses.JSON(w, http.StatusOK, user)
}

// UpdateUser updates an existing user's base data.
//
// This endpoint updates only the user base fields stored in "usuarios".
// Permission assignments must be updated separately through the repository
// method responsible for replacing entries in "usuarios_permissoes".
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

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

	if err = user.Prepare("update"); err != nil {
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

	if err = repo.Update(userID, user); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser removes a user by its ID.
//
// Since the relation table "usuarios_permissoes" uses ON DELETE CASCADE,
// all permission relations linked to the deleted user are automatically
// removed by the database.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, err := strconv.ParseUint(params["userID"], 10, 64)
	if err != nil {
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

	if err = repo.Delete(userID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
