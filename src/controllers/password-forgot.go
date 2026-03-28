package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/mailer"
	"github.com/brunoob35/TreeHouse-API/src/models"
	"github.com/brunoob35/TreeHouse-API/src/persistency"
	repositories "github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var payload struct {
		Email string `json:"email"`
	}

	if err = json.Unmarshal(bodyRequest, &payload); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	usersRepo := repositories.NewUsersRepository(db)
	resetsRepo := repositories.NewPasswordResetsRepository(db)

	userFound, err := usersRepo.FetchByEmail(payload.Email)

	// resposta generica mesmo em caso de email inexistente
	if err != nil {
		log.Printf("FORGOT usuario nao encontrado para email: [%s]", payload.Email)

		responses.JSON(w, http.StatusOK, map[string]string{
			"message": "Se o email existir, um link de recuperacao foi enviado.",
		})
		return
	}

	token, err := security.GenerateSecureToken(32)
	if err != nil {
		log.Printf("FORGOT erro ao gerar token para user_id=%d: %v", userFound.ID, err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	tokenHash := security.HashToken(token)

	reset := models.PasswordReset{
		UserID:    userFound.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(30 * time.Minute),
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)

	log.Printf("FORGOT user_id=%d email=[%s]", userFound.ID, userFound.Email)
	log.Printf("FORGOT token gerado: [%q]", token)
	log.Printf("FORGOT token hash salvo: [%q]", tokenHash)
	log.Printf("FORGOT expires_at: [%s]", reset.ExpiresAt)
	log.Printf("FORGOT reset link: [%q]", resetLink)

	resetID, err := resetsRepo.Create(reset)
	if err != nil {
		log.Printf("FORGOT erro ao salvar password reset para user_id=%d: %v", userFound.ID, err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("FORGOT password reset criado com id=%d para user_id=%d", resetID, userFound.ID)

	err = mailer.SendPasswordResetEmail(userFound.Email, resetLink)
	if err != nil {
		log.Printf("FORGOT erro ao enviar email para user_id=%d email=[%s]: %v", userFound.ID, userFound.Email, err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("FORGOT email de reset enviado para user_id=%d email=[%s]", userFound.ID, userFound.Email)

	responses.JSON(w, http.StatusOK, map[string]string{
		"message": "Se o email existir, um link de recuperacao foi enviado.",
	})
}
