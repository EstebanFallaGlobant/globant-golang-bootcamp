package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
)

type endpoints interface {
	MakeGetUserEndpoint() endpoint.Endpoint
}

func NewHTTPServer(logger kitlog.Logger, endpoints endpoints) *mux.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeErrorResponse),
	}

	getUserHandler := httptransport.NewServer(
		endpoints.MakeGetUserEndpoint(),
		makeDecodeGetUserRequest(logger),
		makeEncodeGetUserResponse(logger),
		options...,
	)

	r := mux.NewRouter()

	subRouter := r.PathPrefix(getUserPath+"/").
		HeadersRegexp(authTknHeaderName, "[[:graph:]]+?").
		Subrouter()

	subRouter.Methods(http.MethodGet).
		Path("/{id:[0-9]+}").
		Handler(getUserHandler)

	return r
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	code := getErrorCode(err)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	response := errorResponse{
		Message: err.Error(),
		Code:    code,
	}

	json.NewEncoder(w).Encode(response)
}

func makeDecodeGetUserRequest(logger kitlog.Logger) httptransport.DecodeRequestFunc {
	return func(c context.Context, r *http.Request) (request interface{}, err error) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		token := r.Header.Get(authTknHeaderName)

		if err != nil {
			return nil, err
		}

		return getUserRequest{authTokent: token, userID: id}, nil
	}
}

func makeEncodeGetUserResponse(logger kitlog.Logger) httptransport.EncodeResponseFunc {
	return func(c context.Context, rw http.ResponseWriter, response interface{}) error {
		rw.WriteHeader(http.StatusOK)
		return json.NewEncoder(rw).Encode(response)
	}
}

func getErrorCode(err error) int {
	switch err.(type) {
	case svcerr.InvalidArgumentError:
		return http.StatusBadRequest
	case svcerr.InvalidRequestError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
