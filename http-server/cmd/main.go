package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"parkpal-web-server/db"
	"parkpal-web-server/internal/business"
	"parkpal-web-server/internal/repository/storage"
	"parkpal-web-server/internal/transport/api"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9091", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "client",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s", err)
	}

	memStore := storage.NewMemStore(dbConn.GetDB())

	// business
	business := business.NewBusiness(memStore, time.Duration(2)*time.Second, l)

	// create the handlers
	apiHandler := api.NewAPI(business, l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/parking-lots", apiHandler.SearchParkingLot).Queries("name")
	getR.HandleFunc("/parking-lots", apiHandler.GetAllParkingLots)
	getR.HandleFunc("/parking-lots/{id:[0-9]+}", apiHandler.GetParkingLot)

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	//
	// // Kafka Consumer setup
	// kafkaConsumer, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "your-group-id", []string{"your-topic"}, business)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// defer kafkaConsumer.Close()
	//
	// start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// go kafkaConsumer.ConsumeMessages(ctx)

	s.Shutdown(ctx)
}
