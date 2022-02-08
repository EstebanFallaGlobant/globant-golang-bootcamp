package client

import (
	context "context"
	"errors"
	"testing"
	"time"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/transform"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetUser(t *testing.T) {
	testCases := []struct {
		name          string
		ID            int64
		serviceError  error
		expectedToken interface{}
		expectedUser  entities.User
		prepareMock   func(ctrl *gomock.Controller, pbUser *pb.User, expectedError error) *MockUserDetailServiceClient
		checkResult   func(t *testing.T, gotUser, expectedUser entities.User, err error)
	}{
		{
			name:          "Test successful get user call",
			ID:            genericID,
			expectedToken: genericToken,
			expectedUser: entities.User{
				ID:       genericID,
				Name:     genericUsrName,
				Password: genericPassword,
				Age:      genericUsrAge,
			},
			prepareMock: func(ctrl *gomock.Controller, pbUser *pb.User, expectedError error) *MockUserDetailServiceClient {
				mock := NewMockUserDetailServiceClient(ctrl)
				{
					mock.EXPECT().
						GetUser(gomock.Any(), gomock.Any(), gomock.Any()).
						DoAndReturn(func(_, _, _ interface{}) (*pb.GetUserResponse, error) {
							time.Sleep(genericSleepTime)
							return &pb.GetUserResponse{
								User: pbUser,
							}, expectedError
						}).
						AnyTimes()
				}

				return mock
			},
			checkResult: func(t *testing.T, gotUser, expectedUser entities.User, err error) {
				assert.NoError(t, err)
				assert.EqualValues(t, expectedUser, gotUser)
			},
		},
		{
			name:          "Test failed call, invalid auth token",
			ID:            genericID,
			expectedToken: 32,
			prepareMock: func(ctrl *gomock.Controller, pbUser *pb.User, expectedError error) *MockUserDetailServiceClient {
				mock := NewMockUserDetailServiceClient(ctrl)
				return mock
			},
			checkResult: func(t *testing.T, gotUser, expectedUser entities.User, err error) {
				expectedError := svcerr.NewInvalidArgumentError(paramAuthTokenName, ruleAuthTokenInvalidType)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name:          "Test failed call, error from gRPC service",
			ID:            genericID,
			expectedToken: genericToken,
			serviceError:  errors.New("some service error"),
			prepareMock: func(ctrl *gomock.Controller, pbUser *pb.User, expectedError error) *MockUserDetailServiceClient {
				mock := NewMockUserDetailServiceClient(ctrl)
				{
					mock.EXPECT().
						GetUser(gomock.Any(), gomock.Any(), gomock.Any()).
						DoAndReturn(func(_, _, _ interface{}) (*pb.GetUserResponse, error) {
							return &pb.GetUserResponse{}, expectedError
						})
				}
				return mock
			},
			checkResult: func(t *testing.T, gotUser, expectedUser entities.User, err error) {
				expectedError := errors.New("")
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
	}

	ctrl := gomock.NewController(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), authTokenCtxKey, tc.expectedToken)

			pbUser := transform.FromUserToPbUser(tc.expectedUser)

			gRPCClient := tc.prepareMock(ctrl, &pbUser, tc.serviceError)

			client := userServiceClient{gRPCClient: gRPCClient}

			result, err := client.GetUser(ctx, tc.ID)

			tc.checkResult(t, result, tc.expectedUser, err)
		})
	}
}
