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
					Err:   fmt.Errorf("%w 200,404,500", rules.ErrIntNotIntList),
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
