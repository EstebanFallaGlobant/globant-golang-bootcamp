package user

import (
	"context"
	"os"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_Service_CreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		userData      entities.User
		expectedId    int64
		expectedError error
		configMock    func(mock *connectionMock, returnError error, user entities.User, returnID int64)
		checkResult   func(t *testing.T, expected, response int64, err error)
	}{
		{
			name: "Test user creation successful",
			userData: entities.User{
				Name:     "Test User",
				PwdHash:  "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:      uint8(34),
				ParentID: 0},
			expectedId:    1,
			expectedError: nil,
			configMock: func(mock *connectionMock, returnError error, user entities.User, returnID int64) {
				mock.On(
					"InsertUser",
					user.Name,
					user.PwdHash,
					user.Age,
					user.ParentID).
					Return(
						int(returnID),
						returnError)

				mock.On(
					"GetUserByName",
					user.Name).
					Return(
						entities.User{},
						nil)
			},
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, response)
			},
		},
		{
			name: "Test user creation failed by repository message",
			userData: entities.User{
				Name:     "Test user",
				PwdHash:  "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:      uint8(25),
				ParentID: 0},
			expectedError: getGenericRepositoryError(),
			configMock: func(mock *connectionMock, returnError error, user entities.User, returnID int64) {
				mock.On(
					"InsertUser",
					user.Name,
					user.PwdHash,
					user.Age,
					user.ParentID).
					Return(
						int(returnID),
						returnError)

				mock.On(
					"GetUserByName",
					user.Name).
					Return(
						entities.User{},
						nil)
			},
			checkResult: func(t *testing.T, expected, response int64, err error) {
				assert.Error(t, err)
				assert.IsType(t, getGenericRepositoryError(), err)
			},
		},
		{
			name: "Test user creation failed due to user already existing",
			userData: entities.User{
				Name:    "Test user",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     20,
			},
			expectedError: getGenericRepositoryError(),
			configMock: func(mock *connectionMock, returnError error, user entities.User, returnID int64) {
				returnUser := user
				{
					returnUser.ID = 1
				}

				mock.On(
					"GetUserByName",
					user.Name).
					Return(
						returnUser,
						nil)
			},
			checkResult: func(t *testing.T, expected, response int64, err error) {
				expectedError := svcerr.NewUserAlreadyExistError("Test user", 1)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
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

			tc.configMock(repo, tc.expectedError, tc.userData, tc.expectedId)

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
		expectedUser    entities.User
		repositoryError error
		checkResult     func(t *testing.T, expected, response entities.User, err error)
	}{
		{
			name:     "Test successfully retrieved user",
			searchId: 3,
			expectedUser: entities.User{
				ID:      3,
				Name:    "Test User",
				PwdHash: "Test Password",
				Age:     uint8(25),
			},
			checkResult: func(t *testing.T, expected, response entities.User, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, expected, response)
			},
		},
		{
			name:         "Test user retrieve failed with id 0",
			searchId:     0,
			expectedUser: entities.User{},
			checkResult: func(t *testing.T, expected, response entities.User, err error) {
				assert.Error(t, err)
				assert.IsType(t, svcerr.NewInvalidArgumentsError("ID", "must be 1 or greater"), err)
				assert.EqualValues(t, expected, response)
			},
		},
		{
			name:            "Test user retrieve failed, error on repository",
			searchId:        999,
			expectedUser:    entities.User{},
			repositoryError: getGenericRepositoryError(),
			checkResult: func(t *testing.T, expected, response entities.User, err error) {
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
