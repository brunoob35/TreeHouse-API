package utils

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
)

//func CPFValidator(cpf string) error {
//	if val := CheckAllEqual(cpf); val == true {
//		log.Println("passou no if 1")
//
//		if CalcularDv1(cpf) {
//			//if CalcularDv2(cpf) {
//			return nil
//		} else {
//			return errors.New("DV1 invalido, CPF Invalido")
//		}
//	} else {
//		return errors.New("Numeros Iguais, CPF Invalido")
//	}
//}

func CPFValidator(cpf string) error {
	if val := CheckAllEqual(cpf); val == true {
		log.Println("passou no if 1")
		if val = CalcularDv1(cpf); val == true {
			log.Println("entrou no if 2")
			if val = CalcularDv2(cpf); val == true {

				log.Println("entrou no if 3")
				return nil
			} else {
				log.Println("nao passou no calculardv2")
				return errors.New("DV2 invalido, CPF Invalido")
			}
			return nil
		} else {
			log.Println("nao passou no calculardv1")
			return errors.New("DV1 invalido, CPF Invalido")

		}
		return nil

	} else {
		log.Println("nao passou no check all equal")
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
	log.Println("debug dig verificador", digitoVerificador)
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
	log.Println("debug dig verificador", digitoVerificador)
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
