package repositories

import (
	"database/sql"
	"errors"

	"github.com/brunoob35/TreeHouse-API/src/models"
)

// PasswordResetsRepository is responsible for password reset operations.
type PasswordResetsRepository struct {
	db *sql.DB
}

// NewPasswordResetsRepository creates a new repository instance.
func NewPasswordResetsRepository(db *sql.DB) *PasswordResetsRepository {
	return &PasswordResetsRepository{db: db}
}

// Create inserts a new password reset token.
func (r *PasswordResetsRepository) Create(reset models.PasswordReset) (uint64, error) {
	query := `
		INSERT INTO treehousedb.password_resets (
			user_id,
			token_hash,
			expires_at
		) VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		reset.UserID,
		reset.TokenHash,
		reset.ExpiresAt,
	)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(insertedID), nil
}

// FetchValidByTokenHash returns a valid (not used and not expired) reset token.
func (r *PasswordResetsRepository) FetchValidByTokenHash(tokenHash string) (models.PasswordReset, error) {
	query := `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			used_at,
			created_at
		FROM treehousedb.password_resets
		WHERE token_hash = ?
		  AND used_at IS NULL
		  AND expires_at > CURRENT_TIMESTAMP
		LIMIT 1
	`

	var reset models.PasswordReset

	err := r.db.QueryRow(query, tokenHash).Scan(
		&reset.ID,
		&reset.UserID,
		&reset.TokenHash,
		&reset.ExpiresAt,
		&reset.UsedAt,
		&reset.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PasswordReset{}, sql.ErrNoRows
		}
		return models.PasswordReset{}, err
	}

	return reset, nil
}

// MarkAsUsed marks a token as used.
func (r *PasswordResetsRepository) MarkAsUsed(id uint64) error {
	query := `
		UPDATE treehousedb.password_resets
		SET used_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("nenhum password reset encontrado")
	}

	return nil
}
