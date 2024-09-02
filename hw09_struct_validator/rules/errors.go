package rules

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyRule    = errors.New("the rule cannot be empty")
	ErrUnknowRule   = errors.New("unknow rule")
	ErrNotImplement = errors.New("the rule has no implementation")
)

// ошибка валидации поля структуры
type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	if len(v.Field) == 0 {
		return fmt.Sprintf("%v", v.Err)
	}
	return fmt.Sprintf("%s: %v", v.Field, v.Err)
}

// слайс ошибок валидации полей структуры
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	cnt := len(v)
	if cnt < 1 {
		return ""
	}
	return fmt.Sprintf("%d structure validation errors found", cnt)
}
