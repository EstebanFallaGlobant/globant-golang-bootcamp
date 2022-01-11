package user

import (
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServiceErrorHandler struct {
	logger kitlog.Logger
}

func (errHandler *userServiceErrorHandler) TogRPCStatus(err error) error {
	level.Info(errHandler.logger).Log("status", "checking error type")
	var result error
	switch err.(type) {
	case invalidRequestError:
		result = status.Error(codes.InvalidArgument, err.Error())
	case argumentsRequiredError:
		result = status.Error(codes.InvalidArgument, err.Error())
	case invalidArgumentsError:
		result = status.Error(codes.InvalidArgument, err.Error())
	case userNotFoundError:
		result = status.Error(codes.NotFound, err.Error())
	case userAlreadyExistsError:
		result = status.Error(codes.AlreadyExists, "user already exists")
	case userNotUpdatedError:
		result = status.Error(codes.Unavailable, "user not updated")
	default:
		result = status.Error(codes.Unavailable, err.Error())
	}

	return result
}
