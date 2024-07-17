package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

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

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Student struct {
		Name   string
		Grades []int `validate:"min:2|max:5"`
	}

	UnknownValidator struct {
		field string `validate:"maxLen:5"`
	}

	InvalidLen struct {
		field string `validate:"len:a"`
	}

	InvalidRegexp struct {
		field string `validate:"regexp:\\e"`
	}

	InvalidMin struct {
		field int `validate:"min:a"`
	}

	InvalidMax struct {
		field int `validate:"max:a"`
	}

	InvalidIn struct {
		field int `validate:"in:1,a"`
	}

	DuplicateValidator struct {
		field string `validate:"len:20|len:21"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "1",
				Name:   "user",
				Age:    17,
				Email:  "invalid email",
				Role:   UserRole("user"),
				Phones: []string{"89012345678", "12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errInvalidLen},
				{Field: "Age", Err: errViolatedMin},
				{Field: "Email", Err: errNoMatchRegexp},
				{Field: "Role", Err: errStringNotIn},
				{Field: "Phones", Err: errInvalidLen},
			},
		},
		{
			in: User{
				ID:     "id1234567890123456789012345678901234",
				Name:   "user",
				Age:    51,
				Email:  "user@email.com",
				Role:   UserRole("admin"),
				Phones: []string{"89012345678"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: errViolatedMax},
			},
		},
		{
			in: User{
				ID:     "id1234567890123456789012345678901234",
				Name:   "user",
				Age:    24,
				Email:  "user@email.com",
				Role:   UserRole("admin"),
				Phones: []string{"89012345678", "89012345679"},
			},
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: errInvalidLen},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
		},
		{
			in: Response{
				Code: 403,
				Body: "Unauthorized",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errIntegerNotIn},
			},
		},
		{
			in: Response{
				Code: 200,
			},
		},
		{
			in:          "string",
			expectedErr: errNotStruct,
		},
		{
			in: Student{
				Name:   "Bob",
				Grades: []int{2, 3, 4, 5, 6},
			},
			expectedErr: ValidationErrors{
				{Field: "Grades", Err: errViolatedMax},
			},
		},
		{
			in: Student{
				Name:   "Bob",
				Grades: []int{2, 3, 4, 5},
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
			} else {
				require.Error(t, err)
				require.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestInvalidTags(t *testing.T) {
	tests := []interface{}{
		UnknownValidator{
			field: "Field",
		},
		InvalidLen{
			field: "Field",
		},
		InvalidRegexp{
			field: "Field",
		},
		InvalidMin{
			field: 1,
		},
		InvalidMax{
			field: 2,
		},
		InvalidIn{
			field: 3,
		},
		DuplicateValidator{
			field: "Field",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)
			require.Error(t, err)
		})
	}
}
