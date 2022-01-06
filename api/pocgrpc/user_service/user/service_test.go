package user

import (
	"context"
	"os"
	"testing"

	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_Service_CreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		userData      User
		expectedId    int64
		expectedError error
		checkResult   func(t *testing.T, expected, response int64, err error)
	}{
		{
			name: "Test user creation successful",
			userData: User{
				Name:    "Test User",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(34),
				Parent:  0},
			expectedId:    1,
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, response)
			},
		},
		{
			name: "Test user creation with empty user name",
			userData: User{
				Name:    " ",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(20),
				Parent:  0},
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, err, NewArgumentsRequiredError("name"))
			},
		},
		{
			name: "Test user creation with empty password",
			userData: User{
				Name:    "Test user",
				PwdHash: " ",
				Age:     uint8(20),
				Parent:  0},
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, err, NewArgumentsRequiredError("password"))
			},
		},
		{
			name: "Test user creation with empty user name and password",
			userData: User{
				Name:    " ",
				PwdHash: " ",
				Age:     uint8(20),
				Parent:  0},
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, err.Error(), NewArgumentsRequiredError("name", "password").Error())
			},
		},
		{
			name: "Test user creation with age 0",
			userData: User{
				Name:    "Test user",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(0),
				Parent:  0},
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, NewInvalidArgumentsError("age", "between 1 and 150"), err)
			},
		},
		{
			name: "Test user creation with age 151",
			userData: User{
				Name:    "Test user",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(151),
				Parent:  0},
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, NewInvalidArgumentsError("age", "between 1 and 150"), err)
			},
		},
		{
			name: "Test user creation failed by repository message",
			userData: User{
				Name:    "Test user",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(25),
				Parent:  0},
			expectedError: getGenericRepositoryError(),
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, getGenericRepositoryError(), err)
			},
		},
	}

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			repo := new(connectionMock)

			repo.On(
				"InsertUser",
				tc.userData.Name,
				tc.userData.PwdHash,
				tc.userData.Age,
				tc.userData.Parent).
				Return(
					int(tc.expectedId),
					tc.expectedError)

			service := NewService(repo, logger)

			id, err := service.CreateUser(ctx, tc.userData)

			tc.checkResult(t, tc.expectedId, id, err)
		})
	}
}

func Test_GetUser(t *testing.T) {
	testCases := []struct {
		name            string
		searchId        int64
		expectedUser    User
		repositoryError error
		checkResult     func(t *testing.T, expected, response User, err error)
	}{
		{
			name:     "Test successfully retrieved user",
			searchId: 3,
			expectedUser: User{
				Id:      3,
				Name:    "Test User",
				PwdHash: "Test Password",
				Age:     uint8(25),
			},
			checkResult: func(t *testing.T, expected, response User, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, expected, response)
			},
		},
		{
			name:         "Test user retrieve failed with id 0",
			searchId:     0,
			expectedUser: User{},
			checkResult: func(t *testing.T, expected, response User, err error) {
				assert.Error(t, err)
				assert.IsType(t, NewInvalidArgumentsError("id", "must be 1 or greater"), err)
				assert.EqualValues(t, expected, response)
			},
		},
		{
			name:            "Test user retrieve failed, error on repository",
			searchId:        999,
			expectedUser:    User{},
			repositoryError: getGenericRepositoryError(),
			checkResult: func(t *testing.T, expected, response User, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, expected, response)
			},
		},
	}

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(connectionMock)

			repo.On(
				"GetUser",
				tc.searchId).
				Return(
					tc.expectedUser,
					tc.repositoryError)

			service := NewService(repo, logger)

			user, err := service.GetUser(ctx, tc.searchId)

			tc.checkResult(t, tc.expectedUser, user, err)
		})
	}
}
