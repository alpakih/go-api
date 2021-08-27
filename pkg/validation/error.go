package validation

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type ErrorValidation struct {
	Field     string `json:"field,omitempty"`
	ActualTag string `json:"tag,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
	Param     string `json:"param,omitempty"`
	Message   string `json:"message"`
}

func WrapValidationErrors(errs validator.ValidationErrors) []ErrorValidation {
	validationErrors := make([]ErrorValidation, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, ErrorValidation{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
			Message:   FormatMessage(validationErr),
		})
	}

	return validationErrors
}

func FormatMessage(err validator.FieldError) string {
	field := err.Field()
	param := err.Param()

	message := fmt.Sprintf("Field validation for '%s' failed on the '%s'", strcase.ToSnake(err.Field()), err.Tag())

	switch err.Tag() {
	case "required":
		message = fmt.Sprintf("The %s field is required.", field)
	case "number":
		message = fmt.Sprintf("The %s must be a number.", field)
	case "email":
		message = fmt.Sprintf("The %s must be a valid email address.", field)
	case "gt":
		message = fmt.Sprintf("The %s must be greater than %s.", field, param)
	case "gte":
		message = fmt.Sprintf("The %s must be greater than or equal %s.", field, param)
	case "lt":
		message = fmt.Sprintf("The %s must be less than %s.", field, param)
	case "lte":
		message = fmt.Sprintf("The %s must be less than or equal %s.", field, param)
	case "phone_number":
		message = fmt.Sprintf("The %s must be a valid phone number.", field)
	case "min":
		message = fmt.Sprintf("The %s must be at least %s", field, param)
	case "max":
		message = fmt.Sprintf("The %s may not be greater than %s.", field, param)
	case "len":
		message = fmt.Sprintf("The %s must be a length %s.", field, param)
	case "eq":
		message = fmt.Sprintf("The %s must be a equals %s.", field, param)
	case "date_only":
		message = fmt.Sprintf("The %s must be valid format %s.", field, "yyyy-mm-dd")
	case "unique":
		message = fmt.Sprintf("The %s %s is already exist.", strcase.ToSnake(field), err.Value())
	case "unique_update":
		message = fmt.Sprintf("The %s %s is already exist.", strcase.ToSnake(field), err.Value())
	case "digit":
		message = fmt.Sprintf("The %s must be a digit number of string.", strcase.ToSnake(field))
	case "enum":
		// first, clean/remove the comma
		cleaned := strings.Replace(param, "_", " ", -1)

		// convert 'clened' comma separated string to slice
		strSlice := strings.Fields(cleaned)
		message = fmt.Sprintf("The %s is must only value %s instead.", strcase.ToSnake(field), strings.Join(strSlice, ","))
	case "rfe":
		params := strings.Split(param, `:`)
		paramField := params[0]
		paramValue := params[1]
		message = fmt.Sprintf("The %s is required if %s = %s.", strcase.ToSnake(field), strcase.ToSnake(paramField), paramValue)
	}

	return message
}
