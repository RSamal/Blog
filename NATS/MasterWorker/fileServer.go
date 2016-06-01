package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"io"
)

func main() {
	m := mux.NewRouter()

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Open(vars["name"])
		if err != nil {
			w.WriteHeader(404)
		}
		if file != nil {
			_, err := io.Copy(w, file)
			if err != nil {
				w.WriteHeader(500)
			}
		}
	}).Methods("GET")

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Create(vars["name"])
		if err != nil {
			w.WriteHeader(500)
		}
		if file != nil {
			_, err := io.Copy(file, r.Body)
			if err != nil {
				w.WriteHeader(500)
			}
		}
	}).Methods("POST")
	http.ListenAndServe(":3000", m)
}
