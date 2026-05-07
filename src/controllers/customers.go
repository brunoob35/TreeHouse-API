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

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var customer models.Customer
	if err = json.Unmarshal(bodyRequest, &customer); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	customer.Ativo = true
	if err = customer.Prepare("create"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewCustomersRepository(db)
	customerID, err := repo.Insert(customer)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	createdCustomer, err := repo.FetchByID(customerID)
	if err != nil {
		customer.ID = customerID
		responses.JSON(w, http.StatusCreated, customer)
		return
	}

	createdStudents, err := repo.FetchStudents(customerID)
	if err == nil {
		createdCustomer.Students = createdStudents
	}
	createdAddresses, err := repo.FetchAddresses(customerID)
	if err == nil {
		createdCustomer.Enderecos = createdAddresses
	}

	responses.JSON(w, http.StatusCreated, createdCustomer)
}

func FetchCustomers(w http.ResponseWriter, r *http.Request) {
	search := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewCustomersRepository(db)
	customers, err := repo.FetchAll(search)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, customers)
}

func FetchCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID, err := strconv.ParseUint(params["customerID"], 10, 64)
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

	repo := repositories.NewCustomersRepository(db)
	customer, err := repo.FetchByID(customerID)
	if err != nil {
		responses.Err(w, http.StatusNotFound, err)
		return
	}

	students, err := repo.FetchStudents(customerID)
	if err == nil {
		customer.Students = students
	}
	addresses, err := repo.FetchAddresses(customerID)
	if err == nil {
		customer.Enderecos = addresses
	}

	responses.JSON(w, http.StatusOK, customer)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID, err := strconv.ParseUint(params["customerID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var customer models.Customer
	if err = json.Unmarshal(bodyRequest, &customer); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = customer.Prepare("update"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewCustomersRepository(db)
	if err = repo.Update(customerID, customer); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	if err = repo.ReplaceAddresses(customerID, customer.Enderecos); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	if err = repo.ReplaceStudents(customerID, customer.Students); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	updatedCustomer, err := repo.FetchByID(customerID)
	if err != nil {
		responses.JSON(w, http.StatusOK, customer)
		return
	}

	students, err := repo.FetchStudents(customerID)
	if err == nil {
		updatedCustomer.Students = students
	}
	addresses, err := repo.FetchAddresses(customerID)
	if err == nil {
		updatedCustomer.Enderecos = addresses
	}

	responses.JSON(w, http.StatusOK, updatedCustomer)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID, err := strconv.ParseUint(params["customerID"], 10, 64)
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

	repo := repositories.NewCustomersRepository(db)
	if err = repo.SoftDelete(customerID); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

func FetchCustomerStudents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID, err := strconv.ParseUint(params["customerID"], 10, 64)
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

	repo := repositories.NewCustomersRepository(db)
	students, err := repo.FetchStudents(customerID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, students)
}

func FetchCustomerAddresses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	customerID, err := strconv.ParseUint(params["customerID"], 10, 64)
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

	repo := repositories.NewCustomersRepository(db)
	addresses, err := repo.FetchAddresses(customerID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, addresses)
}
