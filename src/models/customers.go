package models

import (
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/brunoob35/TreeHouse-API/src/utils"
)

type Customer struct {
	ID             uint64     `json:"id,omitempty"`
	Nome           string     `json:"nome"`
	CPF            string     `json:"cpf"`
	Email          string     `json:"email,omitempty"`
	Telefone       string     `json:"telefone"`
	RG             string     `json:"rg,omitempty"`
	Nascimento     *time.Time `json:"nascimento,omitempty"`
	LGPDAceito     bool       `json:"lgpd_aceito"`
	LGPDAceitoEm   *time.Time `json:"lgpd_aceito_em,omitempty"`
	LGPDFinalidade string     `json:"lgpd_finalidade,omitempty"`
	Enderecos      []Address  `json:"enderecos,omitempty"`
	Ativo          bool       `json:"ativo"`
	StudentsCount  uint64     `json:"students_count,omitempty"`
	ContractsCount uint64     `json:"contracts_count,omitempty"`
	Students       []Student  `json:"students,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty"`
}

func (customer *Customer) Prepare(step string) error {
	customer.format()

	filteredStudents := make([]Student, 0, len(customer.Students))
	for _, student := range customer.Students {
		student.Nome = strings.TrimSpace(student.Nome)
		if student.Nome == "" {
			continue
		}

		student.Ativo = true
		if err := student.Prepare("create"); err != nil {
			return err
		}

		filteredStudents = append(filteredStudents, student)
	}

	customer.Students = filteredStudents

	filteredAddresses := make([]Address, 0, len(customer.Enderecos))
	for _, address := range customer.Enderecos {
		address.CEP = keepOnlyDigits(strings.TrimSpace(address.CEP))
		address.Rua = strings.TrimSpace(address.Rua)
		address.Numero = strings.TrimSpace(address.Numero)
		address.Bairro = strings.TrimSpace(address.Bairro)
		address.Cidade = strings.TrimSpace(address.Cidade)
		address.Estado = strings.TrimSpace(address.Estado)
		address.Pais = strings.TrimSpace(address.Pais)
		address.Complemento = strings.TrimSpace(address.Complemento)

		if address.Rua == "" && address.Numero == "" && address.Bairro == "" &&
			address.Cidade == "" && address.Estado == "" && address.CEP == "" {
			continue
		}

		if address.Pais == "" {
			address.Pais = "Brasil"
		}

		filteredAddresses = append(filteredAddresses, address)
	}

	customer.Enderecos = filteredAddresses
	return customer.validate(step)
}

func (customer *Customer) validate(step string) error {
	if customer.Nome == "" {
		return errors.New("o nome é obrigatório")
	}

	if customer.CPF == "" {
		return errors.New("o cpf é obrigatório")
	}

	if err := utils.CPFValidator(customer.CPF); err != nil {
		return err
	}

	if customer.Telefone == "" {
		return errors.New("o telefone é obrigatório")
	}

	if len(customer.Enderecos) == 0 {
		return errors.New("ao menos um endereço é obrigatório")
	}

	for _, address := range customer.Enderecos {
		if address.CEP == "" {
			return errors.New("o cep é obrigatório para salvar o endereço")
		}

		if address.Rua == "" {
			return errors.New("a rua é obrigatória para salvar o endereço")
		}

		if address.Numero == "" {
			return errors.New("o número é obrigatório para salvar o endereço")
		}

		if address.Bairro == "" {
			return errors.New("o bairro é obrigatório para salvar o endereço")
		}

		if address.Cidade == "" {
			return errors.New("a cidade é obrigatória para salvar o endereço")
		}

		if address.Estado == "" {
			return errors.New("o estado é obrigatório para salvar o endereço")
		}
	}

	return nil
}

func (customer *Customer) format() {
	customer.Nome = strings.TrimSpace(customer.Nome)
	customer.CPF = keepOnlyDigits(strings.TrimSpace(customer.CPF))
	customer.Email = strings.TrimSpace(customer.Email)
	customer.Telefone = strings.TrimSpace(customer.Telefone)
	customer.RG = keepAlphaNumeric(strings.TrimSpace(customer.RG))
	customer.LGPDFinalidade = strings.TrimSpace(customer.LGPDFinalidade)
}

func keepOnlyDigits(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, value)
}

func keepAlphaNumeric(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			return r
		}
		return -1
	}, value)
}
