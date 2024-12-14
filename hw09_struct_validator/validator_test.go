package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DimVlas/otus_hw/hw09_struct_validator/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Response{
				Code: 200,
				Body: "{}",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 100,
				Body: "{}",
			},
			expectedErr: rules.ValidationErrors{
				rules.ValidationError{
					Field: "Code",
					Err:   fmt.Errorf("%w 200,404,500", rules.ErrIntNotInList),
				},
			},
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "qwert",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "qwerty",
			},
			expectedErr: rules.ValidationErrors{
				rules.ValidationError{
					Field: "Version",
					Err:   fmt.Errorf("%w 5", rules.ErrStrLenNotEqual),
				},
			},
		},
		{
			in: User{
				ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
				Name:   "User1",
				Age:    18,
				Email:  "User1@mail.com",
				Role:   "admin",
				Phones: []string{"12345678901", "98765432101"},
				meta:   json.RawMessage(``),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "pD4tNeo-t0OGE_ooz3WqxAcyFeuF6AUk6mQf",
				Name:   "User1",
				Age:    16,
				Email:  "User1@mail.com.dot",
				Role:   "employee",
				Phones: []string{"12345678901", "9876543210"},
				meta:   json.RawMessage(``),
			},
			expectedErr: rules.ValidationErrors{
				rules.ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%w 18", rules.ErrIntCantBeLess),
				},
				rules.ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("%w %s", rules.ErrStrReExpNotMatch, "'^\\w+@\\w+\\.\\w+$'"),
				},
				rules.ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%w %s", rules.ErrStrNotInList, "'admin,stuff'"),
				},
				rules.ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w %v", rules.ErrStrLenNotEqual, 11),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
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
