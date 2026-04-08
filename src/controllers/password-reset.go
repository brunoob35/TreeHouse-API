package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/persistency"
	repositories "github.com/brunoob35/TreeHouse-API/src/repository"
	"github.com/brunoob35/TreeHouse-API/src/responses"
	"github.com/brunoob35/TreeHouse-API/src/security"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("RESET token recebido: [%q]", payload.Token)
	log.Printf("RESET nova senha vazia? %t", payload.NovaSenha == "")

	db, err := persistency.Connect()
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	resetsRepo := repositories.NewPasswordResetsRepository(db)
	usersRepo := repositories.NewUsersRepository(db)

	tokenHash := security.HashToken(payload.Token)

	log.Printf("RESET token hash calculado: [%q]", tokenHash)

	resetData, err := resetsRepo.FetchValidByTokenHash(tokenHash)
	if err != nil {
		log.Printf("RESET falha ao recuperar token hash=[%q]: %v", tokenHash, err)
		responses.Err(w, http.StatusUnauthorized, errors.New("token invalido ou expirado"))
		return
	}

	log.Printf(
		"RESET token encontrado id=%d user_id=%d expires_at=%s used_at=%v",
		resetData.ID,
		resetData.UserID,
		resetData.ExpiresAt.Format(time.RFC3339),
		resetData.UsedAt,
	)

	passwordHash, err := security.Hash(payload.NovaSenha)
	if err != nil {
		log.Printf("RESET erro ao gerar hash de senha para user_id=%d: %v", resetData.UserID, err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	err = usersRepo.UpdatePassword(resetData.UserID, string(passwordHash))
	if err != nil {
		log.Printf("RESET erro ao atualizar senha do user_id=%d: %v", resetData.UserID, err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("RESET senha atualizada com sucesso para user_id=%d", resetData.UserID)

	err = resetsRepo.MarkAsUsed(resetData.ID)
	if err != nil {
		log.Printf("RESET erro ao marcar token id=%d como usado: %v", resetData.ID, err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("RESET token id=%d marcado como usado", resetData.ID)

	responses.JSON(w, http.StatusOK, map[string]string{
		"message": "Senha redefinida com sucesso.",
	})
}
