package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/repository"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
)

const (
	statusMsg = "service status"
)

var address string
var connString string

func main() {
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	flag.StringVar(&connString, "connstr", "POC_gRPC:P0cP4ssw0rd*@tcp(127.0.0.1:3308)/gRPC_Db", "connection string used to connect with a mysql database, must be in format: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]\nMore examples of valid format found in: https://github.com/go-sql-driver/mysql")
	flag.StringVar(&address, "addr", ":5050", "address in which the service will be listening")
	flag.Parse()

	level.Info(logger).Log(statusMsg, "starting")
	defer level.Info(logger).Log(statusMsg, "service finished")

	db, err := sql.Open("mysql", connString)

	level.Info(logger).Log("status", "database connection open")

	if err != nil {
		level.Error(logger).Log("database not connected", err)
		os.Exit(1)
	}

	sqlErrorHandler := repository.MySQLErrorHandler{Logger: logger}

	level.Info(logger).Log("status", "mySQL error handler created")

	repository := repository.NewsqlRepository(logger, db, sqlErrorHandler)

	level.Info(logger).Log("status", "repository layer created")

	svc := user.NewService(repository, logger)

	level.Info(logger).Log("status", "service layer created")

	endpoints := user.MakeEndpoints(svc, logger, nil)

	level.Info(logger).Log("status", "endpoints created")

	grpcServer := user.NewgRPCServer(endpoints, logger, user.UserServiceErrorHandler{Logger: logger})

	level.Info(logger).Log("status", "gRPC server created")

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	grpcListener, err := net.Listen("tcp", address)

	level.Info(logger).Log("status", "Listener created")

	if err != nil {
		level.Error(logger).Log("server listener creation failed", err)
		os.Exit(1)
	}

	go func() {
		server := grpc.NewServer()

		pb.RegisterUserDetailServiceServer(server, grpcServer)

		level.Info(logger).Log("server register", "success")
		level.Info(logger).Log("listening", address)

		if err := server.Serve(grpcListener); err != nil {
			level.Error(logger).Log("server error", err)
		}
	}()

	level.Error(logger).Log("exit", <-errs)
}
