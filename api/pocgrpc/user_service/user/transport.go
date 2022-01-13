package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/go-kit/kit/transport/grpc"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type gRPCUserInfoService struct {
	errHandler errorHandler
	logger     kitlog.Logger
	createUser grpc.Handler
	getUser    grpc.Handler
	pb.UnimplementedUserDetailServiceServer
}

type errorHandler interface {
	TogRPCStatus(err error) error
}

func NewgRPCServer(endpoints Endpoints, logger kitlog.Logger, errorHandler errorHandler) pb.UserDetailServiceServer {
	return &gRPCUserInfoService{
		logger:     logger,
		errHandler: errorHandler,
		createUser: grpc.NewServer(
			endpoints.GetCreateUser,
			makeDecodeCreateUserRequest(logger, errorHandler),
			makeEncodeCreateUserResponse(logger, errorHandler),
		),
		getUser: grpc.NewServer(
			endpoints.GetGetUser,
			makeDecodeGetUserRequest(logger, errorHandler),
			makeEncodeGetUserResponse(logger, errorHandler),
		),
	}
}

func (s *gRPCUserInfoService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(s.logger).Log("failed", err)
		return nil, s.errHandler.TogRPCStatus(err)
	}
	return resp.(*pb.CreateUserResponse), nil
}

func (s *gRPCUserInfoService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	_, resp, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(s.logger).Log("failed", err)
		return nil, s.errHandler.TogRPCStatus(err)
	}
	return resp.(*pb.GetUserResponse), nil
}

func makeDecodeCreateUserRequest(logger kitlog.Logger, errorHandler errorHandler) grpc.DecodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		level.Info(logger).Log("status", "decoding request")
		req, ok := request.(*pb.CreateUserRequest)
		if !ok {
			level.Error(logger).Log("request couldn't be parsed", request)
			return nil, errorHandler.TogRPCStatus(svcerr.NewInvalidRequestError("request could not be parsed"))
		}
		level.Info(logger).Log("status", "request decoded")
		user, err := entities.NewUser(req.User.Name, req.User.PwdHash, uint8(req.User.Age), req.User.Parent)

		if err != nil {
			level.Error(logger).Log("error creating new user", err)
			return nil, errorHandler.TogRPCStatus(svcerr.NewInvalidRequestError(err.Error()))
		}

		return createUserRequest{authToken: req.AuthToken, user: user}, nil

	}
}

func makeEncodeCreateUserResponse(logger kitlog.Logger, errorHandler errorHandler) grpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		res, ok := response.(createUserResponse)

		if !ok {
			level.Error(logger).Log("error", fmt.Sprintf("response could not be parsed: %v", response))
			return nil, errorHandler.TogRPCStatus(errors.New("response could not be parsed"))
		}

		if res.status != nil {
			level.Error(logger).Log("error", res.status)
			return nil, errorHandler.TogRPCStatus(res.status)
		}

		level.Info(logger).Log("message", "response encoded")

		return &pb.CreateUserResponse{
			Id: res.Id,
		}, nil

	}
}

func makeDecodeGetUserRequest(logger kitlog.Logger, errorHandler errorHandler) grpc.DecodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		level.Info(logger).Log("status", "decoding request")
		req, ok := request.(*pb.GetUserRequest)

		if !ok {
			level.Error(logger).Log("request could not be parsed", request)
			return nil, errorHandler.TogRPCStatus(errors.New("request could not be parsed"))
		}

		level.Info(logger).Log("request decoded", req)

		return getUserRequest{
			authToken: req.AuthToken,
			id:        req.Id,
		}, nil
	}
}

func makeEncodeGetUserResponse(logger kitlog.Logger, errorHandler errorHandler) grpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		res, ok := response.(getUserResponse)
		if !ok {
			level.Error(logger).Log("error", fmt.Sprintf("response could not be parsed: %v", response))
			return nil, errorHandler.TogRPCStatus(errors.New("response could not be parsed"))
		}
		if res.status != nil {
			level.Error(logger).Log("error", res.status)
			return nil, errorHandler.TogRPCStatus(res.status)
		}

		return &pb.GetUserResponse{
			User: &pb.User{
				Id:      res.user.ID,
				Name:    res.user.Name,
				PwdHash: res.user.PwdHash,
				Age:     uint32(res.user.Age),
				Parent:  res.user.ParentID,
			},
		}, nil
	}
}
