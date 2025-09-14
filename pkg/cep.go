package pkg

import (
	"regexp"
)

func ValidateCepFormat(cep string) bool {
	cepRegex := regexp.MustCompile(`^\d{5}-\d{3}$`)
	cepWithHyphenRegex := regexp.MustCompile(`^\d{8}$`)

	matchCep := cepRegex.MatchString(cep)
	matchCepWithHyphen := cepWithHyphenRegex.MatchString(cep)

	return matchCep || matchCepWithHyphen
}
