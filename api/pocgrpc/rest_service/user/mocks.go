package user

import (
	"context"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/mock"
)

type mockSvc struct {
	mock.Mock
}

func (mock mockSvc) GetUser(id int64) (entities.User, error) {
	args := mock.Called(id)
	return args.Get(0).(entities.User), args.Error(1)
}

type mockEndpoints struct {
	// mock.Mock
	svc    userService
	logger kitlog.Logger
}

func (mock mockEndpoints) MakeGetUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)

		if !ok {
			return nil, svcerr.NewInvalidRequestError("")
		}

		if err := req.Validate(); err != nil {
			return nil, err
		}

		return mock.svc.GetUser(req.userID)
	}
}
