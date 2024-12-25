package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DimVlas/otus_hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the function on different structures and other types.
type (
	UserRole string

	EmptyStruct struct{}

	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Name    string `validate:"min:5"`
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Product struct {
		Price float32 `validate:"min:1"`
		Name  string  `validate:"len:10"`
	}

	UserResponses struct {
		User      User `validate:"min:1|nested"`
		Responses []Response
	}

	UserTags struct {
		User User `validate:"min:1"`
		Tags string
	}
	UserApp struct {
		User User
		App  App `validate:"nested"`
	}
)

var tests = []struct {
	name        string
	in          interface{}
	expectedErr error
}{
	{
		name: "validate_slice_err",
		in: struct {
			Codes []int `validate:"len:5"`
		}{
			Codes: []int{1, 2, 3},
		},
		expectedErr: fmt.Errorf("'len' %w", rules.ErrUnknowRule),
	},
	{
		name: "validate_nested_struct_slice_err_validation",
		in: UserResponses{
			User: User{
				ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
				Name:   "User1",
				Age:    33,
				Email:  "User1@mail.com",
				Role:   "admin",
				Phones: []string{"1234567890", "9876543210"},
				meta:   json.RawMessage(``),
			},
			Responses: []Response{},
		},
		expectedErr: rules.ValidationErrors{
			rules.ValidationError{
				Field: "Phones",
				Err:   fmt.Errorf("%w 11", rules.ErrStrLenNotEqual),
			},
			rules.ValidationError{
				Field: "Phones",
				Err:   fmt.Errorf("%w 11", rules.ErrStrLenNotEqual),
			},
		},
	},
	{
		name: "validate_nested_struct",
		in: UserResponses{
			User: User{
				ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
				Name:   "User1",
				Age:    33,
				Email:  "User1@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "98765432101"},
				meta:   json.RawMessage(``),
			},
			Responses: []Response{},
		},
		expectedErr: nil,
	},
	{
		name: "validate_no_validate_nested_struct",
		in: UserTags{
			User: User{
				ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
				Name:   "User1",
				Age:    33,
				Email:  "User1@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "98765432101"},
				meta:   json.RawMessage(``),
			},
			Tags: "",
		},
		expectedErr: nil,
	},
	{
		name:        "validate_struct_nil",
		in:          nil,
		expectedErr: nil,
	},
	{
		name:        "validate_not_struct",
		in:          "test",
		expectedErr: rules.ErrRequireStruct,
	},
	{
		name:        "validate_struct_empty",
		in:          EmptyStruct{},
		expectedErr: nil,
	},
	{
		name: "validate_rule_empty",
		in: struct {
			Field string `validate:"rule:cond|rule"`
		}{
			Field: "qwert",
		},
		expectedErr: rules.ErrUnknowRule,
	},
	{
		name: "validate_field_private",
		in: struct {
			field string `validate:"rule:cond"`
		}{
			field: "qwert",
		},
		expectedErr: nil,
	},
	{
		name: "validate_field_no_rules",
		in: struct {
			Field string
		}{
			Field: "qwert",
		},
		expectedErr: nil,
	},
	{
		name: "validate_err_unknow_rule_func",
		in: struct {
			Field string `validate:"rule:cond"`
		}{
			Field: "qwert",
		},
		expectedErr: fmt.Errorf("'%s' %w", "rule", rules.ErrUnknowRule),
	},
	{
		name: "validate_err_bad_condition",
		in: struct {
			Field string `validate:"len:cond"`
		}{
			Field: "qwert",
		},
		expectedErr: fmt.Errorf("'%s' %w '%s'", "cond", rules.ErrInvalidCond, "len"),
	},
	{
		name: "validate_valid_err_cant_be_greate",
		in: User{
			ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
			Name:   "User1",
			Age:    51,
			Email:  "User1@mail.com",
			Role:   "admin",
			Phones: []string{},
			meta:   json.RawMessage(``),
		},
		expectedErr: rules.ValidationErrors{
			rules.ValidationError{
				Field: "Age",
				Err:   fmt.Errorf("%w 50", rules.ErrIntCantBeGreater),
			},
		},
	},
	{
		name: "validate_valid_check_arr",
		in: User{
			ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
			Name:   "User1",
			Age:    25,
			Email:  "User1@mail.dot",
			Role:   "admin",
			Phones: []string{"12345678901", "98765432101"},
			meta:   json.RawMessage(``),
		},
		expectedErr: nil,
	},
	{
		name: "validate_nested_structerr_err",
		in: UserApp{
			User: User{},
			App: App{
				Name:    "App1",
				Version: "qwert",
			},
		},
		expectedErr: fmt.Errorf("'min' %w", rules.ErrUnknowRule),
	},
}

func TestValidate(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d: %s", i, tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			if assert.Error(t, err) {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
