package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/cube2222/Blog/NATS/FrontendBackend"
	"github.com/golang/protobuf/proto"
	"fmt"
	"github.com/nats-io/nats"
	"time"
	"os"
)

var nc *nats.Conn

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}
	var err error

	nc, err = nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	m := mux.NewRouter()
	//m.HandleFunc("/", handleUser)
	m.HandleFunc("/{id}", handleUser)


	http.ListenAndServe(":3000", m)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	myUser := userTransport.User{Id: vars["id"]}
	data, err := proto.Marshal(&myUser)
	if err != nil || len(myUser.Id) == 0 {
		fmt.Println(err)
		w.WriteHeader(500)
		fmt.Println("Problem with parsing the user Id.")
		return
	}

	msg, err := nc.Request("UserNameById", data, 100 * time.Millisecond)
	if err == nil && msg != nil {
		myUserWithName := userTransport.User{}
		err := proto.Unmarshal(msg.Data, &myUserWithName)
		if err == nil {
			myUser = myUserWithName
		}
	}

	fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id, ", the time is ")
}