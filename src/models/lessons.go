package models

import "time"

// Lessons represent every class given. It can have multiple students and multiple teachers, but it belongs to a single Class
type Lessons struct {
	ID       uint64     `json:"id_aula,omitempty"`
	ClassID  uint64     `json:"id_turma"`
	Teachers []Teachers `json:"id_professor"`
	Students []Students `json:"alunos"`
	Status   uint64     `json:"id_status_aula,omitempty"`
	Date     time.Time  `json:"datahora_aula,omitempty"`
	End      time.Time  `json:"datahora_fim_aula,omitempty"`

	// Educational info:
	Adaptations string `json:"adaptacoes,omitempty"`
	Performance string `json:"desempenho,omitempty"`
	Balance     string `json:"saldo,omitempty"`
	Coments     string `json:"comentarios,omitempty"`
}
