package rules

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// типы данных на которые распространяются правила проверки
type TypeValue uint

var Rules = map[reflect.Kind]map[string]func(v reflect.Value, condition string) error{
	reflect.String: {
		// 'len:32' - проверка длины строки должна быть 32 символа
		"len": func(v reflect.Value, condition string) error {
			if v.Kind() != reflect.String {
				return fmt.Errorf("this rule applies only to the string")
			}
			c, err := strconv.Atoi(condition)
			if err != nil {
				return fmt.Errorf("'%s' is not a valid condition for the 'len' rule", condition)
			}

			if utf8.RuneCountInString(v.String()) != c {
				return ValidationError{
					Field: "",
					Err:   fmt.Errorf("length of the string not equal to %s", condition),
				}
			}
			return nil
		},
		"regexp": func(v reflect.Value, condition string) error {
			if v.Kind() != reflect.String {
				return fmt.Errorf("this rule applies only to the string")
			}

			pattern, err := regexp.Compile(condition)
			if err != nil {
				return err
			}
			if !pattern.MatchString(v.String()) {
				return ValidationError{
					Field: "",
					Err:   fmt.Errorf("length of the string not equal to %s", condition),
				}
			}

			return ErrNotImplement
		},
		"in": func(v reflect.Value, condition string) error {
			return ErrNotImplement
		},
	},
}

// const (
// 	Int TypeValue = iota
// 	IntSlice
// 	String
// 	StringSlice
// )

// type Rule interface {
// 	Rule() string
// 	Validate(reflect.Value) error
// }

// func New(rule string) (Rule, error) {
// 	if len(strings.Trim(rule, " ")) == 0 {
// 		return nil, ErrEmptyRule
// 	}

// 	tag := reflect.StructTag(rule)
// 	r, ok := tag.Lookup("len")
// 	if !ok {
// 		return nil, errors.New(" not len rule")
// 	}

// 	l, err := strconv.Atoi(r)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't parse the '%s' rule", rule)
// 	}
// 	lenS := LenStr{
// 		rule: rule,
// 		len:  l,
// 	}
// 	return lenS, nil
// }

// type LenStr struct {
// 	rule string
// 	len  int
// }

// func (r LenStr) Rule() string {
// 	return r.rule
// }
// func (r LenStr) Validate(val reflect.Value) error {
// 	if val.Kind() == reflect.String {
// 		s := val.String()

// 		if len(s) != r.len {
// 			return fmt.Errorf("the length should be equal to %d", r.len)
// 		}
// 		return nil
// 	}

// 	return fmt.Errorf("the '%s' rule only applies to values of type'%s'", r.rule, "string")
// }
