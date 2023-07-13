package models

type Students struct {
	ID            uint64 `json:"id,omitempty"`
	Nome          string `json:"nome_usuario,omitempty"`
	Idade         uint32 `json:"idade,omitempty"`
	Livro         string `json:"livro,omitempty"`
	ProfessoraId  uint32 `json:"professora,omitempty"`
	IdPais        uint32 `json:"id_pais,omitempty"`
	IdAnoLetivo   uint32 `json:"id_ano_letivo,omitempty"`
	Alfabetizacao uint32 `json:"alfabetizacao,omitempty"`
	IdTurma       uint32 `json:"id_turma,omitempty"`
}
