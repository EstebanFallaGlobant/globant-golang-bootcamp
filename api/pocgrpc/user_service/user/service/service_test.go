package service

import (
	"context"
	"errors"
	"os"
	"testing"

	apiErr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/errors"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/repository"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_Service_CreateUser(t *testing.T) {
	TestCases := []struct {
		name             string
		userData         repository.User
		expectedResponse CreateUserResponse
		expectedError    error
		request          func(name, pwd string, age uint8, parent int64) CreateUserRequest
		checkResult      func(t *testing.T, expected *CreateUserResponse, response *CreateUserResponse, err error)
	}{
		{
			name: "Test user creation successful",
			userData: repository.NewUser(
				"Test User",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(34),
				0),
			expectedResponse: CreateUserResponse{Id: 1},
			expectedError:    nil,
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected *CreateUserResponse, response *CreateUserResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected.Id, response.Id)
			},
		},
		{
			name: "Test user creation with empty user name",
			userData: repository.NewUser(
				" ",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(20),
				0),
			expectedResponse: getGenericCreateUserResponse(),
			expectedError:    nil,
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected, response *CreateUserResponse, err error) {
				assert.Error(t, err)
				assert.IsType(t, err, apiErr.NewArgumentsRequiredError("name"))
			},
		},
		{
			name: "Test user creation with empty password",
			userData: repository.NewUser(
				"Test user",
				" ",
				uint8(20),
				0),
			expectedResponse: getGenericCreateUserResponse(),
			expectedError:    nil,
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected, response *CreateUserResponse, err error) {
				assert.Error(t, err)
				assert.IsType(t, err, apiErr.NewArgumentsRequiredError("password"))
			},
		},
		{
			name: "Test user creation with empty user name and password",
			userData: repository.NewUser(
				" ",
				" ",
				uint8(20),
				0),
			expectedResponse: getGenericCreateUserResponse(),
			expectedError:    nil,
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected, response *CreateUserResponse, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, err.Error(), apiErr.NewArgumentsRequiredError("name", "password").Error())
			},
		},
		{
			name: "Test user creation with age 0",
			userData: repository.NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(0),
				0),
			expectedResponse: getGenericCreateUserResponse(),
			expectedError:    nil,
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected, response *CreateUserResponse, err error) {
				assert.Error(t, err)
				assert.IsType(t, apiErr.NewInvalidArgumentsError("age", "between 1 and 150"), err)
			},
		},
		{
			name: "Test user creation with age 151",
			userData: repository.NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(151),
				0),
			expectedResponse: getGenericCreateUserResponse(),
			expectedError:    nil,
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected, response *CreateUserResponse, err error) {
				assert.Error(t, err)
				assert.IsType(t, apiErr.NewInvalidArgumentsError("age", "between 1 and 150"), err)
			},
		},
		{
			name: "Test user creation failed by repository message",
			userData: repository.NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(25),
				0),
			expectedResponse: getGenericCreateUserResponse(),
			expectedError:    getGenericRepositoryError(),
			request:          getGenericCreateUserRequest,
			checkResult: func(t *testing.T, expected, response *CreateUserResponse, err error) {
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
					tc.expectedResponse.Id,
					tc.expectedError)

			service := NewService(repo, logger)

			res, err := service.CreateUser(ctx, tc.request(tc.userData.Name, tc.userData.PwdHash, tc.userData.Age, tc.userData.Parent))

			tc.checkResult(t, &tc.expectedResponse, res, err)
		})
	}
}

func getGenericCreateUserRequest(name, pwd string, age uint8, parent int64) CreateUserRequest {
	return CreateUserRequest{
		Name:   name,
		Pwd:    pwd,
		Age:    age,
		Parent: parent,
	}
}

func getGenericCreateUserResponse() CreateUserResponse {
	return CreateUserResponse{
		Id: 1,
	}
}

func getGenericRepositoryError() error {
	return errors.New("generic repository error")
}
