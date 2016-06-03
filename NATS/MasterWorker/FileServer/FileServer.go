package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"io"
	"fmt"
	"github.com/nats-io/nats"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

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

	RunServiceDiscoverable()

	http.ListenAndServe(":3000", m)
}

func RunServiceDiscoverable() {
	nc, err := nats.Conn{os.Args[1]}
	if err != nil {
		fmt.Println("Can't connect to NATS. Service is not discoverable.")
	}
	nc.Subscribe("Discover.FileServer", func(m *nats.Msg) {
		nc.Publish(m.Reply, )
	})
}
