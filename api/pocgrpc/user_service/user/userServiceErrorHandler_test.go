package user

import (
	"errors"
	"os"
	"testing"

	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func Test_TogRPCError(t *testing.T) {
	testCases := []struct {
		name              string
		sentError         error
		expectedErrorCode uint32
	}{
		{
			name:              "Test generic error",
			sentError:         errors.New("Test error"),
			expectedErrorCode: uint32(codes.Unavailable),
		},
		{
			name:              "Test invalid request error",
			sentError:         svcerr.NewInvalidRequestError(""),
			expectedErrorCode: uint32(codes.InvalidArgument),
		},
		{
			name:              "Test invalid argument error",
			sentError:         svcerr.NewInvalidArgumentsError("test argument", "test rule"),
			expectedErrorCode: uint32(codes.InvalidArgument),
		},
		{
			name:              "Test required argument error",
			sentError:         svcerr.NewArgumentsRequiredError("Test argument"),
			expectedErrorCode: uint32(codes.InvalidArgument),
		},
		{
			name:              "Test user nor found error",
			sentError:         svcerr.NewUserNotFoundError("Test user", 1),
			expectedErrorCode: uint32(codes.NotFound),
		},
		{
			name:              "Test user already exists error",
			sentError:         svcerr.NewUserAlreadyExistError("Test user", 1),
			expectedErrorCode: uint32(codes.AlreadyExists),
		},
		{
			name:              "Test user not updated error",
			sentError:         svcerr.NewUserNotUpdatedError(10, "user not updated"),
			expectedErrorCode: uint32(codes.Unavailable),
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
			errorHandler := UserServiceErrorHandler{
				Logger: logger,
			}

			err := errorHandler.TogRPCStatus(tc.sentError)

			assert.EqualValues(t, tc.expectedErrorCode, err.Code)
		})
	}
}
