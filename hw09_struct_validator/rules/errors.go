package rules

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRequireStruct    = errors.New("'Validate' requires structure")
	ErrEmptyRule        = errors.New("the rule cannot be empty")
	ErrUnknowRule       = errors.New("unknow rule")
	ErrRuleNotImplement = errors.New("the rule has no implementation")
	ErrKindNoRules      = errors.New("for this field kind no validation rules")
)

var (
	// программные ошибки функций валидации.
	// // правило применимо только к строкам
	// ErrOnlyStringRule = errors.New("rule applies only to the string")
	// // правило применимо только к целым
	// ErrOnlyIntRule = errors.New("rule applies only to the int")
	// недопустимое условие для правила.
	ErrInvalidCond = errors.New("invalid condition for the rule")
	// ошибка компиляции регулярного выражения.
	ErrRegexpCompile = errors.New("regex compilation error")
)

// ошибки валидации строк.
var (
	// длина строки не равна.
	ErrStrLenNotEqual = errors.New("length of the string not equal to")
	// строка не содержит совпадений с регулярным выражением.
	ErrStrReExpNotMatch = errors.New("string does not contain any matches to the regular expression")
	// строка на входит в список.
	ErrStrNotInList = errors.New("string is not in the list")
)

// ошибки валидации целых.
var (
	// целое не может быть меньше условия.
	ErrIntCantBeLess = errors.New("cannot be less")
	// целое не содержит совпадений с регулярным выражением.
	ErrIntCantBeGreater = errors.New("cannot be greater")
	// целое на входит в список.
	ErrIntNotInList = errors.New("int is not in the list")
)

// ошибка валидации поля структуры.
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

func (v ValidationError) Unwrap() error {
	return v.Err
}

// слайс ошибок валидации полей структуры.
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	cnt := len(v)
	if cnt < 1 {
		return ""
	}

	return func() string {
		s := strings.Builder{}
		for _, e := range v {
			s.WriteString(fmt.Sprintf("field %s: %s\n", e.Field, e.Err.Error()))
		}
		return s.String()
	}()

	//fmt.Sprintf("%d structure validation errors found", cnt)
}
