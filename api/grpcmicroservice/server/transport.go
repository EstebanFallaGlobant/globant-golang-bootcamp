package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/grpcmicroservice/models"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Transport struct {
	logger log.Logger
}

func NewTransport(logger log.Logger) Transport {
	return Transport{
		logger: logger,
	}
}

func (trs *Transport) GetIsPalHandler(endpoint endpoint.Endpoint, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		endpoint,
		trs.decodeGetIsPalRequest,
		trs.encodeGetIsPalResponse,
		options...,
	)
}

func (trs *Transport) decodeGetIsPalRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.IsPalRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	} else {
		level.Info(trs.logger).Log("parsed request", req.Word)
		return &req, nil
	}
}

func (trs *Transport) encodeGetIsPalResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(*models.IsPalResponse)

	if !ok {
		return errors.New("error decoding")
	}

	return json.NewEncoder(w).Encode(resp)
}

func (trs *Transport) GetReverseHandler(endpoint endpoint.Endpoint, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		endpoint,
		trs.decodeReverseRequest,
		trs.encodeReverseResponse,
		options...,
	)
}

func (trs *Transport) decodeReverseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req models.ReverseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	} else {
		level.Info(trs.logger).Log("parsed request", req.Word)
		return &req, nil
	}
}

func (trs *Transport) encodeReverseResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(*models.ReverseResponse)

	if !ok {
		return errors.New("error decoding")
	}

	return json.NewEncoder(w).Encode(res)
}
