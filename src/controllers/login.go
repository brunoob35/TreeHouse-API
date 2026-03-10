package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

// Login is responsible for validating user credentials
func Login(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repository.UsersNewRepo(db)
	userFound, err := repo.FetchByEmail(user.Email)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = security.ValidatePassword(userFound.Senha, user.Senha); err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		return
	}

	token, err := authentication.GenerateToken(userFound.ID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	//userID := strconv.FormatUint(userFound.ID, 10)
	//
	//responses.JSON(w, http.StatusOK, models.DadosAutenticacao{ID: usuarioID, Token: token})
	w.Write([]byte(token))
}
