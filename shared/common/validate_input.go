package common

import (
  "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
  uni       *ut.UniversalTranslator
  validate  *validator.Validate
)

type Map map[string]string

func ValidateInput(model any) Map {
  en := en.New()
  uni = ut.New(en, en)
  
  trans, _ := uni.GetTranslator("en")
  
  validate = validator.New()
  en_translations.RegisterDefaultTranslations(validate, trans)
  
  errMessage := Map{}
  if err := validate.Struct(model); err != nil {
    for _, errs := range err.(validator.ValidationErrors) {
      errMessage[errs.StructField()] = errs.Translate(trans)
    }
  }
  
  return errMessage
}