package user

import (
	"errors"
	"os"
	"testing"

	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_TogRPCError(t *testing.T) {
	testCases := []struct {
		name          string
		sentError     error
		expectedError error
	}{
		{
			name:          "Test generic error",
			sentError:     errors.New("Test error"),
			expectedError: status.Error(codes.Unavailable, ""),
		},
		{
			name:          "Test invalid request error",
			sentError:     NewInvalidRequestError(),
			expectedError: status.Error(codes.InvalidArgument, "request is invalid"),
		},
		{
			name:          "Test invalid argument error",
			sentError:     NewInvalidArgumentsError("test argument", "test rule"),
			expectedError: status.Error(codes.InvalidArgument, ""),
		},
		{
			name:          "Test required argument error",
			sentError:     NewArgumentsRequiredError("Test argument"),
			expectedError: status.Error(codes.InvalidArgument, ""),
		},
		{
			name:          "Test user nor found error",
			sentError:     NewUserNotFoundError("Test user", 1),
			expectedError: status.Error(codes.NotFound, ""),
		},
		{
			name:          "Test user already exists error",
			sentError:     NewUserAlreadyExistError("Test user", 1),
			expectedError: status.Error(codes.AlreadyExists, ""),
		},
		{
			name: "Test user not updated error",
			sentError: userNotUpdatedError{
				userId:  100,
				message: "not updated ",
			},
			expectedError: status.Error(codes.Unavailable, ""),
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
			errorHandler := userServiceErrorHandler{
				logger: logger,
			}

			err := errorHandler.TogRPCStatus(tc.sentError)

			assert.ErrorAs(t, err, &tc.expectedError)
		})
	}
}
