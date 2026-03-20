package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/gorilla/mux"
)

type addStudentToLessonRequest struct {
	Note string `json:"note,omitempty"`
}

type updateLessonStatusRequest struct {
	StatusID uint64 `json:"status_id"`
}

func CreateLesson(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var lesson models.Lesson
	if err = json.Unmarshal(bodyRequest, &lesson); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = lesson.Prepare(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewLessonsRepository(db)
	lessonID, err := repository.Create(lesson)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	createdLesson, err := repository.FetchByID(lessonID)
	if err != nil {
		responses.JSON(w, http.StatusCreated, map[string]interface{}{
			"id": lessonID,
		})
		return
	}

	responses.JSON(w, http.StatusCreated, createdLesson)
}

func FetchLesson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
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

	repository := repositories.NewLessonsRepository(db)
	lesson, err := repository.FetchByID(lessonID)
	if err != nil {
		responses.Err(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, lesson)
}

func FetchAllLessons(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewLessonsRepository(db)
	lessons, err := repository.FetchAll()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, lessons)
}

func FetchLessonsByClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
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

	repository := repositories.NewLessonsRepository(db)
	lessons, err := repository.FetchByClass(classID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, lessons)
}

func UpdateLesson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var lesson models.Lesson
	if err = json.Unmarshal(bodyRequest, &lesson); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = lesson.PrepareUpdate(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewLessonsRepository(db)
	if err = repository.Update(lessonID, lesson); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func DeleteLesson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
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

	repository := repositories.NewLessonsRepository(db)
	if err = repository.Delete(lessonID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FetchLessonStudents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
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

	repository := repositories.NewLessonsRepository(db)
	students, err := repository.FetchStudents(lessonID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, students)
}

func AddStudentToLesson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	studentID, err := strconv.ParseUint(params["studentID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	var request addStudentToLessonRequest
	bodyRequest, readErr := io.ReadAll(r.Body)
	if readErr == nil && len(bodyRequest) > 0 {
		_ = json.Unmarshal(bodyRequest, &request)
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewLessonsRepository(db)
	if err = repository.AddStudent(lessonID, studentID, request.Note); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"lesson_id":  lessonID,
		"student_id": studentID,
		"origin":     "manual",
	})
}

func RemoveStudentFromLesson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

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

	repository := repositories.NewLessonsRepository(db)
	if err = repository.RemoveStudent(lessonID, studentID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func UpdateLessonStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lessonID, err := strconv.ParseUint(params["lessonID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var request updateLessonStatusRequest
	if err = json.Unmarshal(bodyRequest, &request); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if request.StatusID == 0 {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewLessonsRepository(db)
	if err = repository.UpdateStatus(lessonID, request.StatusID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}
