package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"github.com/cube2222/Blog/NATS/FrontendBackend"
	"github.com/golang/protobuf/proto"
	"os"
	"sync"
)

// We use globals because it's a small application demonstrating NATS.

var users map[string]string
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

	users = make(map[string]string)
	users["1"] = "Bob"
	users["2"] = "John"
	users["3"] = "Dan"
	users["4"] = "Kate"

	nc.QueueSubscribe("UserNameById", "userNameByIdProviders", replyWithUserId)
	wg := sync.WaitGroup{}

	wg.Add(1)
	wg.Wait()
}

func replyWithUserId(m *nats.Msg) {

	myUser := Transport.User{}
	err := proto.Unmarshal(m.Data, &myUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	myUser.Name = users[myUser.Id]
	data, err := proto.Marshal(&myUser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)

}
