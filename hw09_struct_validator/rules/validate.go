package rules

import (
	"errors"
	"reflect"
)

// валидирует структутру.
// возвращает слайс ошибок валидации полей ValidationErrors или програмную ошибку.
// паникует, если v не структура.
func ValidateStruct(v reflect.Value) error {
	cnt := v.NumField()
	if cnt < 1 {
		return nil
	}

	var errStructValid = ValidationErrors{}

	// идем по полям структуры
	for i := range cnt {
		f := v.Type().Field(i)

		if !f.IsExported() { // приватное поле
			continue
		}

		// получаем набор правил для поля
		fieldRules, err := fieldRulesByTag(f.Name, f.Tag.Get("validate"))
		if err != nil {
			return err
		}
		// если нет правил, то и проверять нечего.
		if len(fieldRules.Rules) < 1 {
			return nil
		}

		errField, err := validateField(v, fieldRules)
		if err != nil {
			return err
		}
		errStructValid = append(errStructValid, errField...)
	}

	if len(errStructValid) > 0 {
		return errStructValid
	}

	return nil
}

// валидирует поле структуры.
func validateField(fieldValue reflect.Value, rules FieldRules) (ValidationErrors, error) {
	switch fieldValue.Kind() {
	case reflect.Slice, reflect.Array:
		return validateSlice(fieldValue, rules)
	default:
		return validateValue(fieldValue, rules)
	}
}

func validateSlice(fieldValue reflect.Value, rules FieldRules) (ValidationErrors, error) {
	l := fieldValue.Len()
	if l < 1 {
		return nil, nil
	}

	var vErr = make(ValidationErrors, l)
	for i := 0; i < l; i++ {
		v, err := validateValue(fieldValue.Index(i), rules)
		if err != nil {
			return nil, err
		}
		if v != nil {
			vErr = append(vErr, v...)
		}
	}

	if len(vErr) > 0 {
		return vErr, nil
	}

	return nil, nil
}

// валидируем поле со значение fieldValue, согласно правил описанных fieldRules.
func validateValue(fieldValue reflect.Value, fieldRules FieldRules) (ValidationErrors, error) {
	var errFields = ValidationErrors{}

	// перебираем все правила
	for _, rule := range fieldRules.Rules {
		// получаем функцию валидации.
		vf, err := validationFunction(fieldValue.Kind(), rule.Name)
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
		var e ValidationError
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
