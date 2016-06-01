package main

import (
	"github.com/nats-io/nats"
	"os"
	"fmt"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return 1
	}
	users := make(map[string]string)
	users["1"] = "Bob"
	users["2"] = "John"
	users["3"] = "Dan"
	users["4"] = "Kate"
	var nc *nats.Conn
	var err error
	for nc, err = nats.Connect(os.Args[1]); err != nil; {
		fmt.Println(err)
		fmt.Println("Trying again in 2 seconds...")
		time.Sleep(time.Second * 2)
	}

	nc.Subscribe("UserNameById", )


}
