package main

import (
	"os"
	"fmt"
	"github.com/nats-io/nats"
	natsp "github.com/nats-io/nats/encoders/protobuf"
	"github.com/cube2222/Blog/NATS/EventSubs"
	"time"
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

	for i := 0; i < 5; i++ {
		myMessage := Transport.TextMessage{Id: int32(i), Body: "Hello over standard!"}

		err := ec.Publish("Messaging.Text.Standard", &myMessage)
		if err != nil {
			fmt.Println(err)
		}
	}

	for i := 5; i < 10; i++ {
		myMessage := Transport.TextMessage{Id: int32(i), Body: "Hello, please respond!"}

		res := Transport.TextMessage{}
		err := ec.Request("Messaging.Text.Respond", &myMessage, &res, 200 * time.Millisecond)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res.Body, " with id ", res.Id)

	}

	sendChannel := make(chan *Transport.TextMessage)
	ec.BindSendChan("Messaging.Text.Channel", sendChannel)
	for i := 10; i < 15; i++ {
		myMessage := Transport.TextMessage{Id: int32(i), Body: "Hello over channel!"}

		sendChannel <- &myMessage
	}
}
