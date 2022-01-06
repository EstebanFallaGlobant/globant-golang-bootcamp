package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/grpcmicroservice/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type ServiceInterface interface {
	IsPal(string) string
	Reverse(string) string
}

type Endpoints struct {
	GetIsPalindrome endpoint.Endpoint
	GetReverse      endpoint.Endpoint
}

func MakeEndpoints(svc ServiceInterface, logger log.Logger, middlewares []endpoint.Middleware) Endpoints {
	return Endpoints{
		GetIsPalindrome: wrapEndpoint(makeGetIsPalindromeEndpoint(svc, logger), middlewares),
		GetReverse:      wrapEndpoint(makeGetReverseEndpoint(svc, logger), middlewares),
	}
}

func makeGetIsPalindromeEndpoint(svc ServiceInterface, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*models.IsPalRequest)

		if !ok {
			level.Error(logger).Log("message", "invalid request")
			return nil, errors.New("invalid request")
		}

		msg := svc.IsPal(req.Word)

		return &models.IsPalResponse{
			Message: msg,
		}, nil
	}
}

func makeGetReverseEndpoint(svc ServiceInterface, logger log.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		level.Info(logger).Log("mensaje", "entra a get reverse endpoint")
		req, ok := request.(*models.ReverseRequest)

		msg, _ := json.Marshal(req)

		level.Info(logger).Log("request", string(msg))

		if !ok {
			level.Error(logger).Log("message", "invalid request")
			return nil, errors.New("invalid request")
		}

		level.Info(logger).Log("mensaje", "pasa if")
		reverseString := svc.Reverse(req.Word)
		level.Info(logger).Log("respuesta", reverseString)
		result := new(models.ReverseResponse)
		result.Word = reverseString
		return result, nil
	}
}

func wrapEndpoint(endpoint endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	for _, middleware := range middlewares {
		endpoint = middleware(endpoint)
	}

	return endpoint
}
