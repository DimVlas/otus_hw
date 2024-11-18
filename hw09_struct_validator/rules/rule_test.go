package rules

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func TestParseRulesTag(t *testing.T) {
	// пустой тэг
	t.Run("empty_tag", func(t *testing.T) {
		rules, err := ParseRulesTag("")

		require.Equal(t, []RuleInfo{}, rules, "was expected empty slice RuleInfo")
		require.NoError(t, err, "should no error for empty tag")
	})

	// тэг с одним правилом
	t.Run("one_rule_tag", func(t *testing.T) {
		rules, err := ParseRulesTag("rule:condition")

		require.Equal(t, []RuleInfo{{Name: "rule", Cond: "condition"}}, rules)
		require.NoError(t, err, "should no error for empty tag")
	})

	// тэг с двумя правилами
	t.Run("two_rule_tag", func(t *testing.T) {
		rules, err := ParseRulesTag("rule1:condition1|rule2:condition2")

		exp := []RuleInfo{
			{Name: "rule1", Cond: "condition1"},
			{Name: "rule2", Cond: "condition2"},
		}

		require.Equal(t, exp, rules)
		require.NoError(t, err, "should no error for empty tag")
	})

	// ошибка: пустые правила
	t.Run("error_empty_rules", func(t *testing.T) {
		rules, err := ParseRulesTag("|")

		require.Equal(t, []RuleInfo{}, rules)
		require.EqualError(t, err, ErrEmptyRule.Error())
	})

	// ошибка: некорректное правило
	t.Run("error_incorrect_rules", func(t *testing.T) {
		rules, err := ParseRulesTag("rule:cond|rule")

		require.Equal(t, []RuleInfo{}, rules)
		require.EqualError(t, err, ErrUnknowRule.Error())
	})
}

// тестировани функций-правил для значений типа "string".
func TestStringLen(t *testing.T) {
	// Неверный тип значения, передаем int вместо строки
	t.Run("len bad int value", func(t *testing.T) {
		err := Rules[reflect.String]["len"](reflect.ValueOf(123), "0")

		require.EqualError(t, err, "this rule applies only to the string")
	})
	// Неверный тип значения, передаем указатель вместо строки
	t.Run("len bad &string value", func(t *testing.T) {
		var s string = "asd"
		err := Rules[reflect.String]["len"](reflect.ValueOf(&s), "0")

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
