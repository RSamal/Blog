package main

import (
	"github.com/nats-io/nats"
	natsp "github.com/nats-io/nats/encoders/protobuf"
	"os"
	"fmt"
	"github.com/cube2222/Blog/NATS/EventSubs"
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

	ec.Subscribe("Messaging.Text.Standard", func(m *Transport.TextMessage) {
		fmt.Println("Got standard message: \"", m.Body, "\" with the Id ", m.Id, ".")
	})

	ec.Subscribe("Messaging.Text.Respond", func(subject, reply string, m *Transport.TextMessage) {
		fmt.Println("Got ask for response message: \"", m.Body, "\" with the Id ", m.Id, ".")

		newMessage := Transport.TextMessage{Id: m.Id, Body: "Responding!"}
		ec.Publish(reply, &newMessage)
	})

	receiveChannel := make(chan *Transport.TextMessage)
	ec.BindRecvChan("Messaging.Text.Channel", receiveChannel)

	for m := range receiveChannel {
		fmt.Println("Got channel'ed message: \"", m.Body, "\" with the Id ", m.Id, ".")
	}
}