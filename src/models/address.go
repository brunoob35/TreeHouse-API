package models

import "time"

type Address struct {
	ID          uint64    `json:"id,omitempty"`
	Rua         string    `json:"rua"`
	Numero      string    `json:"numero"`
	Bairro      string    `json:"bairro"`
	Cidade      string    `json:"cidade"`
	Estado      string    `json:"estado"`
	Pais        string    `json:"pais,omitempty"`
	Complemento string    `json:"complemento,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
