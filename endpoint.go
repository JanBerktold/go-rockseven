package rock7

import (
	"encoding/hex"
	"github.com/kellydunn/golang-geo"
	"net/http"
	"strconv"
	"time"
)

// Message represents all parameters of a recieved message.
type Message struct {
	IMEI         string
	MOMSN        int
	TransmitTime time.Time
	IridumPos    *geo.Point
	IridiumCep   int
	HexData      string
	Data         string
}

// Endpoint is being used as the representation of the server which receives incoming messages
// from deployed devices. Should always be created using NewEndpoint.
type Endpoint struct {
	channel chan Message
}

// NewEndpoint creates a new endpoint for recieving messages. Should be added to an HttpMux
// using http.Handle("/your/path", endpoint).
func NewEndpoint() *Endpoint {
	return &Endpoint{
		make(chan Message),
	}
}

// Read returns the next message recieved. Convenience wrapper around the channel.
func (end *Endpoint) Read() Message {
	return <-end.channel
}

// ReadWithTimeout returns the next message recieved or nil, if the time limit has been hit.
// Convenience wrapper around the channel.
func (end *Endpoint) ReadWithTimeout(dur time.Duration) (Message, error) {
	select {
	case msg := <-end.channel:
		return msg, nil
	case <-time.After(dur):
		return Message{}, ErrTimeOut
	}
}

// GetChannel returns a channel containing all recieved messages.
// Getter method is required due to access locking.
func (end *Endpoint) GetChannel() <-chan Message {
	return end.channel
}

func readInt(n string) int {
	num, mErr := strconv.Atoi(n)

	if mErr != nil {
		panic("Invalid request, panic'ing to cancel request")
	}
	return num
}

func readFloat64(n string) float64 {
	num, mErr := strconv.ParseFloat(n, 64)

	if mErr != nil {
		panic("Invalid request, panic'ing to cancel request")
	}
	return num
}

func readHex(n string) string {
	by, err := hex.DecodeString(n)

	if err != nil {
		panic("Invalid request, panic'ing to cancel request")
	}
	return string(by)
}

func readTime(n string) time.Time {
	tim, err := time.Parse("06-02-01 15:04:05", n)

	if err != nil {
		panic("Invalid request, panic'ing to cancel request")
	}
	return tim
}

// Fulfills the requirements for the http.Handler interface. This method should never be called by your code, as it is triggered by the net/http implementation of a HTTP server.
// Allows you to set a created Endpoint as the handler for a URL using http.Handle("recieve", rock7.NewEndpoint()).
func (end *Endpoint) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	end.channel <- Message{
		req.FormValue("imei"),
		readInt(req.FormValue("momsn")),
		readTime(req.FormValue("transmit_time")),
		geo.NewPoint(
			readFloat64(req.FormValue("iridium_latitude")),
			readFloat64(req.FormValue("iridium_longitude")),
		),
		readInt(req.FormValue("iridium_cep")),
		req.FormValue("data"),
		readHex(req.FormValue("data")),
	}
	writer.WriteHeader(http.StatusOK)
}
