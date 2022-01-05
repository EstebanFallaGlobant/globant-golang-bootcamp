package user

import (
	"context"
	"fmt"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/go-kit/kit/transport/grpc"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCUserInfoService struct {
	logger     kitlog.Logger
	createUser grpc.Handler
	pb.UnimplementedUserDetailServiceServer
}

func NewgRPCServer(endpoints Endpoints, logger kitlog.Logger) pb.UserDetailServiceServer {
	return &gRPCUserInfoService{
		logger: logger,
		createUser: grpc.NewServer(
			endpoints.GetCreateUser,
			makeDecodeCreateUserRequest(logger),
			makeEncodeCreateUserResponse(logger),
		),
	}
}

func (s *gRPCUserInfoService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		level.Error(s.logger).Log("error", err)
		return nil, err
	}
	level.Info(s.logger).Log("message", "Create user from transport layer")
	return resp.(*pb.CreateUserResponse), nil
}

func makeDecodeCreateUserRequest(logger kitlog.Logger) grpc.DecodeRequestFunc {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		if req, ok := request.(*pb.CreateUserRequest); !ok {
			level.Error(logger).Log("error", fmt.Sprintf("request couldn't be parsed: %v", request))
			return nil, status.Error(codes.FailedPrecondition, "request couldn't be parsed")
		} else {
			level.Info(logger).Log("message", "request decoded")
			user := req.User
			return CreateUserRequest{AuthToken: req.AuthToken, User: NewUser(user.Name, user.PwdHash, uint8(user.Age), user.Parent)}, nil
		}
	}
}

func makeEncodeCreateUserResponse(logger kitlog.Logger) grpc.EncodeResponseFunc {
	return func(_ context.Context, response interface{}) (interface{}, error) {
		if resp, ok := response.(CreateUserResponse); !ok {
			level.Error(logger).Log("error", fmt.Sprintf("response could not be parsed: %v", response))
			return nil, status.Error(codes.FailedPrecondition, "response could not be parsed")
		} else {
			level.Info(logger).Log("message", "response encoded")
			return &pb.CreateUserResponse{
				Id: resp.Id,
			}, nil
		}
	}
}
