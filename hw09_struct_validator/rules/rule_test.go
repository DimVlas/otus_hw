package rules

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// тестируем только 2 функции пакета rules: fieldRulesByTag и validationFunction,
// и функций валидации из маппы rules
// отстальные функции введены для структуризации кода и не выполняют свои проверки,
// а траслируют результат этих функций

func TestValidationFunction(t *testing.T) {
	type testData struct {
		name     string
		kind     reflect.Kind
		cond     string
		expIsNil bool
		err      error
		mess     string
	}

	tests := []testData{
		{
			name:     "err_kind_no_rules",
			kind:     reflect.Invalid,
			cond:     "rule",
			expIsNil: true,
			err:      fmt.Errorf("'%s' %w", reflect.Invalid, ErrKindNoRules),
			mess:     "expected error " + fmt.Errorf("'%s' %w", reflect.Invalid, ErrKindNoRules).Error(),
		},
		{
			name:     "unknow_rule",
			kind:     reflect.String,
			cond:     "rule",
			expIsNil: true,
			err:      fmt.Errorf("'%s' %w", "rule", ErrUnknowRule),
			mess:     "expected error " + fmt.Errorf("'%s' %w", "rule", ErrUnknowRule).Error(),
		},
		{
			name:     "success",
			kind:     reflect.String,
			cond:     "len",
			expIsNil: false,
			err:      nil,
			mess:     "expected not nil validation function",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, err := validationFunction(test.kind, test.cond)

			if test.expIsNil {
				require.Nil(t, fn, test.mess)
			} else {
				require.NotNil(t, fn, test.mess)
			}

			if test.err != nil {
				require.EqualError(t, err, test.err.Error(), test.mess)
			} else {
				require.NoError(t, err, test.mess)
			}
		})
	}
}

func TestFieldRulesByTag(t *testing.T) {
	type testData struct {
		name  string
		field string
		tag   string
		exp   FieldRules
		err   error
		mess  string
	}

	tests := []testData{
		{
			name:  "empty_tag",
			field: "field",
			tag:   "",
			exp: FieldRules{
				FieldName: "field",
				Rules:     []RuleInfo{},
			},
			err:  nil,
			mess: "should no error for empty tag",
		},
		{
			name:  "one_rule_tag",
			field: "field",
			tag:   "rule:condition",
			exp: FieldRules{
				FieldName: "field",
				Rules:     []RuleInfo{{Name: "rule", Cond: "condition"}},
			},
			err:  nil,
			mess: "should no error for tag with one rule",
		},
		{
			name:  "two_rule_tag",
			field: "field",
			tag:   "rule1:condition1|rule2:condition2",
			exp: FieldRules{
				FieldName: "field",
				Rules: []RuleInfo{
					{Name: "rule1", Cond: "condition1"},
					{Name: "rule2", Cond: "condition2"},
				}},
			err:  nil,
			mess: "should no error for tag with two rule",
		},
		{
			name:  "error_empty_rules",
			field: "field",
			tag:   "|",
			exp: FieldRules{
				FieldName: "field",
				Rules:     []RuleInfo{}},
			err:  ErrEmptyRule,
			mess: "",
		},
		{
			name:  "error_incorrect_rules",
			field: "field",
			tag:   "rule:cond|rule",
			exp: FieldRules{
				FieldName: "field",
				Rules:     []RuleInfo{}},
			err:  ErrUnknowRule,
			mess: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := fieldRulesByTag(test.field, test.tag)

			require.Equal(t, test.exp, r, test.mess)

			if err != nil {
				require.EqualError(t, err, test.err.Error(), test.mess)
			} else {
				require.NoError(t, err, test.mess)
			}
		})
	}
}

type validatorTestData struct {
	name   string
	kind   reflect.Kind
	rule   string
	cond   string
	val    reflect.Value
	expErr error
}

func (v *validatorTestData) validatorFunc() Validator {
	return validators[v.kind][v.rule]
}

var validatorTests = []validatorTestData{
	// Значение типа reflect.String, правило len
	{
		// неверный тип значения
		name:   "string_len__err_bad_type_value_int",
		kind:   reflect.String,
		rule:   "len",
		cond:   "0",
		val:    reflect.ValueOf(123),
		expErr: ErrOnlyStringRule,
	},
	{
		// неверный тип значения
		name: "string_len__err_bad_type_value_&string",
		kind: reflect.String,
		rule: "len",
		cond: "0",
		val: func() reflect.Value {
			s := "abc"
			return reflect.ValueOf(&s)
		}(),
		expErr: ErrOnlyStringRule,
	},
	{
		// неверное условия для правила
		name: "string_len__err_bad_condition",
		kind: reflect.String,
		rule: "len",
		cond: "s",
		val: func() reflect.Value {
			s := "Мой милый дом!"
			return reflect.ValueOf(s)
		}(),
		expErr: ErrInvalidCond,
	},
	{
		// ошибка валидации - длина не соответствует
		name: "string_len__err_validation_len_not_equal",
		kind: reflect.String,
		rule: "len",
		cond: "5",
		val: func() reflect.Value {
			s := "Мой милый дом!"
			return reflect.ValueOf(s)
		}(),
		expErr: ValidationError{
			Err: ErrLenNotEqual,
		},
	},
	{
		// успешная валидация
		name: "string_len__success",
		kind: reflect.String,
		rule: "len",
		cond: "5",
		val: func() reflect.Value {
			s := "милый"
			return reflect.ValueOf(s)
		}(),
		expErr: nil,
	},
	// Значение типа reflect.String, правило regexp
	{
		// неверный тип значения
		name:   "string_regexp__err_bad_type_value_int",
		kind:   reflect.String,
		rule:   "regexp",
		cond:   "",
		val:    reflect.ValueOf(123),
		expErr: ErrOnlyStringRule,
	},
	{
		// неверное условия для правила
		name:   "string_regexp__err_bad_condition",
		kind:   reflect.String,
		rule:   "regexp",
		cond:   "",
		val:    reflect.ValueOf("Дом, милый дом!"),
		expErr: ErrInvalidCond,
	},
	{
		// неверное регулярное выражение
		name:   "string_regexp__err_bad_regexp",
		kind:   reflect.String,
		rule:   "regexp",
		cond:   `/[`,
		val:    reflect.ValueOf("Дом, милый дом!"),
		expErr: ErrRegexpCompile,
	},
	{
		// ошибка валидации - нет совпадения с регулярным выражением
		name: "string_regexp__err_validation_regexp_not_match",
		kind: reflect.String,
		rule: "regexp",
		cond: `dam`,
		val:  reflect.ValueOf("Дом, милый дом!"),
		expErr: ValidationError{
			Err: ErrReExpNotMatch,
		},
	},
	{
		// успешная валидация
		name:   "string_regexp__success",
		kind:   reflect.String,
		rule:   "regexp",
		cond:   `дом`,
		val:    reflect.ValueOf("Дом, милый дом!"),
		expErr: nil,
	},
}

func TestValidator(t *testing.T) {
	for _, test := range validatorTests {
		t.Run(test.name, func(t *testing.T) {
			err := test.validatorFunc()(test.val, test.cond)

			if test.expErr == nil {
				require.NoError(t, err)
				return
			}

			if e, ok := test.expErr.(ValidationError); ok {
				require.ErrorIs(t, err, e.Err)
				return
			}

			require.ErrorIs(t, err, test.expErr)
		})
	}
}
