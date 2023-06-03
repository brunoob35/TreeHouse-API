package utils

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
)

func CPFValidator(cpf string) error {
	if val := CheckAllEqual(cpf); val == true {
		if val = CalcularDv1(cpf); val == true {
			if val = CalcularDv2(cpf); val == true {

				return nil
			} else {
				return errors.New("DV2 invalido, CPF Invalido")
			}
			return nil
		} else {
			return errors.New("DV1 invalido, CPF Invalido")

		}
		return nil

	} else {
		return errors.New("Numeros Iguais, CPF Invalido")
	}
}

func CalcularDv1(cpf string) bool {
	soma := 0
	cpfSlice := strings.Split(cpf, "")
	for i := 0; i <= 8; i++ {
		indexCpf := cpfSlice[i]
		intIndexCpf, err := strconv.Atoi(indexCpf)
		if err != nil {
			// todo: Handle error
		} else {
			// todo: Use the converted integer value (num)
		}
		multiplicacao := intIndexCpf * (10 - i)
		soma += multiplicacao
	}
	resultado := float64(soma) * 10 / 11
	resultadoArredondadoCima := ((math.Ceil(float64(resultado) * 10)) / 10) * 10

	strNumVer := strconv.FormatFloat(resultadoArredondadoCima, 'f', 2, 64)
	digitoVerificador := (strings.Split(strNumVer, ""))[3]
	if digitoVerificador == cpfSlice[9] {
		return true
	}
	return false
}

func CalcularDv2(cpf string) bool {
	soma := 0
	cpfSlice := strings.Split(cpf, "")
	for i := 0; i <= 9; i++ {
		indexCpf := cpfSlice[i]
		intIndexCpf, err := strconv.Atoi(indexCpf)
		if err != nil {
			// todo: Handle error
		} else {
			// todo: Use the converted integer value (num)
		}
		multiplicacao := intIndexCpf * (11 - i)
		soma += multiplicacao
	}
	resultado := float64(soma) * 10 / 11
	resultadoArredondadoCima := ((math.Ceil(float64(resultado) * 10)) / 10) * 10

	strNumVer := strconv.FormatFloat(resultadoArredondadoCima, 'f', 2, 64)
	digitoVerificador := (strings.Split(strNumVer, ""))[3]
	if digitoVerificador == cpfSlice[10] {
		return true
	}
	return false
}

func CheckAllEqual(cpf string) bool {
	cpfSlice := strings.Split(cpf, "")
	FirstNumber := cpfSlice[0]

	for _, i := range cpfSlice {

		if i != FirstNumber {
			return true
		}
	}

	return false
}
