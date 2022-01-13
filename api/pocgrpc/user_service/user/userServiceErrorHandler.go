package user

import (
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceErrorHandler struct {
	Logger kitlog.Logger
}

func (errHandler UserServiceErrorHandler) TogRPCStatus(err error) error {
	level.Info(errHandler.Logger).Log("status", "checking error type")
	var result error
	switch err.(type) {
	case svcerr.InvalidRequestError:
		result = status.Error(codes.InvalidArgument, err.Error())
	case svcerr.ArgumentsRequiredError:
		result = status.Error(codes.InvalidArgument, err.Error())
	case svcerr.InvalidArgumentsError:
		result = status.Error(codes.InvalidArgument, err.Error())
	case svcerr.UserNotFoundError:
		result = status.Error(codes.NotFound, err.Error())
	case svcerr.UserAlreadyExistsError:
		result = status.Error(codes.AlreadyExists, "user already exists")
	case svcerr.UserNotUpdatedError:
		result = status.Error(codes.Unavailable, "user not updated")
	default:
		result = status.Error(codes.Unavailable, err.Error())
	}

	return result
}
