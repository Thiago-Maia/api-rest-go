package models

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name string `json:"name"`
	Cpf  string `json:"cpf"`
	Rg   string `json:"rg"`
}

func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, 100)),

		validation.Field(&s.Cpf,
			validation.Required,
			validation.Length(11, 11),
			validation.By(ValidateCpf(s.Cpf)),
			validation.Match(regexp.MustCompile("^[0-9]{11}$"))),

		validation.Field(&s.Rg,
			validation.Required,
			validation.Length(9, 9),
			validation.Match(regexp.MustCompile("^[0-9]{9}$"))))
}

func ValidateCpf(cpf string) validation.RuleFunc {

	return func(value interface{}) error {
		cpf, _ := value.(string)
		if ok, _ := regexp.Match(`\d{11}`, []byte(cpf)); ok {
			digits := strings.Split(cpf, "")
			for multiplicator := 10.0; multiplicator < 12; multiplicator++ {
				result := CalculateDigit(digits, multiplicator)
				digitToValidate, err := strconv.Atoi(digits[int(multiplicator)-1])

				if err != nil {
					return err
				}

				if digitToValidate != result {
					return errors.New("CPF inválido")
				}
			}

			return nil
		}

		return errors.New("CPF inválido")
	}
}

func CalculateDigit(digits []string, multiplicator float64) int {
	var acumulator float64
	for _, v := range digits {
		digit, _ := strconv.ParseFloat(v, 64)
		acumulator += digit * multiplicator
		multiplicator -= 1

		if multiplicator == 1 {
			break
		}
	}

	if int(math.Mod(acumulator, 11.0)) <= 1 {
		return 0
	}

	return int(11 - math.Mod(acumulator, 11.0))
}
