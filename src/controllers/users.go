package controllers

import (
	"encoding/json"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// CreateUser inserts a new user to the persistency
func CreateUser(w http.ResponseWriter, r *http.Request) {
	//Read Body request firs
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Json unmarshal into user struct
	var newUser models.User
	if err = json.Unmarshal(bodyRequest, &newUser); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = newUser.Prepare("register"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	//Open DB connection
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	//now we create a new instance of the repository
	repo := repository.UsersNewRepo(db)
	newUser.ID, err = repo.Create(newUser)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newUser)

}

// FetchUsers fetch all users from the persistency
/* I've implemented this method for the CRUD purposes, but it doesn't seem too relevant right now
I might consider different uses for this function. Respository function is commented. Adapt query if planning on using*/
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.UsersNewRepo(db)
	users, erro := repo.FetchAllUsers(nomeOuNick)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// FetchUser fetches a user from the persistency by userID
func FetchUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Err(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := persistency.Connect()
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repo := repository.UsersNewRepo(db)
	user, erro := repo.FetchByID(userID)
	if erro != nil {
		responses.Err(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

// UpdateUser updates as user from the persistency by userID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Updates user"))
}

// DeleteUser deletes an usser from the persistency by userID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletes a user"))
}
