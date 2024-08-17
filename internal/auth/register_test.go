package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegisterInput(t *testing.T) {

	testCases := []struct {
		name          string
		input         InputUser
		expectedError bool
		expectedUser  InputUser
	}{
		{
			name: "Valid Input",
			input: InputUser{
				Username:        "validuser123",
				Password:        "ValidPassword123!",
				PasswordConfirm: "ValidPassword123!",
				Email:           "valid@email.com",
			},
			expectedError: false,
			expectedUser: InputUser{
				Username:        "validuser123",
				Password:        "ValidPassword123!",
				PasswordConfirm: "ValidPassword123!",
				Email:           "valid@email.com",
			},
		},
		{
			name: "Invalid Username (too short)",
			input: InputUser{
				Username:        "short",
				Password:        "ValidPassword123!",
				PasswordConfirm: "ValidPassword123!",
				Email:           "valid@email.com",
			},
			expectedError: true,
		},
		{
			name: "Invalid Password (too short)",
			input: InputUser{
				Username:        "validuser123",
				Password:        "short",
				PasswordConfirm: "short",
				Email:           "valid@email.com",
			},
			expectedError: true,
		},
		{
			name: "Invalid Password (no symbol)",
			input: InputUser{
				Username:        "validuser123",
				Password:        "NoSymbol123",
				PasswordConfirm: "NoSymbol123",
				Email:           "valid@email.com",
			},
			expectedError: true,
		},
		{
			name: "Passwords Don't Match",
			input: InputUser{
				Username:        "validuser123",
				Password:        "ValidPassword123!",
				PasswordConfirm: "DifferentPassword",
				Email:           "valid@email.com",
			},
			expectedError: true,
		},
		{
			name: "Invalid Password (no number)",
			input: InputUser{
				Username:        "validuser123",
				Password:        "nonValidPassword!",
				PasswordConfirm: "nonValidPassword!",
				Email:           "valid@email.com",
			},
			expectedError: true,
		},
		{
			name: "2 invalid inputs",
			input: InputUser{
				Username:        "short",
				Password:        "nonValidPassword!",
				PasswordConfirm: "nonValidPassword!",
				Email:           "valid@email.com",
			},
			expectedError: true,
		},
		{
			name: "valid Input",
			input: InputUser{
				Username:        "validuser123",
				Password:        "ValidPassword!2",
				PasswordConfirm: "ValidPassword!2",
				Email:           "valid@email.com",
			},
			expectedError: false,
			expectedUser: InputUser{
				Username:        "validuser123",
				Password:        "ValidPassword!2",
				PasswordConfirm: "ValidPassword!2",
				Email:           "valid@email.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validatedUser, err := validateRegisterInput(tc.input)

			if tc.expectedError {
				assert.Error(t, err)
				// Check if validatedUser is an empty struct
				assert.Equal(t, InputUser{}, validatedUser)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, validatedUser)
			}
		})
	}
}
