# Practical Golang: Getting started with NATS and related patterns

## Introduction

Microservices... the never disappearing buzzword of our times. They promise a lot, but can be slow or complicated if not implemented correctly. One of the main challenges when developing and using a microservice-based architecture is getting the *communication* right. Many will ask, why not **REST**? As I did at some point. Many will actually use it. But the truth is that it leads to tighter coupling, and is *synchronous*. Microservice architectures are meant to be asynchronous. Also, REST is blocking, which also isn't good on many occasions.

What are we meant to use for communication? Usually we use:
  - RPC - Remote Procedure Call
  - Message BUS/Broker

In this article I'll write about one specific Message BUS called **NATS** and using it in Go.

There are also other message BUS'ses/Brokers. Some popular ones are **Kafka** and **RabbitMQ**.

Why NATS? It's simple, and astonishingly fast.

## Setting up NATS

To use NATS you can do one of the following things:
  1. Use the [NATS Docker image][1]
  2. Get the [binaries][2]
  3. Use the public NATS server *nats://demo.nats.io:4222*
  4. Build from [source][3]

Also, remember to
```
go get https://github.com/nats-io/nats
```
the official Go library.

## Getting started

In this article we'll be using protobuffs a lot. So if you want to know more about them, check out my [previous article about protobuffs][4].

First, let's write one of the key usages of microservices. A fronted, that lists information from other micrservices, but doesn't care if one of them is down. It will respond to the user anyways. This makes microservices swappable live, one at a time.

In each of our services we'll need to connect to NATS:

```go
package main

import (
	"github.com/nats-io/nats"
	"fmt"
)

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
}
```

Now, let's write the first provider service. It will receive a *User Id*, and answer with a *user name* For that we'll need a transport structure to send its data over *NATS*. I wrote this short proto file for that:

```proto
syntax = "proto3";
package Transport;

message User {
    string id = 1;
    string name = 2;
}
```

Now we will create the map containing our user names:

```go
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
}
```

and finally the part that's most interesting to us. Subscribing to the topic:

```go
users["4"] = "Kate"

nc.QueueSubscribe("UserNameById", "userNameByIdProviders", replyWithUserId)
```

Notice that it's a **QueueSubscribe**. Which means that if we start 10 instances of this service in the *userNameByIdProviders* group , only one will get each message sent over *UserNameById*. Another thing to note is that this function call is asynchronous, so we need to block somehow. This will provide an endless block:

```go
nc.QueueSubscribe("UserNameById", "userNameByIdProviders", replyWithUserId)
wg := sync.WaitGroup{}

wg.Add(1)
wg.Wait()
}
```

Ok, now to the *replyWithUserId* function:

```go
func replyWithUserId(m *nats.Msg) {
	}
```

Notice that it takes one argument, a pointer to the message.

We'll unmarshal the data:

```go
func replyWithUserId(m *nats.Msg) {

	myUser := Transport.User{}
	err := proto.Unmarshal(m.Data, &myUser)
	if err != nil {
		fmt.Println(err)
		return
	}
```

get the name and marshal back:

```go
myUser.Name = users[myUser.Id]
data, err := proto.Marshal(&myUser)
if err != nil {
  fmt.Println(err)
  return
}
```

And, as this shall be a request we're handling, we respond to the *Reply topic*, a topic created by the caller exactly for this purpose:

```go
if err != nil {
  fmt.Println(err)
  return
}
fmt.Println("Replying to ", m.Reply)
nc.Publish(m.Reply, data)

}
```

Ok, now let's get to the second service. Our time provider service, first the same basic structure:

```go
package main

import (
	"github.com/nats-io/nats"
	"fmt"
	"github.com/cube2222/Blog/NATS/FrontendBackend"
	"github.com/golang/protobuf/proto"
	"os"
	"sync"
	"time"
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
```

This time we're not getting any data from the caller, so we just marshal our time into this proto structure:

```proto
syntax = "proto3";
package Transport;

message Time {
    string time = 1;
}
```

and send it back:

```go
func replyWithTime(m *nats.Msg) {
	curTime := Transport.Time{time.Now().Format(time.RFC3339)}

	data, err := proto.Marshal(&curTime)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Replying to ", m.Reply)
	nc.Publish(m.Reply, data)

}
```

We can now get to our frontend, which will use both those services. First the standard basic structure:
```go
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/cube2222/Blog/NATS/FrontendBackend"
	"github.com/golang/protobuf/proto"
	"fmt"
	"github.com/nats-io/nats"
	"time"
	"os"
	"sync"
)

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

	m := mux.NewRouter()
	m.HandleFunc("/{id}", handleUserWithTime)

	http.ListenAndServe(":3000", m)
}
```

That's a pretty standard web server, now to the interesting bits, the *handleUserWithTime* function, which will respond with the user name and time:

```go
func handleUserWithTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	myUser := Transport.User{Id: vars["id"]}
	curTime := Transport.Time{}
	wg := sync.WaitGroup{}
	wg.Add(2)
}
```
We've parsed the request arguments and started a *WaitGroup* with the value two, as we will do one *asynchronous* request for each of our services. First we'll marshal the user struct:

```go
go func() {
  data, err := proto.Marshal(&myUser)
  if err != nil || len(myUser.Id) == 0 {
    fmt.Println(err)
    w.WriteHeader(500)
    fmt.Println("Problem with parsing the user Id.")
    return
  }
```

and, then we make a request. Sending the user data, and waiting at most 100 ms for the response:

```go
fmt.Println("Problem with parsing the user Id.")
return
}

msg, err := nc.Request("UserNameById", data, 100 * time.Millisecond)
```

now we can check if any error happend, or the response is empty and finish this thread:

```go
msg, err := nc.Request("UserNameById", data, 100 * time.Millisecond)
if err == nil && msg != nil {
  myUserWithName := Transport.User{}
  err := proto.Unmarshal(msg.Data, &myUserWithName)
  if err == nil {
    myUser = myUserWithName
  }
}
wg.Done()
}()
```

Next we'll do the request to the *Time Tellers*.
We again make a request, but its body is nil, as we don't need to pass any data:

```go
go func() {
  msg, err := nc.Request("TimeTeller", nil, 100*time.Millisecond)
  if err == nil && msg != nil {
    receivedTime := Transport.Time{}
    err := proto.Unmarshal(msg.Data, &receivedTime)
    if err == nil {
      curTime = receivedTime
    }
  }
  wg.Done()
}()
```

After both requests finished (or failed) we can just respond to the user:

```go
wg.Wait()

fmt.Fprintln(w, "Hello ", myUser.Name, " with id ", myUser.Id, ", the time is ", curTime.Time, ".")
}
```

Now if you actually test it, you'll notice that if one of the provider services isn't active, the frontend will respond anyways, putting a zero'ed value in place of the non-available resource. You could also make a template that shows an error in that place.

Ok, that was already an interesting architecture. Now we can implement...

## The Master-Slave pattern

This is such a popular pattern, especially in Go, that we really should know how to implement it. The workers will do simple operations on a text file (count the usage amounts of each word in a comma-separated list).

Now you could think that the *Master*, should send the files to the *Workers* over NATS. Wrong. This would lead to a huge slowdown of NATS (at least for bigger files). That's why the *Master* will send the files to a file server over a REST API, and the *Workers* will get it from there. We'll also learn how to do *service discovery* over NATS.

First, the *File Server*. I won't really go through the file handling part, as it's a simple get/post API.I will however, go over the *service discovery* part.

```go
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"io"
	"fmt"
	"github.com/nats-io/nats"
	"github.com/cube2222/Blog/NATS/MasterWorker"
	"github.com/golang/protobuf/proto"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments. Need NATS server address.")
		return
	}

	m := mux.NewRouter()

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Open("/tmp/" + vars["name"])
		defer file.Close()
		if err != nil {
			w.WriteHeader(404)
		}
		if file != nil {
			_, err := io.Copy(w, file)
			if err != nil {
				w.WriteHeader(500)
			}
		}
	}).Methods("GET")

	m.HandleFunc("/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file, err := os.Create("/tmp/" + vars["name"])
		defer file.Close()
		if err != nil {
			w.WriteHeader(500)
		}
		if file != nil {
			_, err := io.Copy(file, r.Body)
			if err != nil {
				w.WriteHeader(500)
			}
		}
	}).Methods("POST")

	RunServiceDiscoverable()

	http.ListenAndServe(":3000", m)
}
```

Now, what does the *RunServiceDiscoverable* function do? It connects to the NATS server and responds with its own http address to incoming requests.

```go
func RunServiceDiscoverable() {
	nc, err := nats.Connect(os.Args[1])
	if err != nil {
		fmt.Println("Can't connect to NATS. Service is not discoverable.")
	}
	nc.Subscribe("Discovery.FileServer", func(m *nats.Msg) {
		serviceAddressTransport := Transport.DiscoverableServiceTransport{"http://localhost:3000"}
		data, err := proto.Marshal(&serviceAddressTransport)
		if err == nil {
			nc.Publish(m.Reply, data)
		}
	})
}
```

The proto file looks like this:

```proto
syntax = "proto3";
package Transport;

message DiscoverableServiceTransport {
    string Address = 1;
}
```



[1]: https://hub.docker.com/_/nats/
[2]: https://github.com/nats-io/gnatsd/releases/
[3]: https://github.com/nats-io/gnatsd
[4]: https://jacobmartins.com/2016/05/24/practical-golang-using-protobuffs/
