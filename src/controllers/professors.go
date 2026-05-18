// src/controllers/professors.go
package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	repositories "github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/gorilla/mux"
)

func CreateProfessor(w http.ResponseWriter, r *http.Request) {
	createUserWithPermission(w, r, authentication.PermProfessor)
}

func FetchCurrentProfessor(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
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
		if err == sql.ErrNoRows {
			responses.Err(w, http.StatusNotFound, err)
			return
		}
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

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

func CountProfessorClasses(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var request models.ProfessorClassCountRequest
	if err = json.Unmarshal(bodyRequest, &request); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if len(request.ProfessorIDs) == 0 {
		responses.Err(w, http.StatusBadRequest, errors.New("professor_ids é obrigatório"))
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewUsersRepository(db)

	counts, err := repo.CountClassesByProfessorIDs(request.ProfessorIDs)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, counts)
}

func AssignProfessorToClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	professorID, err := strconv.ParseUint(params["professorID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	classID, err := strconv.ParseUint(params["classID"], 10, 64)
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

	repo := repositories.NewClassesRepository(db)

	if err = repo.AssignProfessorToClass(classID, professorID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func FetchCurrentProfessorClasses(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewClassesRepository(db)
	classes, err := repo.FetchAllActiveByTeacherID(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, classes)
}

func FetchCurrentProfessorLessons(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewLessonsRepository(db)
	lessons, err := repo.FetchByTeacherID(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, lessons)
}

func FetchCurrentProfessorStudents(w http.ResponseWriter, r *http.Request) {
	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewStudentsRepository(db)
	students, err := repo.FetchByTeacherID(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, students)
}
