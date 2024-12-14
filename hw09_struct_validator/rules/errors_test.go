package rules

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name string
		data ValidationError
		exp  string
	}{
		{
			name: "ValidationError_full",
			data: ValidationError{
				Field: "field",
				Err:   errors.New("test error"),
			},
			exp: "field: test error",
		},
		{
			name: "ValidationError_without_field",
			data: ValidationError{
				Field: "",
				Err:   errors.New("test error"),
			},
			exp: "test error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.data.Error()

			require.Equal(t, test.exp, s)
		})
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		data ValidationErrors
		exp  string
	}{
		{
			name: "ValidationErrors_full",
			data: ValidationErrors{
				{
					Field: "f1",
					Err:   errors.New("error1"),
				},
				{
					Field: "f2",
					Err:   errors.New("error2"),
				},
			},
			exp: "field f1: error1\nfield f2: error2\n",
		},
		{
			name: "ValidationErrors_empty",
			data: ValidationErrors{},
			exp:  "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.data.Error()

			require.Equal(t, test.exp, s)
		})
	}
}
