package rules

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Описание правила проверки.
type RuleInfo struct {
	Name string // правило проверки
	Cond string // условие правила проверки
}

// правила проверки поля.
type FieldRules struct {
	FieldName string     // наименование поля
	Rules     []RuleInfo // слайс правил проверки
}

// описание функции валидации.
// функция валидации задается для каждого правила.
// проверяет соответствие значения v условию condition.
// если kind значения v не соответсвует ожидаемому, то скорее всего будет panic;
// если условие condition будет пустым, то функция вернет nil;
// если условие condition будет некорректным, вернет error;
// если занчение v не удовлетворяет условию, вернется ошибка типа ValidationError с пустым полем ValidationError.Field.
type Validator func(v reflect.Value, condition string) error

// маппа в которой по типам полей содержится маппа с типами правил и функциями валидации для каждого типа правила.
var validators = map[reflect.Kind]map[string]Validator{
	reflect.String: {
		// 'len:32' - проверка длины строки должна быть 32 символа.
		"len": func(v reflect.Value, condition string) error {
			c, err := strconv.Atoi(condition)
			if err != nil {
				// 'condition' недопустимое условие для правила 'len'.
				return fmt.Errorf("'%s' %w '%s'", condition, ErrInvalidCond, "len")
			}

			if utf8.RuneCountInString(v.String()) != c {
				return ValidationError{
					Err: fmt.Errorf("%w %s", ErrStrLenNotEqual, condition),
				}
			}
			return nil
		},
		"regexp": func(v reflect.Value, condition string) error {
			if condition == "" {
				// 'condition' недопустимое условие для правила 'regexp'.
				return fmt.Errorf("'%s' %w '%s'", condition, ErrInvalidCond, "regexp")
			}

			pattern, err := regexp.Compile(condition)
			if err != nil {
				return fmt.Errorf("'%s' %w: %w", condition, ErrRegexpCompile, err)
			}

			if !pattern.MatchString(v.String()) {
				return ValidationError{
					Err: fmt.Errorf("%w '%s'", ErrStrReExpNotMatch, condition),
				}
			}
			return nil
		},
		"in": func(v reflect.Value, condition string) error {
			if condition == "" {
				// 'condition' недопустимое условие для правила 'regexp'.
				return fmt.Errorf("'%s' %w '%s'", condition, ErrInvalidCond, "in")
			}

			if !strings.Contains(condition, v.String()) {
				return ValidationError{
					Err: fmt.Errorf("%w '%s'", ErrStrNotInList, condition),
				}
			}

			return nil
		},
	},
	reflect.Int: {
		// 'min:32' - число не может быть меньше 10.
		"min": func(v reflect.Value, condition string) error {
			c, err := strconv.ParseInt(condition, 0, 0)
			if err != nil {
				// 'condition' недопустимое условие для правила 'min'.
				return fmt.Errorf("'%s' %w '%s': %w", condition, ErrInvalidCond, "min", err)
			}

			if v.Int() < c {
				return ValidationError{
					Err: fmt.Errorf("%w %s", ErrIntCantBeLess, condition),
				}
			}
			return nil
		},
		// 'max:32' - число не может быть больше 10.
		"max": func(v reflect.Value, condition string) error {
			c, err := strconv.ParseInt(condition, 0, 0)
			if err != nil {
				// 'condition' недопустимое условие для правила 'min'.
				return fmt.Errorf("'%s' %w '%s': %w", condition, ErrInvalidCond, "max", err)
			}

			if v.Int() > c {
				return ValidationError{
					Err: fmt.Errorf("%w %s", ErrIntCantBeGreater, condition),
				}
			}
			return nil
		},
		// 'in:32,33' - число не входит в список 32,33.
		"in": func(v reflect.Value, condition string) error {
			cl := strings.Split(condition, ",")
			var isValid bool
			for _, c := range cl {
				i, err := strconv.ParseInt(c, 0, 0)
				if err != nil {
					// 'condition' недопустимое условие для правила 'in'.
					return fmt.Errorf("'%s' %w '%s': %w", condition, ErrInvalidCond, "in", err)
				}

				if v.Int() == i {
					isValid = true
					break
				}
			}

			if !isValid {
				return ValidationError{
					Err: fmt.Errorf("%w %s", ErrIntNotInList, condition),
				}
			}

			return nil
		},
	},
}

// возвращает функцию валидации для типа kind и правила rule.
func ValidationFunction(kind reflect.Kind, rule string) (Validator, error) {
	r, ok := validators[kind]
	if !ok {
		return nil, fmt.Errorf("'%s' %w", kind, ErrKindNoRules)
	}

	fv, ok := r[rule]
	if !ok {
		return nil, fmt.Errorf("'%s' %w", rule, ErrUnknowRule)
	}

	return fv, nil
}

// получает из тэга fieldTag струтуру FieldRules с правилами валидации для поля с именем fieldName.
func TagRules(fieldName string, fieldTag string) (FieldRules, error) {
	rls, err := parseRulesTag(fieldTag)
	if err != nil {
		return FieldRules{
			FieldName: fieldName,
			Rules:     []RuleInfo{},
		}, err
	}

	return FieldRules{
		FieldName: fieldName,
		Rules:     rls,
	}, nil
}

// парсит полученную строку, возвращая массив структур с описанием правил проверки.
// ожидается, что строка имеет вид 'правило:условие|правило:условие|...'.
func parseRulesTag(rulesTag string) ([]RuleInfo, error) {
	rulesTag = strings.Trim(rulesTag, " ")
	if rulesTag == "" {
		return []RuleInfo{}, nil
	}

	// Разбили на отдельные описания правила: строки вида 'правило:условие'.
	rs := strings.Split(rulesTag, "|")

	ri := []RuleInfo{}
	// из каждого описания правила выделяем имя правила и условие.
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
