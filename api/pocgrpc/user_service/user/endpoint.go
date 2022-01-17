package user

import (
	"context"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type service interface {
	CreateUser(ctx context.Context, user entities.User) (int64, error)
	GetUser(ctx context.Context, id int64) (entities.User, error)
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
			level.Error(logger).Log(msgRqstNotParsed, request)
			err := svcerr.NewInvalidRequestError(msgRqstNotParsed)
			return createUserResponse{status: err}, err
		}

		if err := req.Validate(); err != nil {
			level.Error(logger).Log(msgRqstInvalid, err)
			err = svcerr.NewInvalidRequestError(err.Error())
			return createUserResponse{status: err}, err
		}

		id, err := svc.CreateUser(ctx, req.user)

		if err != nil {
			level.Error(logger).Log(errStatusKey, err)
			return createUserResponse{status: err}, err
		}

		return createUserResponse{
			id: id,
		}, nil
	}
}

func makeGetGetUserEndpoint(svc service, logger kitlog.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)

		if !ok {
			level.Error(logger).Log(msgRqstNotParsed, request)
			err := svcerr.NewInvalidRequestError(msgRqstNotParsed)
			return getUserResponse{status: err}, err
		}

		if err := req.Validate(); err != nil {
			level.Error(logger).Log(msgRqstInvalid, err)
			err := svcerr.NewInvalidRequestError(err.Error())
			return getUserResponse{status: err}, err
		}

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
