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

	nc.QueueSubscribe("TimeTeller", "TimeTellers", replyWithTime)
	wg := sync.WaitGroup{}

	wg.Add(1)
	wg.Wait()
}

func ReplyWithTime(m *nats.Msg) {
	err := proto.Unmarshal(m.Data, &myUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := proto.Marshal(&myUser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)

}
