package user

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_CreateUserTransport(t *testing.T) {
	testCases := []struct {
		name          string
		createRequest func() *pb.CreateUserRequest
		createSvcUser func(...InitializationOption) User
		checkResponse func(t *testing.T, response interface{}, err error)
		svcId         int
		svcError      error
		handlerError  error
	}{
		{
			name:          "Test successful response",
			svcId:         1,
			createSvcUser: getNewUser,
			createRequest: func() *pb.CreateUserRequest {
				return &pb.CreateUserRequest{
					AuthToken: "Test token",
					User:      getNewgRPCUser(),
				}
			},
			checkResponse: func(t *testing.T, response interface{}, err error) {
				if resp, ok := response.(*pb.CreateUserResponse); !ok {
					t.Fatal("Response couln't be parsed")
				} else {
					assert.NoError(t, err)
					assert.EqualValues(t, 1, resp.GetId())
				}
			},
		},
		{
			name:          "Test with error from service",
			svcId:         0,
			svcError:      errors.New("Some service error"),
			handlerError:  status.Error(codes.Unavailable, "test error"),
			createSvcUser: getNewUser,
			createRequest: func() *pb.CreateUserRequest {
				return &pb.CreateUserRequest{
					AuthToken: "Test token",
					User:      getNewgRPCUser(),
				}
			},
			checkResponse: func(t *testing.T, response interface{}, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
			},
		},
	}

	ctx := context.Background()
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mockService)
			svc.On("CreateUser", mock.Anything, tc.createSvcUser()).Return(tc.svcId, tc.svcError)

			errHandler := new(mockErrorHandler)
			errHandler.On("TogRPCStatus", tc.svcError).Return(tc.handlerError)

			server := NewgRPCServer(MakeEndpoints(svc, logger, nil), logger, errHandler)

			resp, err := server.CreateUser(ctx, tc.createRequest())

			tc.checkResponse(t, resp, err)
		})
	}
}

func Test_GetUserTransport(t *testing.T) {
	testCases := []struct {
		name          string
		userId        int64
		svcUser       User
		svcError      error
		checkResponse func(t *testing.T, expectedUser User, response interface{}, err error)
		getRequest    func() *pb.GetUserRequest
		handlerError  error
	}{
		{
			name:   "Tests successful request",
			userId: 2,
			svcUser: getNewUser(func(user *User) error {
				user.Id = 2
				return nil
			}),
			getRequest: func() *pb.GetUserRequest {
				return &pb.GetUserRequest{
					AuthToken: "Test Tokent",
					Id:        2,
				}
			},
			checkResponse: func(t *testing.T, expectedUser User, response interface{}, err error) {
				if resp, ok := response.(*pb.GetUserResponse); !ok {
					t.Fatalf("Failed to parse response: %v", response)
				} else {
					assert.NoError(t, err)
					assert.EqualValues(t, expectedUser.Id, resp.User.Id)
					assert.EqualValues(t, expectedUser.Name, resp.User.Name)
					assert.EqualValues(t, expectedUser.PwdHash, resp.User.PwdHash)
					assert.EqualValues(t, expectedUser.Age, resp.User.Age)
					assert.EqualValues(t, expectedUser.Parent, resp.User.Parent)
				}
			},
		},
		{
			name:         "Test service error",
			userId:       999,
			svcUser:      User{},
			svcError:     errors.New("Service error"),
			handlerError: status.Error(codes.Unavailable, "test error"),
			getRequest: func() *pb.GetUserRequest {
				return &pb.GetUserRequest{
					AuthToken: "Test Token",
					Id:        999,
				}
			},
			checkResponse: func(t *testing.T, expectedUser User, response interface{}, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)
			},
		},
	}

	ctx := context.Background()
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mockService)
			svc.On("GetUser", mock.Anything, tc.userId).Return(tc.svcUser, tc.svcError)

			errHandler := new(mockErrorHandler)
			errHandler.On("TogRPCStatus", tc.svcError).Return(tc.handlerError)

			server := NewgRPCServer(MakeEndpoints(svc, logger, nil), logger, errHandler)

			resp, err := server.GetUser(ctx, tc.getRequest())

			tc.checkResponse(t, tc.svcUser, resp, err)
		})
	}
}
