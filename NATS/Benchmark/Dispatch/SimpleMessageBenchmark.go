package main

import (
	"github.com/nats-io/nats"
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

	someData := []byte("myData")
	for i := 0; i < 10000000; i++ {
		nc.Publish("Benchmark", someData)
	}
}
