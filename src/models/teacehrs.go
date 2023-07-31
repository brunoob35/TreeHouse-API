package models

// Teachers have their own ID and their UserID representing their status as system user and their status as a Teacher
type Teachers struct {
	ID     int64  `json:"id_teacehr"`
	UserID int64  `json:"user_id"`
	Name   string `json:"nome_usuario,omitempty"`
}
