package user

import (
	"context"
	"errors"
	"os"
	"testing"

	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type connectionMock struct {
	mock.Mock
}

func (connection *connectionMock) InsertUser(user User) (int64, error) {
	args := connection.Called(user.Name, user.PwdHash, user.Age, user.Parent)
	return int64(args.Int(0)), args.Error(1)
}

func Test_Service_CreateUser(t *testing.T) {
	TestCases := []struct {
		name          string
		userData      User
		expectedId    int64
		expectedError error
		checkResult   func(t *testing.T, expected, response int64, err error)
	}{
		{
			name: "Test user creation successful",
			userData: NewUser(
				"Test User",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(34),
				0),
			expectedId:    1,
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, response)
			},
		},
		{
			name: "Test user creation with empty user name",
			userData: NewUser(
				" ",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(20),
				0),
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, err, NewArgumentsRequiredError("name"))
			},
		},
		{
			name: "Test user creation with empty password",
			userData: NewUser(
				"Test user",
				" ",
				uint8(20),
				0),
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, err, NewArgumentsRequiredError("password"))
			},
		},
		{
			name: "Test user creation with empty user name and password",
			userData: NewUser(
				" ",
				" ",
				uint8(20),
				0),
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, err.Error(), NewArgumentsRequiredError("name", "password").Error())
			},
		},
		{
			name: "Test user creation with age 0",
			userData: NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(0),
				0),
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, NewInvalidArgumentsError("age", "between 1 and 150"), err)
			},
		},
		{
			name: "Test user creation with age 151",
			userData: NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(151),
				0),
			expectedError: nil,
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, NewInvalidArgumentsError("age", "between 1 and 150"), err)
			},
		},
		{
			name: "Test user creation failed by repository message",
			userData: NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(25),
				0),
			expectedError: getGenericRepositoryError(),
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, getGenericRepositoryError(), err)
			},
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var logger kitlog.Logger
			{
				logger = kitlog.NewLogfmtLogger(os.Stderr)
				logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
				logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
			}
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

func getGenericRepositoryError() error {
	return errors.New("generic repository error")
}
