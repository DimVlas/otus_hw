package hw09structvalidator

import (
	"reflect"

	"github.com/DimVlas/otus_hw/hw09_struct_validator/rules"
)

// implemented in errors.go file
// type ValidationError struct {
// 	Field string
// 	Err   error
// }

// type ValidationErrors []ValidationError

// func (v ValidationErrors) Error() string {
// 	panic("implement me")
// }

func Validate(v interface{}) error {
	// nothing to validate
	if v == nil {
		return nil
	}

	rval := reflect.ValueOf(v)

	if rval.Kind() != reflect.Struct {
		return rules.ErrRequireStruct
	}

	return validate(rval)
}
