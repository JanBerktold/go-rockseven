# go-rockseven
[![Build Status](https://travis-ci.org/JanBerktold/go-rockseven.svg)](https://travis-ci.org/JanBerktold/go-rockseven) [![Coverage Status](https://coveralls.io/repos/JanBerktold/go-rockseven/badge.svg?branch=master)](https://coveralls.io/r/JanBerktold/go-rockseven?branch=master) [![GoDoc](http://godoc.org/github.com/janberktold/go-rockseven?status.svg)](http://godoc.org/github.com/janberktold/go-rockseven)


A simple package for interacting with your http://www.rock7mobile.com/ devices using HTTP requests.

## Installation

Package is designed to be compatible with every major Go release since 1.2. It can be installed using Go's native toolchain.

		go get github.com/janberktold/go-rockseven

## TODO

- Perform final tests with RockBLOCK device
- Write further GoDoc documentation

## Sending

Sending a message to an endpoint, providing the IMEI number:

```go
import (
	"fmt"
	"github.com/janberktold/go-rockseven"
)

func main() {
	client := rock7.NewClient("user", "pass")

	if code, err := client.SendString("1234689", "Hello, world!"); err == nil {
		fmt.Printf("Sent message, assigned messageId: %v\n", code)
	} else {
		fmt.Printf("Failed sending message %q\n", err.Error())
	}
}
```

Alternatively, you can set a default IMEI number.

```go
client.SetDefaultIMEI("1234689")
code, err := client.SendStringToDefault("Hello, world!")
```

Sending a byte slice can be done using the corresponding methods:

```go
client.SetDefaultIMEI("1234689")
code, err := client.SendToDefault([]byte{79, 75})
```

or

```go
code, err := client.Send("1234689", []byte{79, 75})
```

## Receiving

The endpoint is designed to fit nicely into golang's net/http package and can therefore be used as part of a standard HTTP server. The example below spawns a endpoint which listens on /recieve and prints all incoming messages to the stdout.

```go
import (
	"net/http"
	"github.com/janberktold/go-rockseven"
	"fmt"
)

func printMessages(end *rock7.Endpoint) {
	for {
		msg := <-end.GetChannel()
		fmt.Printf("Recieved message %q\n", msg.Data)
	}
}

func main() {
	endpoint := rock7.NewEndpoint()
	go printMessages(endpoint)
	http.Handle("recieve", endpoint)
	http.ListenAndServe(":80", nil)
}
```

In general, there are two different approaches to access incoming messages: On the one hand, you can obtain a recieve-only channel via client.GetChannel() and implement your own reading mechanism. Example:

```go
for {
	msg := <-endpoint.GetChannel()
	fmt.Printf("Recieved message %q\n", msg.Data)
}
```

Alternatively, you can use the provided convenience methods Read and ReadWithTimeout(time.Duration) to keep track of incoming messages. Both methods are blocking, as set in the Golang coding guidelines. Simple read:


```go
for {
	msg := endpoint.Read()
	fmt.Printf("Recieved message %q\n", msg.Data)
}
```

Simple read with timeout:

```go
for {
	if msg, err := endpoint.ReadWithTimeout(2 * time.Second); err == nil {
		fmt.Printf("Recieved message %q\n", msg.Data)
	} else {
		fmt.Println(err.Error())
	}
}
```

A timeout can also be implemented using the raw channel:

```go
for {
	select {
	case msg := <-endpoint.GetChannel():
		fmt.Printf("Recieved message %q\n", msg.Data)
	case <-time.After(2 * time.Second):
		fmt.Println("Hit time limit")
	}
}
```
