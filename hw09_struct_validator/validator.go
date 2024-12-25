package hw09structvalidator

import (
	"errors"
	"reflect"

	r "github.com/DimVlas/otus_hw/hw09_struct_validator/rules"
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
		return r.ErrRequireStruct
	}

	return validateStruct(rval)
}

// валидирует структутру.
// возвращает слайс ошибок валидации полей ValidationErrors или програмную ошибку.
// паникует, если v не структура.
func validateStruct(v reflect.Value) error {
	cnt := v.NumField()
	if cnt < 1 {
		return nil
	}

	var errStructValid r.ValidationErrors
	// идем по полям структуры
	for i := range cnt {
		f := v.Type().Field(i)

		if !f.IsExported() { // приватное поле
			continue
		}

		// получаем набор правил для поля
		fieldRules, err := r.TagRules(f.Name, f.Tag.Get("validate"))
		if err != nil {
			return err
		}
		// если нет правил, то и проверять нечего.
		if len(fieldRules.Rules) < 1 {
			continue
		}

		errField, err := validateField(v.Field(i), fieldRules)
		if err != nil {
			return err
		}
		if len(errField) > 0 {
			errStructValid = append(errStructValid, errField...)
		}
	}

	if len(errStructValid) > 0 {
		return errStructValid
	}

	return nil
}

// валидирует поле структуры.
func validateField(fieldValue reflect.Value, rules r.FieldRules) (r.ValidationErrors, error) {
	switch fieldValue.Kind() {
	case reflect.Slice, reflect.Array:
		return validateSlice(fieldValue, rules)
	case reflect.Struct:
		for _, r := range rules.Rules {
			if r.Name == "nested" {
				return validateStructF(fieldValue, rules)
			}
		}
		return nil, nil
	default:
		return validateValue(fieldValue, rules)
	}
}

// валидирует вложенную структуру.
func validateStructF(fieldValue reflect.Value, _ r.FieldRules) (r.ValidationErrors, error) {
	err := validateStruct(fieldValue)
	if err != nil {
		var e r.ValidationErrors
		if errors.As(err, &e) {
			return e, nil
		}

		return nil, err
	}

	return nil, nil
}

// валидирует slice или массив.
func validateSlice(fieldValue reflect.Value, rules r.FieldRules) (r.ValidationErrors, error) {
	ln := fieldValue.Len()
	if ln < 1 {
		return nil, nil
	}

	var vErr r.ValidationErrors
	for i := 0; i < ln; i++ {
		v, err := validateValue(fieldValue.Index(i), rules)
		if err != nil {
			return nil, err
		}
		if v != nil {
			if len(vErr) > 0 {
				vErr = append(vErr, v...)
				continue
			}

			vErr = v
		}
	}

	if len(vErr) > 0 {
		return vErr, nil
	}

	return nil, nil
}

// валидируем поле со значение fieldValue, согласно правил описанных fieldRules.
func validateValue(fieldValue reflect.Value, fieldRules r.FieldRules) (r.ValidationErrors, error) {
	var errFields r.ValidationErrors
	// перебираем все правила
	for _, rule := range fieldRules.Rules {
		// получаем функцию валидации.
		vf, err := r.ValidationFunction(fieldValue.Kind(), rule.Name)
		if err != nil {
			return nil, err
		}
		// проверяем fieldValue функцией валидации.
		err = vf(fieldValue, rule.Cond)

		// если нет ошибок - переходим к следующему правилу.
		if err == nil {
			continue
		}
		// если ошибка валидации, сохраняем ее в массив ошибок валидации.
		var e r.ValidationError
		if errors.As(err, &e) {
			e.Field = fieldRules.FieldName
			errFields = append(errFields, e)
			continue
		}
		// если программная ошибка - возращаем ее, выходим.
		return nil, err
	}

	// если были ошибки валидации - возращаем их.
	if len(errFields) > 0 {
		return errFields, nil
	}

	return nil, nil
}
