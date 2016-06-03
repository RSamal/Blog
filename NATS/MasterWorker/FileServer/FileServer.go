package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"io"
	"fmt"
	"github.com/nats-io/nats"
	"github.com/cube2222/Blog/NATS/MasterWorker"
	"github.com/golang/protobuf/proto"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	m := mux.NewRouter()

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Open("/tmp/" + vars["name"])
		defer file.Close()
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
		file, err := os.Create("/tmp/" + vars["name"])
		defer file.Close()
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
	nc, err := nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println("Can't connect to NATS. Service is not discoverable.")
	}
	nc.Subscribe("Discovery.FileServer", func(m *nats.Msg) {
		serviceAddressTransport := Transport.DiscoverableServiceTransport{"http://localhost:3000"}
		data, err := proto.Marshal(&serviceAddressTransport)
		if err == nil {
			nc.Publish(m.Reply, data)
		}
	})
}
