package user

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateUserTransport(t *testing.T) {
	testCases := []struct {
		name          string
		createRequest func() *pb.CreateUserRequest
		createSvcUser func() User
		checkResponse func(t *testing.T, response interface{}, err error)
		svcId         int
		svcError      error
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
			createSvcUser: getNewUser,
			createRequest: func() *pb.CreateUserRequest {
				return &pb.CreateUserRequest{
					AuthToken: "Test token",
					User:      getNewgRPCUser(),
				}
			},
			checkResponse: func(t *testing.T, response interface{}, err error) {
				if resp, ok := response.(*pb.CreateUserResponse); !ok {
					t.Fatal("Response couldn't be parsed")
				} else {
					assert.EqualValues(t, 0, resp.GetId())
				}
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

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mockService)

			level.Info(logger).Log("svcId", tc.svcId)
			level.Info(logger).Log("svcError", tc.svcError)

			svc.On("CreateUser", mock.Anything, tc.createSvcUser()).Return(tc.svcId, tc.svcError)

			endpoints := MakeEndpoints(svc, logger, nil)
			s := NewgRPCServer(endpoints, logger)

			resp, err := s.CreateUser(ctx, tc.createRequest())

			tc.checkResponse(t, resp, err)
		})
	}
}

func getNewgRPCUser() *pb.User {
	user := getNewUser()

	return &pb.User{
		Name:    user.Name,
		PwdHash: user.PwdHash,
		Age:     uint32(user.Age),
		Parent:  user.Parent,
	}
}
