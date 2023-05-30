package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/brunoob35/TreeHouse-API/src/model"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"io"
	"log"
	"net/http"
)

// CreateUser inserts a new user to the persistency
func CreateUser(w http.ResponseWriter, r *http.Request) {
	//Read Body request firs
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		//temporary error handling
		log.Fatal(err)
	}

	//Json unmarshal into user struct
	var newUser model.User
	if err = json.Unmarshal(bodyRequest, &newUser); err != nil {
		//temporary error handling
		log.Fatal("Erro no unmarshal ", err)
	}

	//Open DB connection
	db, err := persistency.Connect()
	if err != nil {
		//temporary error handling
		log.Fatal("Erro para conectar no banco: ", err)
	}

	//now we create a new instance of the repository
	repo := repository.UsersNewRepo(db)
	userID, err := repo.Create(newUser)
	if err != nil {
		//temporary error handling
		log.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf("Id inserido: %d", userID)))

}

// FetchUsers fetch all users from the persistency
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Fetches all users"))
}

// FetchUser fetch an un user from the persistency by userID
func FetchUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Fetch an user"))
}

// UpdateUser updates as user from the persistency by userID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Updates user"))
}

// DeleteUser deletes an usser from the persistency by userID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deletes a user"))
}
