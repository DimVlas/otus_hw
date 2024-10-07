package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/DimVlas/otus_hw/hw09_struct_validator/rules"
)

func validate(v reflect.Value, vTag string) error {
	switch v.Kind() {
	case reflect.Struct:
		return validateStruct(v, vTag)
	case reflect.String:
		return validateString(v, vTag)

	default:
		return fmt.Errorf("'%s' unknown data type for validation", v.Type())
	}
}

// Валидирует структутру.
// Принимает
// v типа reflect.Value струтуры, если v.Kind() != reflect.Struct - panic;
// vTag - тэг с правилами проверки
// Возвращает error, если произошла программная ошибка
// или ValidationErrors, если были ошибки валидации.
// Валидация проходит в 2 этапа
// сначала валидируется целиком структура, согласно правилам описанным в vTag,
// далее запускается валидация каждого публичного поля структуры
func validateStruct(v reflect.Value, vTag string) error {
	var errStruct rules.ValidationErrors

	err := validStructR(v, vTag)
	if err != nil {
		var ok bool
		if errStruct, ok = err.(rules.ValidationErrors); !ok {
			return err
		}
	}

	err = validStructF(v)
	if err != nil {
		errVaild, ok := err.(rules.ValidationErrors)
		if !ok {
			return err
		}
		if errStruct == nil {
			return errVaild
		}

		errStruct = append(errStruct, errVaild...)
		return errStruct
	}

	return err
}

// валидирует структуру "целиком" соглано правилам валидации структуры
// Validation Structure Rules
func validStructR(v reflect.Value, vTag string) error {
	if len(vTag) == 0 {
		return nil
	}

	return nil
}

// валидирует поля структуры
// Validation Structure Fields
func validStructF(v reflect.Value) error {
	cnt := v.NumField()
	if cnt < 1 {
		return nil
	}

	var errStruct = rules.ValidationErrors{}
	for i := 0; i < cnt; i++ {
		f := v.Type().Field(i)

		if f.PkgPath != "" { // приватное поле
			continue
		}

		tag := f.Tag.Get("validate")
		if len(tag) == 0 {
			continue
		}
		// рекурсивно выполняем валидацию значения поля
		errField := validate(v.Field(i), tag)
		if errField != nil {
			errValid, ok := errField.(rules.ValidationErrors)
			if !ok {
				return errField
			}
			for i := range errValid {
				errValid[i].Field = f.Name
			}
			errStruct = append(errStruct, errValid...)
		}
	}
	if len(errStruct) > 0 {
		return errStruct
	}
	return nil
}

func validateString(v reflect.Value, vTag string) error {
	if len(vTag) == 0 {
		return rules.ErrEmptyRule
	}
	rule := strings.Split(vTag, ":")
	if len(rule) != 2 {
		return rules.ErrUnknowRule
	}

	f, ok := rules.Rules[v.Kind()][rule[0]]
	if !ok {
		return rules.ErrUnknowRule
	}

	return f(v, rule[1])
}

// func checkStruct(v interface{}) error {
// 	tp := reflect.TypeOf(v)
// 	vl := reflect.ValueOf(v)

// 	if tp.Kind() != reflect.Struct {
// 		return rules.ErrRequireStruct
// 	}

// 	cnt := tp.NumField()
// 	if cnt < 1 {
// 		return nil
// 	}

// 	for i := 0; i < cnt; i++ {
// 		field := tp.Field(i)
// 		if field.PkgPath != "" { // приватное поле
// 			continue
// 		}
// 		tag := field.Tag
// 		tagValidate := tag.Get("validate")
// 		if tagValidate == "" {
// 			continue
// 		}

// 		fv := vl.FieldByName(field.Name)

// 		fmt.Println(fv)
// 		fmt.Printf("field: %s\n", field.Name)
// 		fmt.Print("rules:\n")
// 		rs := strings.Split(tagValidate, "|")
// 		for _, r := range rs {
// 			fmt.Printf("       %s\n", r)
// 		}
// 		fmt.Println()
// 	}

// 	return nil
// }
