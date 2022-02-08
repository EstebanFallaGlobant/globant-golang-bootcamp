package entities

import (
	"testing"

	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/stretchr/testify/assert"
)

func Test_UserValidate_StringParams(t *testing.T) {
	testCases := []struct {
		name        string
		testUser    User
		checkResult func(result []error, t *testing.T)
	}{
		{
			name: "Test valid user",
			testUser: User{
				Name:     testUsrName,
				Password: testUsrName,
				Age:      testUsrAge,
			},
			checkResult: func(result []error, t *testing.T) {
				assert.Empty(t, result)
			},
		},
		{
			name: "Test empty name, fails validation",
			testUser: User{
				Name:     " ",
				Password: testPassword,
				Age:      testUsrAge,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedError := svcerr.NewInvalidArgumentError(paramNameStr, ruleEmptyStr)

				assert.EqualValues(t, 1, len(result))
				assert.Error(t, result[0])
				assert.ErrorAs(t, result[0], &expectedError)
			},
		},
		{
			name: "Test empty password, fails validation",
			testUser: User{
				Name:     testUsrName,
				Password: " ",
				Age:      testUsrAge,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedError := svcerr.NewInvalidArgumentError(paramPasswordStr, ruleEmptyStr)

				assert.EqualValues(t, 1, len(result))
				assert.Error(t, result[0])
				assert.ErrorAs(t, result[0], &expectedError)
			},
		},
		{
			name: "Test empty name and password, fails validation",
			testUser: User{
				Age: testUsrAge,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedErrors := []error{
					svcerr.NewInvalidArgumentError(paramNameStr, ruleEmptyStr),
					svcerr.NewInvalidArgumentError(paramPasswordStr, ruleEmptyStr)}

				assert.EqualValues(t, 2, len(result))
				assert.ElementsMatch(t, expectedErrors, result)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.testUser.Validate()

			tc.checkResult(result, t)
		})
	}
}

func Test_UserValidate_NumericParams(t *testing.T) {
	testCases := []struct {
		name        string
		testUser    User
		checkResult func(result []error, t *testing.T)
	}{
		{
			name: "Test age less than 1, fails validation",
			testUser: User{
				Name:     testUsrName,
				Password: testPassword,
				Age:      0,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedError := svcerr.NewInvalidArgumentError(paramAgeStr, ruleLessThanOne)
				assert.GreaterOrEqual(t, 1, len(result))
				assert.ErrorAs(t, result[0], &expectedError)
			},
		},
		{
			name: "Test age greater than allowed, fails validation",
			testUser: User{
				Name:     testUsrName,
				Password: testPassword,
				Age:      MaxAllowedAge + 1,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedError := svcerr.NewInvalidArgumentError(paramAgeStr, ruleGreaterThanAllowedAge)
				assert.GreaterOrEqual(t, 1, len(result))
				assert.ErrorAs(t, result[0], &expectedError)
			},
		},
		{
			name: "Test ID less than 0, fails verification",
			testUser: User{
				Name:     testUsrName,
				Password: testPassword,
				Age:      testUsrAge,
				ID:       -1,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedError := svcerr.NewInvalidArgumentError(paramIDStr, ruleLessThanZero)
				assert.GreaterOrEqual(t, 1, len(result))
				assert.ErrorAs(t, result[0], &expectedError)
			},
		},
		{
			name: "Test parentID less than 0, fails verification",
			testUser: User{
				Name:     testUsrName,
				Password: testPassword,
				Age:      testUsrAge,
				ParentID: -1,
			},
			checkResult: func(result []error, t *testing.T) {
				expectedError := svcerr.NewInvalidArgumentError(paramParentIDStr, ruleLessThanZero)
				assert.GreaterOrEqual(t, 1, len(result))
				assert.ErrorAs(t, result[0], &expectedError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.testUser.Validate()

			tc.checkResult(result, t)
		})
	}
}
