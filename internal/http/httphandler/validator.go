package httphandler

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate           = validator.New()
	validateTranslator = getValidatorTranslator()
)

func init() {
	_ = en_translations.RegisterDefaultTranslations(validate, validateTranslator)
}

func getValidatorTranslator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)

	validateTranslator, ok := uni.GetTranslator("en")
	if !ok {
		panic(errors.New("validator translator not found"))
	}

	return validateTranslator
}
