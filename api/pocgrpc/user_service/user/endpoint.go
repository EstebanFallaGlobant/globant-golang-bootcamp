package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Service interface {
	CreateUser(ctx context.Context, user User) (int64, error)
}

type Endpoints struct {
	GetCreateUser endpoint.Endpoint
}

type CreateUserResponse struct {
	Id int64
}

type CreateUserRequest struct {
	AuthToken string
	User      User
}

func MakeEndpoints(svc Service, logger kitlog.Logger, middlewares []endpoint.Middleware) Endpoints {
	return Endpoints{
		GetCreateUser: wrapEndpoint(makeGetCreateUserEndpoint(svc, logger), middlewares),
	}
}

func makeGetCreateUserEndpoint(svc Service, logger kitlog.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateUserRequest)

		if !ok {
			level.Error(logger).Log("message", "invalid request")
			return nil, NewInvalidRequestError()
		}

		id, err := svc.CreateUser(ctx, req.User)

		if err != nil {
			level.Error(logger).Log("message", "error creating user")
			return nil, errors.New("error creating user")
		}

		return &CreateUserResponse{
			Id: id,
		}, nil
	}
}

func wrapEndpoint(endpoint endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	for _, middleware := range middlewares {
		endpoint = middleware(endpoint)
	}
	return endpoint
}
