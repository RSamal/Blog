package main

import (
	"github.com/nats-io/nats"
	natsp "github.com/nats-io/nats/encoders/protobuf"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	nc, err := nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	ec, err := nats.NewEncodedConn(nc, natsp.PROTOBUF_ENCODER)
	defer ec.Close()


}
