package validation

import (
	"fmt"
	"github.com/alpakih/go-api/pkg/database"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/guregu/null.v4"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func nullFloatValidator(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Float); ok {
		if valuer.Valid {
			return valuer.Float64
		}
	}
	return nil
}

func nullIntValidator(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Int); ok {
		if valuer.Valid {
			return valuer.Int64
		}
	}
	return nil
}

func nullTimeValidator(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Time); ok {
		if valuer.Valid {
			return valuer.Time
		}
	}
	return nil
}

func validateOnlyNumber(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^[0-9]*$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func validateDateOnly(fl validator.FieldLevel) bool {
	if fl.Field().String() != "" {
		regex := regexp.MustCompile(`^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$`)
		return regex.MatchString(fl.Field().String())
	}
	return true
}

func validateUnique(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramTable := param[1]

	if paramField == `` {
		return true
	}

	count := int64(0)

	if err := database.Conn().Table(paramTable).Where(paramField+"= ?", fl.Field().String()).
		Count(&count).Error; err != nil {
		log.Fatal(err)

	}
	if count > 0 {
		return false
	}
	return true
}

//ValidateCustom -- ValidateCustom
func validateEnum(field validator.FieldLevel) bool {

	if field.Param() == `` {
		return true
	}

	// first, clean/remove the comma
	cleaned := strings.Replace(field.Param(), "_", " ", -1)

	// convert 'clened' comma separated string to slice
	strSlice := strings.Fields(cleaned)

	if !itemExists(strSlice, field.Field().String()) {
		return false
	}
	return true
}

func validateUpdateUnique(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramFieldValue := param[0]
	paramTable := param[1]
	paramField := param[2]
	paramFieldCond := param[3]

	if paramFieldValue == `` {
		return true
	}

	// param field reflect.Value.
	var paramReflectValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramReflectValue = fl.Parent().Elem().FieldByName(paramFieldValue)
	} else {
		paramReflectValue = fl.Parent().FieldByName(paramFieldValue)
	}

	count := int64(0)

	if err := database.Conn().Table(paramTable).Where(paramField+" =?", fl.Field().String()).Where(paramFieldCond+" <> ?", paramReflectValue.String()).
		Count(&count).Error; err != nil {
		log.Fatal(err)

	}
	if count > 0 {
		return false
	}
	return true
}

func validateRequireIfAnotherField(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	// param field reflect.Value.
	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}

	if isEq(paramFieldValue, paramValue) == false {
		return true
	}
	return hasValue(fl)
}

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isEq(field reflect.Value, value string) bool {
	switch field.Kind() {
	case reflect.String:
		return field.String() == value
	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)
		return int64(field.Len()) == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)
		return field.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)
		return field.Uint() == p
	case reflect.Float32, reflect.Float64:
		p := asFloat(value)
		return field.Float() == p
	}
	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)
	return i
}

func asUint(param string) uint64 {
	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)
	return i
}

func asFloat(param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)
	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}