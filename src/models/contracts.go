package models

import (
	"errors"
	"strings"
	"time"

	"github.com/brunoob35/TreeHouse-API/src/utils"
)

const (
	ContractStatusActive          uint64 = 1
	ContractStatusPending         uint64 = 2
	ContractStatusExpired         uint64 = 3
	ContractStatusUpcomingDue     uint64 = 4
	ContractTypeAnnual            uint64 = 1
	ContractTypeSemiannual        uint64 = 2
	ContractTypeQuarterly         uint64 = 3
	ContractTypeMonthly           uint64 = 4
	ContractTypeTemporary         uint64 = 5
	ContractUpcomingDueWindowDays        = 7
)

type Contract struct {
	ID                      uint64     `json:"id,omitempty"`
	RepresentativeCustomerID uint64    `json:"id_cliente_representante,omitempty"`
	ResponsibleCustomerID   uint64     `json:"id_cliente_responsavel,omitempty"`
	StudentID               uint64     `json:"id_aluno,omitempty"`
	ContractTypeID          uint64     `json:"id_tipo_contrato"`
	StatusID                uint64     `json:"id_status,omitempty"`
	ClassID                 *uint64    `json:"id_turma,omitempty"`
	Value                   float64    `json:"valor"`
	RepresentativeEmail     string     `json:"email_representante,omitempty"`
	RepresentativeCPF       string     `json:"cpf_representante,omitempty"`
	RG                      string     `json:"rg,omitempty"`
	RepresentativePhone     string     `json:"telefone_representante,omitempty"`
	RepresentativeCivilStatus string   `json:"est_civil_representante,omitempty"`
	DiscountPercentage      float64    `json:"desconto_porcentagem,omitempty"`
	FinalValue              float64    `json:"valor_final,omitempty"`
	Installments            *uint64    `json:"parcelas,omitempty"`
	InstallmentsDescription string     `json:"parcelas_descricao,omitempty"`
	LessonsCount            *uint64    `json:"numero_aulas,omitempty"`
	Periodicity             string     `json:"periodicidade,omitempty"`
	LessonDuration          string     `json:"tempo_aula,omitempty"`
	ContractDuration        string     `json:"tempo_contrato,omitempty"`
	StartDate               *time.Time `json:"inicio_contrato,omitempty"`
	DueDate                 *time.Time `json:"vencimento_contrato,omitempty"`
	FirstLessonDate         *time.Time `json:"primeira_aula,omitempty"`
	CreatedAt               time.Time  `json:"created_at,omitempty"`
	UpdatedAt               time.Time  `json:"updated_at,omitempty"`

	StudentName             string     `json:"student_name,omitempty"`
	ResponsibleName         string     `json:"responsible_name,omitempty"`
	RepresentativeName      string     `json:"representative_name,omitempty"`
	ContractTypeName        string     `json:"contract_type_name,omitempty"`
	StatusName              string     `json:"status_name,omitempty"`
	EffectiveStatusID       uint64     `json:"effective_status_id,omitempty"`
	EffectiveStatusName     string     `json:"effective_status_name,omitempty"`

	ResponsibleCustomer     *Customer  `json:"responsavel,omitempty"`
	RepresentativeCustomer  *Customer  `json:"representante,omitempty"`
	Student                 *Student   `json:"aluno,omitempty"`
}

type ContractStatus struct {
	ID   uint64 `json:"id"`
	Name string `json:"nome_status"`
}

type ContractType struct {
	ID   uint64 `json:"id"`
	Name string `json:"nome_tipo"`
}

func (contract *Contract) Prepare(step string) error {
	contract.format()
	contract.applyDerivedValues()

	if err := contract.validate(step); err != nil {
		return err
	}

	return nil
}

func (contract *Contract) format() {
	contract.RepresentativeEmail = strings.TrimSpace(strings.ToLower(contract.RepresentativeEmail))
	contract.RepresentativeCPF = keepOnlyDigits(strings.TrimSpace(contract.RepresentativeCPF))
	contract.RG = keepAlphaNumeric(strings.TrimSpace(contract.RG))
	contract.RepresentativePhone = strings.TrimSpace(contract.RepresentativePhone)
	contract.RepresentativeCivilStatus = strings.TrimSpace(contract.RepresentativeCivilStatus)
	contract.InstallmentsDescription = strings.TrimSpace(contract.InstallmentsDescription)
	contract.Periodicity = strings.TrimSpace(contract.Periodicity)
	contract.LessonDuration = strings.TrimSpace(contract.LessonDuration)
	contract.ContractDuration = strings.TrimSpace(contract.ContractDuration)

	if contract.ResponsibleCustomer != nil {
		contract.ResponsibleCustomer.Ativo = true
		_ = contract.ResponsibleCustomer.Prepare("create")
	}

	if contract.RepresentativeCustomer != nil {
		contract.RepresentativeCustomer.Ativo = true
		_ = contract.RepresentativeCustomer.Prepare("create")
	}

	if contract.Student != nil {
		contract.Student.Ativo = true
		_ = contract.Student.Prepare("create")
	}
}

func (contract *Contract) applyDerivedValues() {
	if contract.StatusID == 0 {
		contract.StatusID = ContractStatusPending
	}

	if contract.FinalValue <= 0 && contract.Value > 0 {
		discountFactor := 1 - (contract.DiscountPercentage / 100)
		if discountFactor < 0 {
			discountFactor = 0
		}
		contract.FinalValue = contract.Value * discountFactor
	}

	if contract.DueDate == nil && contract.StartDate != nil {
		contract.DueDate = deriveContractDueDate(*contract.StartDate, contract.ContractTypeID)
	}
}

func (contract *Contract) validate(step string) error {
	if contract.ContractTypeID == 0 {
		return errors.New("id_tipo_contrato é obrigatório")
	}

	if contract.Value <= 0 {
		return errors.New("o valor do contrato deve ser maior que zero")
	}

	if contract.FinalValue <= 0 {
		return errors.New("o valor final do contrato deve ser maior que zero")
	}

	if contract.ResponsibleCustomerID == 0 && contract.ResponsibleCustomer == nil {
		return errors.New("um cliente responsável é obrigatório")
	}

	if contract.StudentID == 0 && contract.Student == nil {
		return errors.New("um aluno é obrigatório")
	}

	if contract.RepresentativeCPF != "" {
		if err := utils.CPFValidator(contract.RepresentativeCPF); err != nil {
			return err
		}
	}

	if contract.ResponsibleCustomer != nil {
		if err := contract.ResponsibleCustomer.Prepare("create"); err != nil {
			return err
		}
	}

	if contract.RepresentativeCustomer != nil {
		if err := contract.RepresentativeCustomer.Prepare("create"); err != nil {
			return err
		}
	}

	if contract.Student != nil {
		if err := contract.Student.Prepare("create"); err != nil {
			return err
		}
	}

	if step == "create" && contract.RepresentativeCustomerID == 0 && contract.RepresentativeCustomer == nil {
		contract.RepresentativeCustomerID = contract.ResponsibleCustomerID
	}

	return nil
}

func deriveContractDueDate(startDate time.Time, contractTypeID uint64) *time.Time {
	var dueDate time.Time

	switch contractTypeID {
	case ContractTypeAnnual:
		dueDate = startDate.AddDate(1, 0, 0)
	case ContractTypeSemiannual:
		dueDate = startDate.AddDate(0, 6, 0)
	case ContractTypeQuarterly:
		dueDate = startDate.AddDate(0, 3, 0)
	case ContractTypeMonthly:
		dueDate = startDate.AddDate(0, 1, 0)
	default:
		return nil
	}

	return &dueDate
}

func ComputeEffectiveContractStatus(contract Contract, now time.Time) (uint64, string) {
	if contract.StatusID == ContractStatusExpired {
		return ContractStatusExpired, "Vencido"
	}

	if contract.DueDate != nil {
		if contract.DueDate.Before(now) {
			return ContractStatusExpired, "Vencido"
		}

		if contract.StatusID == ContractStatusActive {
			windowEnd := now.AddDate(0, 0, ContractUpcomingDueWindowDays)
			if !contract.DueDate.After(windowEnd) {
				return ContractStatusUpcomingDue, "Prox. Vencimento"
			}
		}
	}

	switch contract.StatusID {
	case ContractStatusActive:
		return ContractStatusActive, "Ativo"
	case ContractStatusPending:
		return ContractStatusPending, "Pendente"
	case ContractStatusExpired:
		return ContractStatusExpired, "Vencido"
	default:
		if contract.StatusName != "" {
			return contract.StatusID, contract.StatusName
		}
		return contract.StatusPendingFallback()
	}
}

func (contract Contract) StatusPendingFallback() (uint64, string) {
	return ContractStatusPending, "Pendente"
}
