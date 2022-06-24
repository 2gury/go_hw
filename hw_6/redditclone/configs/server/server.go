package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	server *http.Server
}

func ConfigureStatic(mx *mux.Router) {
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	mx.PathPrefix("/static/").Handler(staticHandler)

	mx.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/")
	})

	mx.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/")
	})
}

func NewHTTPServer(handler http.Handler, port string) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
	}
}

func (s *HTTPServer) Launch() error {
	log.Println("Starting server at " + s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
