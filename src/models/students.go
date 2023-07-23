package models

type Students struct {
	ID            int64  `json:"id,omitempty"`
	Nome          string `json:"nome_usuario,omitempty"`
	Idade         uint32 `json:"idade,omitempty"`
	Livro         string `json:"livro,omitempty"`
	ProfessoraId  int32  `json:"professora,omitempty"`
	IdPais        int32  `json:"id_pais,omitempty"`
	IdAnoLetivo   int32  `json:"id_ano_letivo,omitempty"`
	Alfabetizacao int32  `json:"alfabetizacao,omitempty"`
	IdTurma       int32  `json:"id_turma,omitempty"`
}
