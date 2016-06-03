package main

import (
	"github.com/satori/go.uuid"
	"github.com/cube2222/Blog/NATS/MasterWorker"
	"os"
	"fmt"
)

var Tasks []Transport.Task

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	Tasks = make([]Transport.Task, 0, 20)
}

func initTestTasks() {
	for i := 0; i < 20; i++ {
		newTask := Transport.Task{Uuid: uuid.NewV4(), State: 0}


	}
}
