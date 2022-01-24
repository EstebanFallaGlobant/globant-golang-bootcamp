package entities

import (
	"testing"

	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/stretchr/testify/assert"
)

const (
	ageValidationErrorMsg = "must be between 1 and 150"
)

func Test_UserValidate(t *testing.T) {
	testCases := []struct {
		name       string
		user       User
		checkError func(t *testing.T, err error)
	}{
		{
			name: "Test user with successful validation",
			user: User{
				ID:      1,
				Name:    "Test User",
				PwdHash: "Test Password",
				Age:     20,
			},
			checkError: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "Test user with empty name",
			user: User{
				Name:    " ",
				PwdHash: "Test password",
				Age:     20,
			},
			checkError: func(t *testing.T, err error) {
				expectedError := svcerr.NewArgumentsRequiredError("name")
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name: "Test user with empty password hash",
			user: User{
				Name:    "Test user",
				PwdHash: " ",
				Age:     20,
			},
			checkError: func(t *testing.T, err error) {
				expectedError := svcerr.NewArgumentsRequiredError("password")
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name: "Test user with empty name and password",
			user: User{
				Name:    " ",
				PwdHash: " ",
				Age:     20,
			},
			checkError: func(t *testing.T, err error) {
				expectedError := svcerr.NewArgumentsRequiredError("name", "password")
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name: "Test user with age 0",
			user: User{
				Name:    "Test user",
				PwdHash: "Test password",
			},
			checkError: func(t *testing.T, err error) {
				expectedError := svcerr.NewInvalidArgumentsError("age", ageValidationErrorMsg)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name: "Test user with age 151",
			user: User{
				Name:    "Test user",
				PwdHash: "Test password",
				Age:     151,
			},
			checkError: func(t *testing.T, err error) {
				expectedError := svcerr.NewInvalidArgumentsError("age", ageValidationErrorMsg)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name: "Test user with ID -1",
			user: User{
				ID:      -1,
				Name:    "Test user",
				PwdHash: "Test password",
				Age:     20,
			},
			checkError: func(t *testing.T, err error) {
				expectedError := svcerr.NewInvalidArgumentsError("ID", "must be 0 or greater")
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.user.Validate()

			tc.checkError(t, err)
		})
	}
}
