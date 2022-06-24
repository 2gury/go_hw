package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Add(x, y int) (res int) {
	return x + y
}

func main() {
	mx := mux.NewRouter()
	mx.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hi"))
		if err != nil {
			return
		}
	}).Methods(http.MethodGet)

	log.Println("Start server at :8080 port")
	err := http.ListenAndServe(":8080", mx)
	if err != nil {
		log.Fatalln(err)
	}
}
