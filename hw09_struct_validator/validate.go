package hw09structvalidator

import (
	"reflect"

	"github.com/DimVlas/otus_hw/hw09_struct_validator/rules"
)

// валидирует структутру.
// возвращает слайс ошибок валидации полей ValidationErrors или програмную ошибку
// Паникует, если v не структура
func validate(v reflect.Value) error {
	cnt := v.NumField()
	if cnt < 1 {
		return nil
	}

	var errStructValid = rules.ValidationErrors{}

	for i := range cnt {
		f := v.Type().Field(i)

		if !f.IsExported() { // приватное поле
			continue
		}
		fieldRules, err := rules.FieldRulesByTag(f.Name, f.Tag.Get("validate"))
		if err != nil {
			return err
		}

		if len(fieldRules.Rules) < 1 {
			continue
		}

		err = validateField(v.Field(i), fieldRules)
		if err != nil {
			switch e := err.(type) {
			case rules.ValidationError:
				errStructValid = append(errStructValid, e)
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

// валидируем поле со значение fieldValue, согласно правил описанных fieldRules
func validateField(fieldValue reflect.Value, fieldRules rules.FieldRules) error {
	kind := fieldValue.Kind()

	var errFields = rules.ValidationErrors{}

	for _, rule := range fieldRules.Rules {
		f := rules.Rules[kind][rule.Name]
		err := f(fieldRules.FieldName, fieldValue, rule.Cond)
		if err != nil {
			switch e := err.(type) {
			case rules.ValidationError:
				errFields = append(errFields, e)
			default:
				return err
			}
		}
	}

	if len(errFields) > 0 {
		return errFields
	}

	return nil
}
