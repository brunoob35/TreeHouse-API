package controllers

import (
	"encoding/json"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

// CreateStudent inserts a new student into the persistency
func CreateStudent(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var newStudent models.Students
	if err := json.Unmarshal(bodyRequest, &newStudent); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Perform any necessary preparations on the new student
	// (e.g., validation, formatting, etc.)

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewStudentsRepository(db)

	newStudent, err = repo.Create(newStudent)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newStudent)
}

// FetchStudentByID fetches a student from the DB by ID
func FetchStudentByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentID, err := strconv.Atoi(params["id"])
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

	repo := repository.NewStudentsRepository(db)
	student, err := repo.FetchByID(studentID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, student)
}

// FetchAllStudents fetches all students from the DB
func FetchAllStudents(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewStudentsRepository(db)
	students, err := repo.FetchAll()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, students)
}

// UpdateStudent updates an existing student in the persistency
func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	studentID, err := strconv.Atoi(params["id"])
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var updatedStudent models.Students
	if err := json.Unmarshal(bodyRequest, &updatedStudent); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Perform any necessary preparations on the updated student
	// (e.g., validation, formatting, etc.)

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.NewStudentsRepository(db)
	updatedStudent.ID = int64(studentID)
	err = repo.Update(updatedStudent)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, updatedStudent)
}

//In the above code, we have implemented the `Student` struct controller with the following functions:
//
//1. `CreateStudent`: It creates a new student by reading the request body, unmarshaling it into a `Student` struct, performing any necessary preparations, opening a DB connection, creating a new instance of the student repository, and calling the `Create` method to insert the student into the persistency.
//
//2. `FetchStudentByID`: It fetches a student from the persistency by ID. It reads the student ID from the request parameters, opens a DB connection, creates a new instance of the student repository, and calls the `FetchByID` method to retrieve the student.
//
//3. `FetchAllStudents`: It fetches all students from the persistency. It opens a DB connection, creates a new instance of the student repository, and calls the `FetchAll` method to retrieve all students.
//
//4. `UpdateStudent`: It updates an existing student in the persistency. It reads the student ID from the request parameters, reads the request body, unmarshals it into a `Student` struct, performs any necessary preparations, opens a DB connection, creates a new instance of the student repository, and calls the `Update` method to update the student.
//
//Note: Replace `"github.com/your-package"` with the actual import path of your project.
