package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
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
		Phones []string        `validate:"len:11|regexp:^\\d+$"`
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
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"1234567890I"},
			},
			expectedErr: fmt.Errorf("Phones: validate: validate: must match regexp ^\\d+$"),
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Alice",
				Age:    17,
				Email:  "alice@example.com",
				Role:   "stuff",
				Phones: []string{"12345678901"},
			},

			expectedErr: fmt.Errorf("Age: validate: must be at least 18"),
		},
		{
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "1.0"},
			expectedErr: fmt.Errorf("Version: validate: must be 5 characters long"),
		},
		{
			in:          Response{Code: 200, Body: "OK"},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 666, Body: "Burn in hell"},
			expectedErr: fmt.Errorf("Code: validate: must be one of [200 404 500]"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt, i := tt, i
			t.Parallel()

			err := Validate(tt.in)

			if (err != nil && tt.expectedErr == nil) ||
				(err == nil && tt.expectedErr != nil) ||
				(err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Test case %d failed: expected error '%v', got '%v'", i, tt.expectedErr, err)
			}

			// Не научился, не получилось...
			// switch {
			// case tt.expectedErr == nil && err != nil:
			// 	t.Errorf("Test case %d: unexpected error: %v", i, err)
			// case tt.expectedErr != nil && err == nil:
			// 	t.Errorf("Test case %d: expected error: %v, got nil", i, tt.expectedErr)
			// case tt.expectedErr != nil && err != nil:
			// 	if !errors.Is(err, tt.expectedErr) {
			// 		t.Errorf("Test case %d: expected error: %v, got: %v", i, tt.expectedErr, err)
			// 	}

			// 	var validationErrors ValidationErrors
			// 	if errors.As(err, &validationErrors) {
			// 		fmt.Println("DEBUG:", validationErrors)
			// 	}
			// }
		})
	}
}
