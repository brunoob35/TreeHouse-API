package controllers

import (
	"encoding/json"
	"errors"
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

func CreateClass(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var class models.Class
	if err = json.Unmarshal(bodyRequest, &class); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = class.Prepare(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewClassesRepository(db)
	classID, generatedLessonsCount, err := repository.Create(class)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	fetchedClass, err := repository.FetchByID(classID)
	if err != nil {
		class.ID = classID
		responses.JSON(w, http.StatusCreated, map[string]interface{}{
			"id":                      class.ID,
			"teacher_id":              class.TeacherID,
			"name":                    class.Name,
			"recurrence_desc":         class.RecurrenceDesc,
			"recurrence_json":         class.RecurrenceJSON,
			"endereco":                class.Endereco,
			"generated_lessons_count": generatedLessonsCount,
		})
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":                      fetchedClass.ID,
		"teacher_id":              fetchedClass.TeacherID,
		"id_endereco":             fetchedClass.IDEndereco,
		"name":                    fetchedClass.Name,
		"recurrence_desc":         fetchedClass.RecurrenceDesc,
		"recurrence_json":         fetchedClass.RecurrenceJSON,
		"endereco":                fetchedClass.Endereco,
		"generated_lessons_count": generatedLessonsCount,
	})
}

func FetchClass(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewClassesRepository(db)
	class, err := repository.FetchByID(classID)
	if err != nil {
		responses.Err(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, class)
}

func FetchAllActiveClasses(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewClassesRepository(db)
	classes, err := repository.FetchAllActive()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, classes)
}

func FetchAllClasses(w http.ResponseWriter, r *http.Request) {
	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewClassesRepository(db)
	classes, err := repository.FetchAll()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, classes)
}

func UpdateClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	classID, err := strconv.ParseUint(params["classID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var class models.Class
	if err = json.Unmarshal(bodyRequest, &class); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = class.Prepare(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewClassesRepository(db)
	if err = repository.Update(classID, class); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	updatedClass, err := repository.FetchByID(classID)
	if err != nil {
		responses.JSON(w, http.StatusOK, class)
		return
	}

	responses.JSON(w, http.StatusOK, updatedClass)
}

func DeleteClass(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewClassesRepository(db)
	if err = repository.SoftDelete(classID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func AddStudentToClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	classID, err := strconv.ParseUint(params["classID"], 10, 64)
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

	repository := repositories.NewClassesRepository(db)
	if err = repository.AddStudent(classID, studentID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"class_id":   classID,
		"student_id": studentID,
	})
}

func RemoveStudentFromClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	classID, err := strconv.ParseUint(params["classID"], 10, 64)
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

	repository := repositories.NewClassesRepository(db)
	if err = repository.RemoveStudent(classID, studentID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FetchClassStudents(w http.ResponseWriter, r *http.Request) {
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

	repository := repositories.NewClassesRepository(db)
	students, err := repository.FetchStudents(classID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, students)
}

func CreatePrivateClassFromStudent(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var request models.CreatePrivateClassRequest
	if err = json.Unmarshal(bodyRequest, &request); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if request.StudentID == 0 {
		responses.Err(w, http.StatusBadRequest, errors.New("student_id is required"))
		return
	}

	class := models.Class{
		TeacherID:      request.TeacherID,
		Name:           request.Name,
		RecurrenceDesc: request.RecurrenceDesc,
		RecurrenceJSON: request.RecurrenceJSON,
		Endereco:       request.Endereco,
	}
	if strings.TrimSpace(class.Name) == "" {
		class.Name = "Turma"
	}

	if err = class.Prepare(); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repository := repositories.NewClassesRepository(db)
	classID, generatedLessonsCount, err := repository.CreatePrivateClassFromStudent(request.StudentID, class)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	fetchedClass, err := repository.FetchByID(classID)
	if err != nil {
		responses.JSON(w, http.StatusCreated, map[string]interface{}{
			"id":                      classID,
			"generated_lessons_count": generatedLessonsCount,
		})
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":                      fetchedClass.ID,
		"teacher_id":              fetchedClass.TeacherID,
		"id_endereco":             fetchedClass.IDEndereco,
		"name":                    fetchedClass.Name,
		"recurrence_desc":         fetchedClass.RecurrenceDesc,
		"recurrence_json":         fetchedClass.RecurrenceJSON,
		"endereco":                fetchedClass.Endereco,
		"generated_lessons_count": generatedLessonsCount,
	})
}
