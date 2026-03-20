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

// CreateStudent is responsible for creating a new student.
//   - reads the request body
//   - parses the incoming JSON into a student struct
//   - prepares and validates the student data
//   - inserts the student record into the database
//   - all new students start as active
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var newStudent models.Student
	if err = json.Unmarshal(bodyRequest, &newStudent); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = newStudent.Prepare("create"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// All new students start as active
	newStudent.Ativo = true

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewStudentsRepository(db)

	newStudent.ID, err = repo.Insert(newStudent)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newStudent)
}

// FetchStudents returns a list of students optionally filtered by name.
// The optional "nome" query parameter performs a case-insensitive
// search using a LIKE clause in the database query.
func FetchStudents(w http.ResponseWriter, r *http.Request) {
	nome := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("nome")))

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewStudentsRepository(db)

	students, err := repo.FetchAll(nome)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, students)
}

// FetchStudent returns a single student by its ID.
// The student ID must be provided as a path parameter.
// Example: /students/{studentID}
func FetchStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	studentID, err := strconv.ParseUint(params["studentID"], 10, 64)
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

	repo := repositories.NewStudentsRepository(db)

	student, err := repo.FetchByID(studentID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, student)
}

// UpdateStudent updates the base data of an existing student.
// This endpoint updates the editable fields stored in the "alunos" table.
// The student ID is provided as a path parameter.
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	studentID, err := strconv.ParseUint(params["studentID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var student models.Student
	if err = json.Unmarshal(bodyRequest, &student); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = student.Prepare("update"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewStudentsRepository(db)

	if err = repo.Update(studentID, student); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteStudent performs a soft delete on a student.
// this operation sets the "ativo" field to false.
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	studentID, err := strconv.ParseUint(params["studentID"], 10, 64)
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

	repo := repositories.NewStudentsRepository(db)

	if err = repo.SoftDelete(studentID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
