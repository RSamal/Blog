# Practical Golang: Using Protobuffers.

## Introduction

Most apps we make need a means of communication. We usually use *JSON*, or just plain text. *JSON* has got especially popular because of the rise of *Node.js*. The truth though, is, that JSON isn't really a fast format. The marshaller in Go also isn't that fast. That's why in this article we'll learn how to use [google protocol buffers][1]. They are in fact very easy to use, and are much faster than JSON.

Regarding the performance gains, here they are, according to [this benchmark][2]:

| benchmark                    | iter    | time/iter  | bytes alloc | allocs       |
|------------------------------|---------|------------|-------------|--------------|
| BenchmarkJsonMarshal-8       | 500000  | 3714 ns/op | 1232 B/op   | 10 allocs/op |
| BenchmarkJsonUnmarshal-8     | 500000  | 4125 ns/op | 416 B/op    | 7 allocs/op  |
| BenchmarkProtobufMarshal-8   | 1000000 | 1554 ns/op | 200 B/op    | 7 allocs/op  |
| BenchmarkProtobufUnmarshal-8 | 1000000 | 1055 ns/op | 192 B/op    | 10 allocs/op |

Ok, now let's set up the environment.

## Setup

First we'll need to get the *protobuffer compiler* binaries from here:
https://github.com/google/protobuf/releases/tag/v3.0.0-beta-3
Unpack them somewhere in your **PATH**.

The next step is to get the golang plugin. Make sure that **GOPATH**/bin is in your **PATH**.
```
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Writing .proto files

Now it's time to define our structure we'll use. I'll create mine in my project root. I'll call it *clientStructure.proto*.

First we need to define the version of *protobuffers* we will use. Here we will use the newest - **proto3**. We'll also define the package of the file. This will also be our go package name of the generated file.

```proto
syntax = "proto3";
package main;
```

Ok, now we'll define our main structure in the file. The *Client* structure:

```proto
message Client {

}
```

Now it's time to define the available fields. Fields are refered to by id, so for each field we define the type, name and id like this:

```proto
type name = id;
```

We'll start with our first 4 fields of our client:

```proto
message Client {
    int32 id = 1;
    string name = 2;
    string email = 3;
    string country = 4;
}
```

We will also define an inner structure *Mail*:

```proto
    string country = 4;

    message Mail {
        string remoteEmail = 1;
        string body = 2;
    }
```

and finally define the inbox field. It's an array of mails, which we create using the *repeated* keyword:

```proto
    message Mail {
        string remoteEmail = 1;
        string body = 2;
    }

    repeated Mail inbox = 5;
}
```






[1]: https://github.com/google/protobuf
[2]: https://github.com/alecthomas/go_serialization_benchmarks




Let's give him the creative name john doe
