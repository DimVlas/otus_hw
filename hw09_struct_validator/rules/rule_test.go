package rules

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

// тестировани функций-правил для значений типа "string".
func TestStringLen(t *testing.T) {
	// Неверный тип значения, передаем int вместо строки
	t.Run("len bad value", func(t *testing.T) {
		err := Rules[reflect.String]["len"](reflect.ValueOf(123), "0")

		require.EqualError(t, err, "this rule applies only to the string")
	})

	// неверное значение условия для правила
	t.Run("len bad condition", func(t *testing.T) {
		s := "Мой милый дом!"

		err := Rules[reflect.String]["len"](reflect.ValueOf(s), "s")

		require.EqualError(t, err, "'s' is not a valid condition for the 'len' rule")
	})

	// проверка провалена
	t.Run("len not equal", func(t *testing.T) {
		f := Rules[reflect.String]["len"]

		s := "Мой милый дом!"
		l := utf8.RuneCountInString(s) - 1

		err := f(reflect.ValueOf(s), strconv.Itoa(l))

		// TODO: здесь должна быть не просто ошибка, а ошибка ValidationError, это тож надо проверить

		require.EqualError(t, err, fmt.Sprintf("length of the string not equal to %d", l))
	})

	// проверка успешна
	t.Run("len success", func(t *testing.T) {
		s := "Мой милый дом!"

		err := Rules[reflect.String]["len"](reflect.ValueOf(s), strconv.Itoa(utf8.RuneCountInString(s)))
		require.NoError(t, err)
	})

}
