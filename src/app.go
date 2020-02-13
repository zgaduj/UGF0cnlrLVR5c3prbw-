package app

import (
	"app/config"
	"app/src/controllers"
	"app/src/providers"
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Config      *config.Config
	Router      *chi.Mux
	DB          *xorm.Engine
	Controllers *controllers.Controllers
}

func (a *App) SetConfig(config *config.Config) {
	a.Config = config
}

func (a *App) LoadRoutes() {
	log.Print("LoadRoutes")
	//a.Router = controllers.InitRoutes(a)
	flag.Parse()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/api", func(r chi.Router) {
		log.Print("test")
		//r.Use()
		r.Route("/fetcher", func(r chi.Router) {
			fetcher := controllers.FetcherController()
			// / post
			r.Get("/", fetcher.Methods.Get)
			r.Post("/", fetcher.Methods.Save)

			//r.Route("/{fetcherId}", func(r chi.Router) {
			//
			//})
			// /$id delete

		})
	})
	a.Router = r
}

func (a *App) SetDB() bool {
	return providers.Connect(a, a.Config.DB.FilePath)
}

func (a *App) InitApp() {
	_port := strconv.Itoa(a.Config.APP.ListeningPort)
	log.Print("Start listening on port: ", _port)
	http.ListenAndServe(":"+_port, a.Router)
}
