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

// парсит полученную строку, возвращая массив структур с описанием правил проверки
// ожидается, что строка имеет вид 'правило:условие|правило:условие|...'
func parseRulesTag(rulesTag string) ([]RuleInfo, error) {
	rulesTag = strings.Trim(rulesTag, " ")
	if rulesTag == "" {
		return []RuleInfo{}, nil
	}

	// Разбили на отдельные описания правила: строки вида 'правило:условие'
	rs := strings.Split(rulesTag, "|")

	ri := []RuleInfo{}
	// из каждого описания правила выделяем имя правила и условие
	for _, r := range rs {
		if len(r) == 0 {
			return []RuleInfo{}, ErrEmptyRule
		}
		rule := strings.Split(r, ":")
		if len(rule) != 2 {
			return []RuleInfo{}, ErrUnknowRule
		}

		ri = append(ri, RuleInfo{Name: rule[0], Cond: rule[1]})
	}

	return ri, nil
}

func FieldRulesByTag(fieldName string, fieldTag string) (FieldRules, error) {
	frs := FieldRules{
		FieldName: fieldName,
		Rules:     []RuleInfo{},
	}
	if fieldTag == "" {
		return frs, nil
	}

	rls, err := parseRulesTag(fieldTag)
	if err != nil {
		return frs, err
	}

	frs.Rules = rls

	return frs, nil
}

// типы данных на которые распространяются правила проверки
//type TypeValue uint

var Rules = map[reflect.Kind]map[string]func(fName string, v reflect.Value, condition string) error{
	reflect.String: {
		// 'len:32' - проверка длины строки должна быть 32 символа
		"len": func(fName string, v reflect.Value, condition string) error {
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
					Field: fName,
					Err:   fmt.Errorf("length of the string not equal to %s", condition),
				}
			}
			return nil
		},
		"regexp": func(fName string, v reflect.Value, condition string) error {
			if v.Kind() != reflect.String {
				return fmt.Errorf("this rule applies only to the string")
			}

			pattern, err := regexp.Compile(condition)
			if err != nil {
				return err
			}
			if !pattern.MatchString(v.String()) {
				return ValidationError{
					Field: fName,
					Err:   fmt.Errorf("length of the string not equal to %s", condition),
				}
			}

			return ErrRuleNotImplement
		},
		"in": func(fName string, v reflect.Value, condition string) error {
			return ErrRuleNotImplement
		},
	},
}
