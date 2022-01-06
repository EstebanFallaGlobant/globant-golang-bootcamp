package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type service interface {
	CreateUser(ctx context.Context, user User) (int64, error)
	GetUser(ctx context.Context, id int64) (User, error)
}

type Endpoints struct {
	GetCreateUser endpoint.Endpoint
	GetGetUser    endpoint.Endpoint
}

func MakeEndpoints(svc service, logger kitlog.Logger, middlewares []endpoint.Middleware) Endpoints {
	return Endpoints{
		GetCreateUser: wrapEndpoint(makeGetCreateUserEndpoint(svc, logger), middlewares),
		GetGetUser:    wrapEndpoint(makeGetGetUserEndpoint(svc, logger), middlewares),
	}
}

func makeGetCreateUserEndpoint(svc service, logger kitlog.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(createUserRequest)

		if !ok {
			level.Error(logger).Log("message", "invalid request")
			err := NewInvalidRequestError()
			return createUserResponse{status: err}, err
		}

		id, err := svc.CreateUser(ctx, req.user)

		if err != nil {
			level.Error(logger).Log("message", "error creating user")
			return createUserResponse{status: err}, errors.New("error creating user")
		}

		return createUserResponse{
			Id: id,
		}, nil
	}
}

func makeGetGetUserEndpoint(svc service, logger kitlog.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
		user, err := svc.GetUser(ctx, req.id)
		return getUserResponse{
			status: err,
			user:   user,
		}, err
	}
}

func wrapEndpoint(endpoint endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	for _, middleware := range middlewares {
		endpoint = middleware(endpoint)
	}
	return endpoint
}
