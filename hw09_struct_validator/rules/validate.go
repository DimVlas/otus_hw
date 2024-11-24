package rules

import (
	"reflect"
)

// валидирует структутру.
// возвращает слайс ошибок валидации полей ValidationErrors или програмную ошибку
// Паникует, если v не структура
func ValidateStruct(v reflect.Value) error {
	cnt := v.NumField()
	if cnt < 1 {
		return nil
	}

	var errStructValid = ValidationErrors{}

	for i := range cnt {

		f := v.Type().Field(i)

		if !f.IsExported() { // приватное поле
			continue
		}

		if err := validateField(f, v.Field(i)); err != nil {
			switch e := err.(type) {
			case ValidationErrors:
				errStructValid = append(errStructValid, e...)
			default:
				return err
			}
		}
	}

	if len(errStructValid) > 0 {
		return errStructValid
	}

	return nil
}

// валидирует поле структуры
func validateField(fieldInfo reflect.StructField, fieldValue reflect.Value) error {
	fieldRules, err := FieldRulesByTag(fieldInfo.Name, fieldInfo.Tag.Get("validate"))
	if err != nil {
		return err
	}
	// если не правил, то и проверять нечего
	if len(fieldRules.Rules) < 1 {
		return nil
	}

	return validateFieldValue(fieldValue, fieldRules)
}

// валидируем поле со значение fieldValue, согласно правил описанных fieldRules
func validateFieldValue(fieldValue reflect.Value, fieldRules FieldRules) error {
	var errFields = ValidationErrors{}

	// перебираем все правила
	for _, rule := range fieldRules.Rules {
		// проверяем fieldValue очередным правилом
		err := validateFieldRules(fieldValue, fieldValue.Kind(), rule.Name, rule.Cond)
		// если нет ошибок - переходим к следующему правилу
		if err == nil {
			continue
		}
		// если ошибка валидации, сохраняем ее в массив ошибок валидации
		if vErr, ok := err.(ValidationError); ok {
			vErr.Field = fieldRules.FieldName
			errFields = append(errFields, vErr)
			continue
		}
		// если программная ошибка - возращаем ее, выходим
		return err
	}

	// если были ошибки валидации - возращаем их
	if len(errFields) > 0 {
		return errFields
	}

	return nil
}
