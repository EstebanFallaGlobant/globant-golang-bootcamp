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

type mockService struct {
	mock.Mock
}

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
				return CreateUserRequest{
					AuthToken: "Test token",
					User:      NewUser("Test user", "Test password", 10, 0),
				}
			},
			srvUser: getNewUser(),
			checkResponse: func(response interface{}, err error, t *testing.T) {
				if res, ok := response.(CreateUserResponse); !ok {
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
				return CreateUserRequest{
					AuthToken: "Test token",
					User:      NewUser("Test user", "Test password", 10, 0),
				}
			},
			srvUser:       getNewUser(),
			expectedError: errors.New("Some error"),
			checkResponse: func(response interface{}, err error, t *testing.T) {
				if res, ok := response.(CreateUserResponse); !ok {
					t.Fatal(errors.New("response could not be parsed"))
				} else {
					assert.Error(t, err)
					assert.EqualValues(t, 0, res.Id)
				}

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

			svr := new(mockService)

			svr.On("CreateUser", ctx, tc.srvUser).Return(tc.expectedId, tc.expectedError)

			endpoints := MakeEndpoints(svr, logger, nil)

			response, err := endpoints.GetCreateUser(ctx, tc.getRequest())

			tc.checkResponse(response, err, t)
		})
	}
}

func (mock *mockService) CreateUser(ctx context.Context, user User) (int64, error) {
	args := mock.Called(ctx, user)

	return int64(args.Int(0)), args.Error(1)
}

func getNewUser() User {
	return NewUser("Test user", "Test password", 10, 0)
}
