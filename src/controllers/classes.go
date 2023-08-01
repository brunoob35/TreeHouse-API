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
)

func CreateClass(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var newClass models.Classes
	if err := json.Unmarshal(bodyRequest, &newClass); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = newClass.Prepare(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	repo := repository.ClassesNewRepo(db)
	newClass.ID, err = repo.Create(newClass)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	responses.JSON(w, http.StatusCreated, newClass)
}

// FetchClassByID fetches a class from the persistency by ID
func FetchClassByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classID, err := strconv.Atoi(params["id"])
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

	repo := repository.ClassesNewRepo(db)
	class, err := repo.FetchByID(classID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, class)
}

// FetchAllClasses fetches all classes from the persistency
func FetchAllClasses(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	repo := repository.ClassesNewRepo(db)
	classes, err := repo.FetchAll()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	for _, class := range classes {
		class.Students, err = repo.SelectClassStudents(class)
	}

	defer db.Close()

	responses.JSON(w, http.StatusOK, classes)
}

func FetchActiveClasses(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	repo := repository.ClassesNewRepo(db)
	classes, err := repo.FetchAllActive()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	for _, class := range classes {
		class.Students, err = repo.SelectClassStudents(class)
	}

	defer db.Close()

	responses.JSON(w, http.StatusOK, classes)
}

// UpdateClass updates an existing class in the persistency
func UpdateClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classID, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Read the request body
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Unmarshal the request body into a class struct
	var updatedClass models.Classes
	if err := json.Unmarshal(bodyRequest, &updatedClass); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Perform any necessary preparations on the updated class
	// (e.g., validation, formatting, etc.)

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.ClassesNewRepo(db)
	updatedClass.ID = uint64(classID)
	err = repo.Update(updatedClass)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, updatedClass)
}

func SetStudentToClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classID, err := strconv.Atoi(params["classID"])
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var student models.Students
	if err := json.Unmarshal(bodyRequest, &student); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	repo := repository.ClassesNewRepo(db)
	classesUpdated, err := repo.SetStudent(classID, student)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	responses.JSON(w, http.StatusCreated, classesUpdated)
}

func RemoveStudentFromClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classID, err := strconv.Atoi(params["classID"])
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var student models.Students
	if err := json.Unmarshal(bodyRequest, &student); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	repo := repository.ClassesNewRepo(db)
	classesUpdated, err := repo.RemoveStudent(classID, student)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	defer db.Close()

	responses.JSON(w, http.StatusCreated, classesUpdated)
}
