package models

type ProfessorClassCountRequest struct {
	ProfessorIDs []uint64 `json:"professor_ids"`
}

type ProfessorClassCount struct {
	ProfessorID  uint64 `json:"professor_id"`
	ClassesCount uint64 `json:"classes_count"`
}
