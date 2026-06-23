package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
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

type requestLessonRescheduleRequest struct {
	RequestedLessonDate time.Time `json:"requested_lesson_date"`
}

func FetchLessonStatuses(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
	statuses, err := repository.FetchStatuses()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, statuses)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	repository := repositories.NewLessonsRepository(db).WithAuditUser(userID)
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

	permissions, err := authentication.ExtractPermissions(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	if authentication.HasPermission(permissions, authentication.PermProfessor) {
		userID, authErr := authentication.ExtractUserID(r)
		if authErr != nil {
			responses.Err(w, http.StatusUnauthorized, authErr)
			return
		}

		allowedStatus := request.StatusID == 2 || request.StatusID == 5
		if !allowedStatus {
			responses.Err(w, http.StatusForbidden, errors.New("professor pode apenas marcar presença ou sinalizar reagendamento"))
			return
		}

		belongsToTeacher, belongsErr := repository.BelongsToTeacher(lessonID, userID)
		if belongsErr != nil {
			responses.Err(w, http.StatusInternalServerError, belongsErr)
			return
		}

		if !belongsToTeacher {
			responses.Err(w, http.StatusForbidden, errors.New("aula não vinculada ao professor autenticado"))
			return
		}
	}

	if err = repository.UpdateStatus(lessonID, request.StatusID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func RequestLessonReschedule(w http.ResponseWriter, r *http.Request) {
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

	var request requestLessonRescheduleRequest
	if err = json.Unmarshal(bodyRequest, &request); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if request.RequestedLessonDate.IsZero() {
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
	if err = repository.RequestReschedule(lessonID, request.RequestedLessonDate); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	updatedLesson, fetchErr := repository.FetchByID(lessonID)
	if fetchErr != nil {
		responses.JSON(w, http.StatusOK, map[string]any{"id": lessonID})
		return
	}

	responses.JSON(w, http.StatusOK, updatedLesson)
}

func ApproveLessonReschedule(w http.ResponseWriter, r *http.Request) {
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
	if err = repository.ApproveReschedule(lessonID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	updatedLesson, fetchErr := repository.FetchByID(lessonID)
	if fetchErr != nil {
		responses.JSON(w, http.StatusOK, map[string]any{"id": lessonID})
		return
	}

	responses.JSON(w, http.StatusOK, updatedLesson)
}

func RejectLessonReschedule(w http.ResponseWriter, r *http.Request) {
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
	if err = repository.RejectReschedule(lessonID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	updatedLesson, fetchErr := repository.FetchByID(lessonID)
	if fetchErr != nil {
		responses.JSON(w, http.StatusOK, map[string]any{"id": lessonID})
		return
	}

	responses.JSON(w, http.StatusOK, updatedLesson)
}
