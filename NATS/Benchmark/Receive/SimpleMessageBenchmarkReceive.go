package main

import (
	"github.com/nats-io/nats"
	"os"
	"fmt"
	"time"
	"sync/atomic"
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
	i := int32(0)

	endChan := make(chan int)
	sub, _ := nc.SubscribeSync("Benchmark")
	sub.NextMsg(5 * time.Second)

	startTime := time.Now()

	nc.Subscribe("Benchmark", func(m *nats.Msg) {
		atomic.AddInt32(&i,1)
		if i > 1000000 {
			endChan <- 1
		}
	})

	<- endChan
	dur := time.Since(startTime).Nanoseconds()/9000000
	fmt.Println(dur)
	fmt.Println(i)
}