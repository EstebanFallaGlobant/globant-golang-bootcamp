package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/repository"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
)

const (
	statusMsg = "service status"
)

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	connString := os.Getenv("udsConnStr")
	if util.IsEmptyString(connString) {
		level.Error(logger).Log(statusMsg, "connections string not found")
		os.Exit(1)
	}

	address := os.Getenv("udsListenAddr")
	if util.IsEmptyString(address) {
		address = ":5050"
	}

	level.Info(logger).Log(statusMsg, "starting")
	defer level.Info(logger).Log(statusMsg, "service finished")

	db, err := sql.Open("mysql", connString)

	if err != nil {
		level.Error(logger).Log("database not connected", err)
		os.Exit(1)
	}

	sqlErrorHandler := repository.MySQLErrorHandler{Logger: logger}
	repository := repository.NewsqlRepository(logger, db, sqlErrorHandler)
	svc := user.NewService(repository, logger)
	endpoints := user.MakeEndpoints(svc, logger, nil)
	grpcServer := user.NewgRPCServer(endpoints, logger, user.UserServiceErrorHandler{Logger: logger})
	level.Info(logger).Log(statusMsg, "server created")

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", address)

	if err != nil {
		level.Error(logger).Log("server listener creation failed", err)
		os.Exit(1)
	}
	level.Info(logger).Log(statusMsg, "Listener created")

	go func() {
		baseServer := grpc.NewServer()

		pb.RegisterUserDetailServiceServer(baseServer, grpcServer)

		level.Info(logger).Log(statusMsg, "server registered")
		level.Info(logger).Log("listening at", address)

		if err := baseServer.Serve(grpcListener); err != nil {
			level.Error(logger).Log("server error", err)
		}
	}()

	level.Error(logger).Log("exit", <-errs)
}
