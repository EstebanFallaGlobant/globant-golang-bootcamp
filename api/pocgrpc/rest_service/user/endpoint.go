package user

import (
	"context"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type userService interface {
	GetUser(id int64) (entities.User, error)
}

type getUserEndpoint struct {
	svc    userService
	logger log.Logger

	getUserMiddlewares []endpoint.Middleware
}

type TargetEndpoint uint8

const (
	Undefined TargetEndpoint = iota
	GetUser
)

func (e *getUserEndpoint) MakeGetUserEndpoint() endpoint.Endpoint {

	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)

		if !ok {
			err := svcerr.NewInvalidRequestError(invRqstNonParsed)
			level.Error(e.logger).Log(err, request)
			return nil, err
		}

		if err := req.Validate(); err != nil {
			level.Error(e.logger).Log(err, req)
			return nil, err
		}

		user, err := e.svc.GetUser(req.userID)

		if err != nil {
			level.Error(e.logger).Log(userNotFoundKeyName, err)
			return nil, err
		}

		return getUserResponse{user: user}, nil
	}

	return wrapEndpoint(endpoint, e.getUserMiddlewares)
}

func (e *getUserEndpoint) AddMiddlewares(target TargetEndpoint, middlewares ...endpoint.Middleware) error {
	switch target {
	case GetUser:
		e.getUserMiddlewares = append(e.getUserMiddlewares, middlewares...)
	default:
		return svcerr.NewInvalidArgumentError(endpointTargetParamName, invArgEndpointTarget)
	}

	return nil
}

func (e *getUserEndpoint) countMiddlewares(target TargetEndpoint) int {
	switch target {
	case GetUser:
		return len(e.getUserMiddlewares)
	default:
		return 0
	}
}

func wrapEndpoint(endpoint endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	for _, middleware := range middlewares {
		endpoint = middleware(endpoint)
	}

	return endpoint
}
