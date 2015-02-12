# go-rockseven [![Build Status](https://travis-ci.org/JanBerktold/go-rockseven.svg)](https://travis-ci.org/JanBerktold/go-rockseven)

A simple package for interacting with your http://www.rock7mobile.com/ devices using a web interface.

## Sending

Sending a message to an endpoint, providing the IMEI number:

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

Alternatively, you can set a default IMEI number.

	client.SetDefaultIMEI("1234689")
	code, err := client.SendStringToDefault("Hello, world!")


Sending a byte slice can be done using the corresponding methods:

	client.SetDefaultIMEI("1234689")
	code, err := client.SendToDefault([]byte{79, 75})

or

	code, err := client.Send("1234689", []byte{79, 75})

## Receiving (Draft)


	import (
		"net/http"
		"github.com/janberktold/go-rockseven"
		"fmt"
	)

	func printMessages(end *rock7.Endpoint) {
		for {
			msg := <-end.GetChannel()
			fmt.Printf("Recieved message %q\n", msg)
		}
	}

	func main() {
		endpoint := rock7.NewEndpoint()
		go printMessages(endpoint)
		http.Handle("recieve", endpoint)
		http.ListenAndServe(":80", nil)
	}
