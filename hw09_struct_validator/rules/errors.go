package rules

import (
	"errors"
	"fmt"
)

var (
	ErrRequireStruct    = errors.New("'Validate' requires structure")
	ErrEmptyRule        = errors.New("the rule cannot be empty")
	ErrUnknowRule       = errors.New("unknow rule")
	ErrRuleNotImplement = errors.New("the rule has no implementation")
	ErrKindNoRules      = errors.New("for this field kind no validation rules")
)

// программные ошибки функций валидации
var (
	// правило применимо только к строкам
	ErrOnlyStringRule = errors.New("rule applies only to the string")
	// недопустимое условие для правила
	ErrInvalidCond = errors.New("invalid condition for the rule")
)

// ошибки валидации строк
var (
	// длина строки не равна
	ErrNotEqualLen = errors.New("length of the string not equal to")
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
