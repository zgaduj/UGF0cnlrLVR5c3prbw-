package middlewares

import (
	"log"
	"net/http"
)

func BodyMaxSize(maxSize int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print("ContentLength == ", r.ContentLength)
			log.Print("maxSize == ", maxSize)
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			err := r.ParseForm()
			if err != nil || r.ContentLength > maxSize {
				// redirect or set error status code.
				log.Print("oversize")
				w.WriteHeader(413)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
