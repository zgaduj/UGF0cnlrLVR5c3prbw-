package main

import (
	"app/config"
	"app/src/database"
	"app/src/handlers"
	"app/src/middlewares"
	"app/src/worker"
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

type Env struct {
	db     *gorm.DB
	config config.Config
}

func main() {
	worker.RunWorker()

	_config := config.GetConfig()

	env := &Env{
		config: *_config,
	}

	dbStore, error := database.NewConnection(env.config.DB.FilePath)
	if error != nil {
		log.Panic(error)
	} else {

		defer dbStore.Close()

		env.db = dbStore.DB

		flag.Parse()
		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)
		r.Use(middlewares.BodyMaxSize(_config.APP.MaxBodySize))
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
			})
		})

		_port := strconv.Itoa(env.config.APP.ListeningPort)
		log.Print("Start listening on port: ", _port)
		http.ListenAndServe(":"+_port, r)
	}
	log.Fatal("Error connect with DB")
}
