package rules

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Описание правила проверки
type RuleInfo struct {
	Name string // правило проверки
	Cond string // условие правила проверки
}

// правила проверки поля
type FieldRules struct {
	FieldName string     // наименование поля
	Rules     []RuleInfo // слайс правил проверки
}

func FieldRulesByTag(fieldName string, fieldTag string) (FieldRules, error) {
	frs := FieldRules{
		FieldName: fieldName,
		Rules:     []RuleInfo{},
	}
	if fieldTag == "" {
		return frs, nil
	}

	rs := strings.Split(fieldTag, "|")

	if len(rs) < 1 {
		return frs, ErrEmptyRule
	}
	for _, r := range rs {
		rule := strings.Split(r, ":")
		if len(rule) < 2 {
			return frs, ErrUnknowRule
		}

		frs.Rules = append(frs.Rules, RuleInfo{Name: rule[0], Cond: rule[1]})
	}

	return frs, nil
}

// типы данных на которые распространяются правила проверки
//type TypeValue uint

var Rules = map[reflect.Kind]map[string]func(v reflect.Value, condition string) error{
	reflect.String: {
		// 'len:32' - проверка длины строки должна быть 32 символа
		"len": func(v reflect.Value, condition string) error {
			if v.Kind() != reflect.String {
				// это правило применимо только к строкам
				return fmt.Errorf("this rule applies only to the string")
			}
			c, err := strconv.Atoi(condition)
			if err != nil {
				// строка не является валидным условием для правила 'len'
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

			return ErrRuleNotImplement
		},
		"in": func(v reflect.Value, condition string) error {
			return ErrRuleNotImplement
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
