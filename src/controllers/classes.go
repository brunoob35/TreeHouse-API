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
	classID, err := repository.Create(class)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	class.ID = classID
	responses.JSON(w, http.StatusCreated, class)
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

	responses.JSON(w, http.StatusNoContent, nil)
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
	classID, err := repository.CreatePrivateClassFromStudent(request.StudentID, request.TeacherID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	class, err := repository.FetchByID(classID)
	if err != nil {
		responses.JSON(w, http.StatusCreated, map[string]interface{}{
			"id": classID,
		})
		return
	}

	responses.JSON(w, http.StatusCreated, class)
}
