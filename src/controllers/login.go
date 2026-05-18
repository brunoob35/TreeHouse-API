package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/brunoob35/TreeHouse-API/src/authentication"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	"github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

type loginResponse struct {
	Token       string `json:"token,omitempty"`
	FirstAccess bool   `json:"first_access,omitempty"`
	Email       string `json:"email,omitempty"`
}

// Login is responsible for validating user credentials.
//
// This flow performs the following steps:
//   - reads the request body
//   - parses the incoming credentials
//   - loads the user by email
//   - validates the provided password
//   - loads the user's permission IDs from the database
//   - aggregates those permissions into a numeric bitmask
//   - generates a JWT token containing the user ID and permission mask
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

	repo := repositories.NewUsersRepository(db)

	// FetchByEmail loads the user base data required for authentication.
	userFound, err := repo.FetchByEmail(user.Email)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, errors.New("credenciais inválidas"))
		return
	}

	if userFound.Senha == "" {
		firstAccessToken, tokenErr := authentication.GenerateFirstAccessToken(userFound.ID)
		if tokenErr != nil {
			responses.Err(w, http.StatusInternalServerError, tokenErr)
			return
		}

		responses.JSON(w, http.StatusOK, loginResponse{
			Token:       firstAccessToken,
			FirstAccess: true,
			Email:       userFound.Email,
		})
		return
	}

	// ValidatePassword compares the stored password hash with the plain password
	// received in the login request.
	if err = security.ValidatePassword(userFound.Senha, user.Senha); err != nil {
		responses.Err(w, http.StatusUnauthorized, err)
		log.Println("Unauthorized")
		return
	}

	// FetchPermissionMaskByUserID loads the user's permission IDs from the
	// relationship table and aggregates them into a single numeric bitmask.
	permissionMask, err := repo.FetchPermissionMaskByUser(userFound.ID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	token, err := authentication.GenerateToken(userFound.ID, permissionMask)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, loginResponse{
		Token: token,
	})
}

func CompleteFirstAccess(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var payload struct {
		Token     string `json:"token"`
		NovaSenha string `json:"nova_senha"`
	}

	if err = json.Unmarshal(bodyRequest, &payload); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	if payload.Token == "" {
		responses.Err(w, http.StatusBadRequest, errors.New("token é obrigatório"))
		return
	}

	if len(payload.NovaSenha) < 6 {
		responses.Err(w, http.StatusBadRequest, errors.New("a nova senha deve ter pelo menos 6 caracteres"))
		return
	}

	userID, err := authentication.ExtractFirstAccessUserID(payload.Token)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, errors.New("token inválido ou expirado"))
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repo := repositories.NewUsersRepository(db)
	userFound, err := repo.FetchByID(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if userFound.Senha != "" {
		responses.Err(w, http.StatusConflict, errors.New("senha inicial já foi definida"))
		return
	}

	passwordHash, err := security.Hash(payload.NovaSenha)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	if err = repo.UpdatePassword(userID, string(passwordHash)); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{
		"message": "Primeira senha cadastrada com sucesso.",
	})
}
