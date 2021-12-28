package service

import (
	"context"
	"os"
	"testing"

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
		request          func(name string, pwd string, age uint8, parent int64) CreateUserRequest
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
			request: func(name, pwd string, age uint8, parent int64) CreateUserRequest {
				return CreateUserRequest{
					Name:   name,
					Pwd:    pwd,
					Age:    age,
					Parent: parent,
				}
			},
			checkResult: func(t *testing.T, expected *CreateUserResponse, response *CreateUserResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected.Id, response.Id)
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
