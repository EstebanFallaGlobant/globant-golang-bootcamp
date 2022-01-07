package user

import (
	"context"
	"errors"
	"os"
	"testing"

	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUserEndpoint(t *testing.T) {
	TestCases := []struct {
		name          string
		getRequest    func() interface{}
		checkResponse func(res interface{}, err error, t *testing.T)
		srvUser       User
		expectedId    int
		expectedError error
	}{
		{
			name: "Test user creation good request",
			getRequest: func() interface{} {
				return createUserRequest{
					authToken: "Test token",
					user:      getNewUser(),
				}
			},
			srvUser: getNewUser(),
			checkResponse: func(response interface{}, err error, t *testing.T) {
				if res, ok := response.(createUserResponse); !ok {
					t.Fatal(errors.New("response could not be parsed"))
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, res)
					assert.EqualValues(t, 1, res.Id)
				}
			},
			expectedId: 1,
		},
		{
			name: "Test bad request",
			getRequest: func() interface{} {
				request := struct {
					Auth string
					Id   int64
				}{
					Auth: "",
					Id:   1,
				}
				return request
			},
			checkResponse: func(res interface{}, err error, t *testing.T) {
				assert.Error(t, err)
				assert.ErrorAs(t, NewInvalidRequestError(), &err)
			},
			srvUser: getNewUser(),
		},
		{
			name: "Test service error",
			getRequest: func() interface{} {
				return createUserRequest{
					authToken: "Test token",
					user:      getNewUser(),
				}
			},
			srvUser:       getNewUser(),
			expectedError: errors.New("Some error"),
			checkResponse: func(response interface{}, err error, t *testing.T) {
				if res, ok := response.(createUserResponse); !ok {
					t.Fatal(errors.New("response could not be parsed"))
				} else {
					assert.Error(t, err)
					assert.EqualValues(t, 0, res.Id)
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

	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {

			svr := new(mockService)

			svr.On("CreateUser", ctx, tc.srvUser).Return(tc.expectedId, tc.expectedError)

			endpoints := MakeEndpoints(svr, logger, nil)

			response, err := endpoints.GetCreateUser(ctx, tc.getRequest())

			tc.checkResponse(response, err, t)
		})
	}
}

func Test_GetUserEndpoint(t *testing.T) {
	testCases := []struct {
		name          string
		getRequest    func() interface{}
		checkResponse func(t *testing.T, response interface{}, expectedUser User, err error)
		userId        int64
		expectedUser  User
		expectedError error
	}{
		{
			name: "Test succsessful endpoint",
			getRequest: func() interface{} {
				return getUserRequest{
					authToken: "Test Tokent",
					id:        2,
				}
			},
			userId: 2,
			expectedUser: User{
				Id:      2,
				Name:    "Test User",
				PwdHash: "Test Password",
				Age:     20,
			},
			checkResponse: func(t *testing.T, response interface{}, expectedUser User, err error) {
				if res, ok := response.(getUserResponse); !ok {
					t.Fatal("response coul not be parsed")
				} else {
					assert.NoError(t, err)
					assert.NoError(t, res.status)
					assert.EqualValues(t, expectedUser, res.user)
				}

			},
		},
		{
			name: "Test invalid request",
			getRequest: func() interface{} {
				request := struct {
					authToken string
					date      string
				}{
					authToken: "Test Tokent",
					date:      "20-02-1998",
				}
				return request
			},
			checkResponse: func(t *testing.T, response interface{}, expectedUser User, err error) {
				if res, ok := response.(getUserResponse); !ok {
					t.Fatal("response could not be parsed")
				} else {
					assert.Error(t, err)
					assert.Error(t, res.status)
					assert.EqualValues(t, expectedUser, res.user)
					assert.IsType(t, NewInvalidRequestError(), err)
				}

			},
		},
		{
			name: "Test service error",
			getRequest: func() interface{} {
				return getUserRequest{
					authToken: "Test Token",
					id:        0,
				}
			},
			expectedError: errors.New("Service error"),
			checkResponse: func(t *testing.T, response interface{}, expectedUser User, err error) {
				if res, ok := response.(getUserResponse); !ok {
					t.Fatal("response could not be parsed")
				} else {
					assert.Error(t, err)
					assert.Error(t, res.status)
					assert.IsType(t, errors.New("Service error"), err)
					assert.EqualValues(t, expectedUser, res.user)
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mockService)

			svc.On("GetUser", ctx, tc.userId).Return(tc.expectedUser, tc.expectedError)

			endpoints := MakeEndpoints(svc, logger, nil)

			response, err := endpoints.GetGetUser(ctx, tc.getRequest())

			tc.checkResponse(t, response, tc.expectedUser, err)
		})
	}
}
