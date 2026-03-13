package controllers

import (
	"net/http"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	repositories "github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
)

// CreateProfessor creates a new user and automatically associates permission 2.
// The function calls createUserWithPermission and gives the respective valid permission ID
func CreateProfessor(w http.ResponseWriter, r *http.Request) {
	createUserWithPermission(w, r, authentication.PermProfessor)
}

// FetchProfessors returns all active professors optionally filtered by name.
func FetchProfessors(w http.ResponseWriter, r *http.Request) {
	nome := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("nome")))

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewUsersRepository(db)

	users, err := repo.FetchProfessors(nome)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// ReturnAllProfessors returns all professors optionally filtered by name.
func ReturnAllProfessors(w http.ResponseWriter, r *http.Request) {
	nome := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("nome")))

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewUsersRepository(db)

	users, err := repo.ReturnAllProfessors(nome)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}
