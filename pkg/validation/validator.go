package validation

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/guregu/null.v4"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {

	if err := v.validator.RegisterValidation("date_only", validateDateOnly); err != nil {
		log.Error("failed register validation date_only")
		return err
	}
	if err := v.validator.RegisterValidation("unique", validateUnique); err != nil {
		log.Error("failed register validation unique")
		return err
	}
	if err := v.validator.RegisterValidation("enum", validateEnum); err != nil {
		log.Error("failed register validation enum")
		return err
	}
	if err := v.validator.RegisterValidation("digit", validateOnlyNumber); err != nil {
		log.Error("failed register validation digit")
		return err
	}
	if err := v.validator.RegisterValidation("unique_update", validateUpdateUnique); err != nil {
		log.Error("failed register validation unique_update")
		return err
	}
	if err := v.validator.RegisterValidation("rfe", validateRequireIfAnotherField); err != nil {
		log.Error("failed register validation rfe")
		return err
	}
	v.validator.RegisterCustomTypeFunc(nullFloatValidator, null.Float{})
	v.validator.RegisterCustomTypeFunc(nullIntValidator, null.Int{})
	v.validator.RegisterCustomTypeFunc(nullTimeValidator, null.Time{})
	return v.validator.Struct(i)
}
