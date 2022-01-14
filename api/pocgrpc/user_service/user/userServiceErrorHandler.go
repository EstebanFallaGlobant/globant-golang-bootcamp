package user

import (
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/codes"
)

type UserServiceErrorHandler struct {
	Logger kitlog.Logger
}

func (errHandler UserServiceErrorHandler) TogRPCStatus(err error) *pb.Err {
	level.Info(errHandler.Logger).Log("status", "checking error type")
	var result pb.Err
	switch err.(type) {
	case svcerr.InvalidRequestError:
		result.Code = uint32(codes.InvalidArgument)
		result.ErrMsg = err.Error()
	case svcerr.ArgumentsRequiredError:
		result.Code = uint32(codes.InvalidArgument)
		result.ErrMsg = err.Error()
	case svcerr.InvalidArgumentsError:
		result.Code = uint32(codes.InvalidArgument)
		result.ErrMsg = err.Error()
	case svcerr.UserNotFoundError:
		result.Code = uint32(codes.NotFound)
		result.ErrMsg = err.Error()
	case svcerr.UserAlreadyExistsError:
		result.Code = uint32(codes.AlreadyExists)
		result.ErrMsg = err.Error()
	case svcerr.UserNotUpdatedError:
		result.Code = uint32(codes.Unavailable)
		result.ErrMsg = err.Error()
	default:
		result.Code = uint32(codes.Unavailable)
		result.ErrMsg = err.Error()
	}

	return &result
}
