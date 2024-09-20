package rules

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestStringLen(t *testing.T) {
	t.Run("len bad value", func(t *testing.T) {
		f := Rules[reflect.String]["len"]

		err := f(reflect.ValueOf(123), "0")

		require.EqualError(t, err, "this rule applies only to the string")
	})

	t.Run("len bad condition", func(t *testing.T) {
		f := Rules[reflect.String]["len"]

		s := "Мой милый дом!"

		err := f(reflect.ValueOf(s), "s")

		require.EqualError(t, err, "'s' is not a valid condition for the 'len' rule")
	})

	t.Run("len not equal", func(t *testing.T) {
		f := Rules[reflect.String]["len"]

		s := "Мой милый дом!"
		l := utf8.RuneCountInString(s) - 1

		err := f(reflect.ValueOf(s), strconv.Itoa(l))

		require.EqualError(t, err, fmt.Sprintf("length of the string not equal to %d", l))
	})

	t.Run("len success", func(t *testing.T) {
		f := Rules[reflect.String]["len"]

		s := "Мой милый дом!"

		err := f(reflect.ValueOf(s), strconv.Itoa(utf8.RuneCountInString(s)))
		require.NoError(t, err)
	})

}
