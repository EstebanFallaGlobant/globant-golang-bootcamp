package user

import (
	"context"
	"errors"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/go-kit/kit/transport/grpc"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCUserInfoService struct {
	errHandler errorHandler
	logger     kitlog.Logger
	createUser grpc.Handler
	getUser    grpc.Handler
	pb.UnimplementedUserDetailServiceServer
}

type errorHandler interface {
	TogRPCStatus(err error) *pb.Err
}

func NewgRPCServer(endpoints Endpoints, logger kitlog.Logger, errorHandler errorHandler) pb.UserDetailServiceServer {
	return &gRPCUserInfoService{
		logger:     logger,
		errHandler: errorHandler,
		createUser: grpc.NewServer(
			endpoints.GetCreateUser,
			makeDecodeCreateUserRequest(logger),
			makeEncodeCreateUserResponse(logger),
		),
		getUser: grpc.NewServer(
			endpoints.GetGetUser,
			makeDecodeGetUserRequest(logger),
			makeEncodeGetUserResponse(logger),
		),
	}
}

func (s *gRPCUserInfoService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(s.logger).Log(failStatusKey, err)
		responseError := s.errHandler.TogRPCStatus(err)
		return &pb.CreateUserResponse{
			Status: responseError,
		}, status.Error(codes.Code(responseError.Code), responseError.ErrMsg)
	}
	return resp.(*pb.CreateUserResponse), nil
}

func (s *gRPCUserInfoService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	_, resp, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(s.logger).Log(failStatusKey, err)
		responseError := s.errHandler.TogRPCStatus(err)
		return &pb.GetUserResponse{
			Status: responseError,
		}, status.Error(codes.Code(responseError.Code), responseError.ErrMsg)
	}
	return resp.(*pb.GetUserResponse), nil
}

func makeDecodeCreateUserRequest(logger kitlog.Logger) grpc.DecodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*pb.CreateUserRequest)
		if !ok {
			level.Error(logger).Log(msgRqstNotParsed, request)
			return nil, svcerr.NewInvalidRequestError(msgRqstNotParsed)
		}
		level.Info(logger).Log(nrmStatusKey, msgRqstDecoded)
		user, err := entities.NewUser(req.User.Name, req.User.PwdHash, uint8(req.User.Age), req.User.ParentId)

		if err != nil {
			level.Error(logger).Log(errStatusKey, err)
			return nil, svcerr.NewInvalidRequestError(err.Error())
		}

		return createUserRequest{authToken: req.AuthToken, user: user}, nil

	}
}

func makeEncodeCreateUserResponse(logger kitlog.Logger) grpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		res, ok := response.(createUserResponse)

		if !ok {
			level.Error(logger).Log(msgRspNotParsed, response)
			return nil, errors.New(msgRspNotParsed)
		}

		if res.status != nil {
			level.Error(logger).Log(errStatusKey, res.status)
			return nil, res.status
		}

		level.Info(logger).Log(nrmStatusKey, msgRspEncoded)

		return &pb.CreateUserResponse{
			Id: res.id,
		}, nil

	}
}

func makeDecodeGetUserRequest(logger kitlog.Logger) grpc.DecodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*pb.GetUserRequest)

		if !ok {
			level.Error(logger).Log(msgRqstNotParsed, request)
			return nil, errors.New(msgRqstNotParsed)
		}

		level.Info(logger).Log(nrmStatusKey, msgRqstDecoded)

		return getUserRequest{
			authToken: req.AuthToken,
			id:        req.Id,
		}, nil
	}
}

func makeEncodeGetUserResponse(logger kitlog.Logger) grpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		res, ok := response.(getUserResponse)
		if !ok {
			level.Error(logger).Log(msgRspNotParsed, response)
			return nil, errors.New(msgRspNotParsed)
		}
		if res.status != nil {
			level.Error(logger).Log(errStatusKey, res.status)
			return nil, res.status
		}

		level.Info(logger).Log(nrmStatusKey, msgRspEncoded)

		return &pb.GetUserResponse{
			Status: &pb.Err{
				Code: uint32(codes.OK),
			},
			User: &pb.User{
				Id:       res.user.ID,
				Name:     res.user.Name,
				PwdHash:  res.user.PwdHash,
				Age:      uint32(res.user.Age),
				ParentId: res.user.ParentID,
			},
		}, nil
	}
}
