package main

import (
	"net/http"
	"os"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/grpcmicroservice/server"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
)

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}
	var options []httptransport.ServerOption

	svc := server.NewService(logger)
	endpoints := server.MakeEndpoints(svc, logger, nil)
	trs := server.NewTransport(logger)

	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/palindrome").Handler(trs.GetIsPalHandler(endpoints.GetIsPalindrome, options))
	router.Methods(http.MethodGet).Path("/reverse").Handler(trs.GetReverseHandler(endpoints.GetReverse, options))

	addr := "127.0.0.1:8080"
	level.Info(logger).Log("status", "listening", "address", addr)
	svr := http.Server{
		Addr:    addr,
		Handler: router,
	}

	level.Error(logger).Log(svr.ListenAndServe())
}
