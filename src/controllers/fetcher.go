package controllers

import (
	"github.com/go-chi/chi"
	"net/http"
)

type FetcherMethods struct {
	Get        func(w http.ResponseWriter, r *http.Request)
	Save       func(w http.ResponseWriter, r *http.Request)
	Delete     func(w http.ResponseWriter, r *http.Request)
	GetHistory func(w http.ResponseWriter, r *http.Request)
}

type FetcherControllerStruct struct {
	//app     *app.App
	Routes  *chi.Mux
	Methods *FetcherMethods
}

func FetcherController() *FetcherControllerStruct {
	return &FetcherControllerStruct{
		//app: app,
		Methods: &FetcherMethods{
			Get:        fetcherGet,
			Save:       fetcherSave,
			Delete:     fetcherDelete,
			GetHistory: fetcherHistory,
		},
	}
}

func fetcherSave(w http.ResponseWriter, r *http.Request) {

}

func fetcherDelete(w http.ResponseWriter, r *http.Request) {

}

func fetcherGet(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("test asd"))
}

func fetcherHistory(w http.ResponseWriter, r *http.Request) {

}
