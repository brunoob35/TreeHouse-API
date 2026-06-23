package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/gorilla/mux"
)

func FetchContractTypes(w http.ResponseWriter, r *http.Request) {
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

	repo := repositories.NewContractsRepository(db).WithAuditUser(userID)
	types, err := repo.FetchTypes()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, types)
}

func FetchContractStatuses(w http.ResponseWriter, r *http.Request) {
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

	repo := repositories.NewContractsRepository(db).WithAuditUser(userID)
	statuses, err := repo.FetchStatuses()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, statuses)
}

func FetchContracts(w http.ResponseWriter, r *http.Request) {
	search := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	statusFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("status")))
	groupFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("group")))

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewContractsRepository(db)
	contracts, err := repo.FetchAll(search)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if statusFilter != "" {
		filtered := make([]models.Contract, 0, len(contracts))
		for _, contract := range contracts {
			effectiveStatus := strings.ToLower(strings.TrimSpace(contract.EffectiveStatusName))
			if effectiveStatus == strings.ToLower(statusFilter) {
				filtered = append(filtered, contract)
			}
		}
		contracts = filtered
	}

	if groupFilter != "" {
		filtered := make([]models.Contract, 0, len(contracts))
		for _, contract := range contracts {
			isInactive := contract.StatusID == models.ContractStatusInactive

			switch groupFilter {
			case "main":
				if !isInactive {
					filtered = append(filtered, contract)
				}
			case "inactive":
				if isInactive {
					filtered = append(filtered, contract)
				}
			default:
				filtered = append(filtered, contract)
			}
		}
		contracts = filtered
	}

	responses.JSON(w, http.StatusOK, contracts)
}

func FetchContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contractID, err := strconv.ParseUint(params["contractID"], 10, 64)
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

	repo := repositories.NewContractsRepository(db)
	contract, err := repo.FetchByID(contractID)
	if err != nil {
		responses.Err(w, http.StatusNotFound, err)
		return
	}

	responses.JSON(w, http.StatusOK, contract)
}

func CreateContract(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var contract models.Contract
	if err = json.Unmarshal(bodyRequest, &contract); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = contract.Prepare("create"); err != nil {
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

	repo := repositories.NewContractsRepository(db).WithAuditUser(userID)
	contractID, err := repo.Insert(contract)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	createdContract, err := repo.FetchByID(contractID)
	if err != nil {
		contract.ID = contractID
		responses.JSON(w, http.StatusCreated, contract)
		return
	}

	responses.JSON(w, http.StatusCreated, createdContract)
}

func UpdateContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contractID, err := strconv.ParseUint(params["contractID"], 10, 64)
	if err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var contract models.Contract
	if err = json.Unmarshal(bodyRequest, &contract); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if err = contract.Prepare("update"); err != nil {
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

	repo := repositories.NewContractsRepository(db).WithAuditUser(userID)
	if err = repo.Update(contractID, contract); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	updatedContract, err := repo.FetchByID(contractID)
	if err != nil {
		contract.ID = contractID
		responses.JSON(w, http.StatusOK, contract)
		return
	}

	responses.JSON(w, http.StatusOK, updatedContract)
}

func CreateClassFromContract(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	contractID, err := strconv.ParseUint(params["contractID"], 10, 64)
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

	userID, err := authentication.ExtractUserID(r)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	contractsRepo := repositories.NewContractsRepository(db).WithAuditUser(userID)
	classID, generatedLessonsCount, err := contractsRepo.CreateClassFromContract(contractID, class)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	classesRepo := repositories.NewClassesRepository(db)
	createdClass, err := classesRepo.FetchByID(classID)
	if err != nil {
		responses.JSON(w, http.StatusCreated, map[string]interface{}{
			"contract_id":             contractID,
			"class_id":                classID,
			"generated_lessons_count": generatedLessonsCount,
		})
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"contract_id":             contractID,
		"class_id":                classID,
		"class":                   createdClass,
		"generated_lessons_count": generatedLessonsCount,
	})
}
