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

// CreateLessons inserts a new lesson into the persistency
func CreateLesson(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Unmarshal the request body into a lesson struct
	var newLesson models.Lessons
	if err := json.Unmarshal(bodyRequest, &newLesson); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Perform any necessary preparations on the new lesson
	// (e.g., validation, formatting, etc.)

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.LessonsNewRepo(db)
	newLesson.ID, err = repo.Create(newLesson)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newLesson)
}

func CreateLessonWStudent(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Unmarshal the request body into a lesson struct
	var newLesson models.Lessons
	if err := json.Unmarshal(bodyRequest, &newLesson); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.LessonsNewRepo(db)
	newLesson.ID, err = repo.Create(newLesson)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	students, err := GetClassStudents(w, newLesson.ClassID)

	for _, student := range students {
		repo.SetStudentLesson(newLesson, student)
	}

	responses.JSON(w, http.StatusCreated, newLesson)
}

// FetchLessonByID fetches a lesson from the persistency by ID
func FetchLessonByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.Atoi(params["id"])
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

	repo := repository.LessonsNewRepo(db)
	lesson, err := repo.FetchByID(lessonID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, lesson)
}

// FetchAllLessons fetches all lessons from the persistency
func FetchAllLessons(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.LessonsNewRepo(db)
	lessons, err := repo.FetchAll()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, lessons)
}

// UpdateLesson updates an existing lesson in the persistency
func UpdateLesson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.Atoi(params["id"])
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

	// Unmarshal the request body into a lesson struct
	var updatedLesson models.Lessons
	if err := json.Unmarshal(bodyRequest, &updatedLesson); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Perform any necessary preparations on the updated lesson
	// (e.g., validation, formatting, etc.)

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.LessonsNewRepo(db)
	updatedLesson.ID = lessonID
	err = repo.Update(updatedLesson)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, updatedLesson)
}
