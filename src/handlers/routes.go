package handlers

//
//import (
//	"flag"
//	"github.com/go-chi/chi"
//	"github.com/go-chi/chi/middleware"
//	"github.com/go-chi/render"
//)
//
//func InitRoutes() *chi.Mux {
//	flag.Parse()
//	r := chi.NewRouter()
//	r.Use(middleware.RequestID)
//	r.Use(middleware.Logger)
//	r.Use(middleware.Recoverer)
//	r.Use(middleware.URLFormat)
//	r.Use(render.SetContentType(render.ContentTypeJSON))
//	r.Route("/api", func(r chi.Router) {
//		r.Route("/fetcher", func(r chi.Router) {
//			fetcher := FetcherController()
//			// / post
//			r.Get("/", fetcher.methods.get)
//			r.Post("/", fetcher.methods.save)
//
//			// /$id delete
//
//		})
//	})
//	return r
//}
