package repositories

import "database/sql"

func setAuditUserTx(tx *sql.Tx, userID *uint64) error {
	if userID == nil || *userID == 0 {
		_, err := tx.Exec(`SET @app_usuario_id = NULL`)
		return err
	}

	_, err := tx.Exec(`SET @app_usuario_id = ?`, *userID)
	return err
}
