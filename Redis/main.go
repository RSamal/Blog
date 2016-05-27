package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"github.com/golang/protobuf/proto"
)

func main() {
	c, err := redis.Dial("tcp", "docker:6379")
	if err != nil {
		println(err)
		defer c.Close()
	}


	m, err := c.Do("SET", "myLovelyKey", "Its lovely value2")
	fmt.Println(m)
	m, err = c.Do("SET", "myLovelyInt", 4)
	fmt.Println(m)
	val, err := redis.Int(c.Do("GET", "myLovelyInt"))
	fmt.Println(val + 1)

	mySession := Session{UserName: "cube", LoggedIn: true}
	data, _ := proto.Marshal(&mySession)

	m, _ = c.Do("SET", "123123123", data)
	fmt.Println("proto message: ",m)

	dataReceived, _ := redis.Bytes(c.Do("GET", "123123123"))

	newSession := Session{}

	proto.Unmarshal(dataReceived, &newSession)

	fmt.Println(newSession.UserName, ":", newSession.LoggedIn)
}
