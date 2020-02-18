package main

import (
	"app/app/database"
	"app/app/handlers"
	"app/app/middlewares"
	"app/config"
	"context"
	"flag"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

type APIConf struct {
	API *config.APIConfig
	DB  *config.DBConfig
}

type Env struct {
	db     *gorm.DB
	config *APIConf
}

func main() {
	log.Print("[API] START")

	apiConfig := &config.APIConfig{}
	if err := env.Parse(apiConfig); err != nil {
		log.Fatal("[ENV-API] Error parse", err)
	}

	dbConfig := &config.DBConfig{}
	if err := env.Parse(dbConfig); err != nil {
		log.Fatal("[ENV-DB] Error parse", err)
	}

	appCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	env := &Env{ // @todo: move to single config?
		config: &APIConf{
			API: apiConfig,
			DB:  dbConfig,
		},
	}

	dbStore, error := database.NewConnection(env.config.DB)

	r := chi.NewRouter()

	_port := strconv.Itoa(env.config.API.ListeningPort)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	if error != nil {
		log.Fatal(error)
	} else {

		defer dbStore.Close()

		env.db = dbStore.DB

		flag.Parse()
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)
		r.Use(middlewares.BodyMaxSize(env.config.API.MaxBodySize))
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Route("/api", func(r chi.Router) {
			r.Route("/fetcher", func(r chi.Router) {
				r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
					handlers.FetcherGet(env.db, writer, request)
				})
				r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
					handlers.FetcherSave(env.db, writer, request)
				})
				r.Delete("/{id:[0-9]+}", func(writer http.ResponseWriter, request *http.Request) {
					handlers.FetcherDelete(env.db, writer, request)
				})
				r.Get("/{id:[0-9]+}/history", func(writer http.ResponseWriter, request *http.Request) {
					handlers.FetcherHistory(env.db, writer, request)
				})
			})
		})

		log.Print("Start listening on port: ", _port)

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}
	log.Fatal("Error connect with DB")

	go func() {
		<-signals
		if err := srv.Shutdown(context.TODO()); err != nil {
			log.Fatal(err) // failure/timeout shutting down the server gracefully
		}
		cancelCtx()
	}()

	<-appCtx.Done()
}
