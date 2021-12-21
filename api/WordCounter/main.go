package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/WordCounter/wordcounterapi"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/WordCounter/wordcounterapi/structs"
)

func main() {
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	Args := os.Args

	var addr string

	if len(Args) > 1 {
		addr = Args[1]
	} else {
		addr = "127.0.0.1:8080"
	}

	var api wordcounterapi.App

	api.Initialize(new(structs.WordCounter))

	svr := api.Run(addr)

	go func() {
		if err := svr.ListenAndServe(); err != nil {
			log.Fatal(err)
		}

	}()

	fmt.Printf("server running on: %s\n", svr.Addr)
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	svr.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}
