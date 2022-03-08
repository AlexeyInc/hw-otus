package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
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
		Version string `validate:"len:4|in:beta,test,prod"`
	}

	Token struct {
		Header    []byte
		Payload   []byte `validate:"min:18|max:50"`
		Signature []byte `validate:"in:20,44,100"`
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
			// default valid case for User
			in: User{
				"Some___Random___Id___With___Length36",
				"Alex", 24, "test@mail.com",
				"admin",
				[]string{"12345678901"},
				json.RawMessage(`{"precomputed": true}`),
			},
			expectedErr: errors.New(""),
		},
		{
			// default valid case for App
			in:          App{"test"},
			expectedErr: errors.New(""),
		},
		{
			// default valid case for Token
			in:          Token{[]byte{1, 2, 3}, []byte{18, 23, 45}, []byte{20, 44}},
			expectedErr: errors.New(""),
		},
		{
			// default valid case for Response
			in:          Response{200, "someText"},
			expectedErr: errors.New(""),
		},
		{
			// non valid Id, Age, Email, Role, Phones
			in: User{
				"someRandomIdWithLength24",
				"Alex", 51, "test@mail",
				"testRole",
				[]string{"123"},
				json.RawMessage(`{"precomputed": true}`),
			},
			expectedErr: errors.New(
				`Field Name: ID\n
				Errors: string len should be: 36\n
				Field Name: Age\n
				Errors: max should be: 50\n
				Field Name: Email\n
				Errors: string should match pattern: ^\w+@\w+\.\w+$\n
				Field Name: Role\n
				Errors: value should be in: admin,stuff\n
				Field Name: Phones\n
				Errors: string len should be: 11`),
		},
		{
			// non valid Id and Role
			in: User{
				"someRandomIdWithLength24",
				"Alex", 24, "test@mail.com",
				"otherRole",
				[]string{"12345678901"},
				json.RawMessage(`{"precomputed": true}`),
			},
			expectedErr: errors.New(
				`Field Name: ID\n
				Errors: string len should be: 36\n
				Field Name: Role\n
				Errors: value should be in: admin,stuff`),
		},
		{
			// check multiple errors
			in: App{"dev"},
			expectedErr: errors.New(
				`Field Name: Version\n
				Errors: string len should be: 4; value should be in: beta,test,prod`),
		},
		{
			// not valid byte slices (min/max err)
			in: Token{[]byte{1, 2, 3}, []byte{12, 23, 55}, []byte{20, 44}},
			expectedErr: errors.New(`Field Name: Payload\n
			Errors: min should be: 18; max should be: 50`),
		},
		{
			// not valid byte slices ('in' err)
			in: Token{[]byte{1, 2, 3}, []byte{18, 23, 45}, []byte{20, 21}},
			expectedErr: errors.New(`Field Name: Signature\n
			Errors: value should be in: 20,44,100`),
		},
		{
			// not valid int ('in' err)
			in: Response{201, "someText"},
			expectedErr: errors.New(`Field Name: Code\n
			Errors: value should be in: 200,404,500`),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			resultErr := Validate(tt.in)

			expectedErr := formatErrMsg(tt.expectedErr.Error())
			actualErr := resultErr.Error()

			require.Equal(t, expectedErr, actualErr,
				"Error should be: %v, got: %v", tt.expectedErr, resultErr)

			_ = tt
		})
	}
}

func formatErrMsg(err string) string {
	err = strings.ReplaceAll(err, "\t", "")
	formattedExpectedErr := strings.ReplaceAll(err, "\\n", "")
	return formattedExpectedErr
}
