package rock7

import (
	"github.com/kellydunn/golang-geo"
	"net/http"
	"time"
)

type Message struct {
	IMEI         string
	MOMSN        int
	TransmitTime time.Time
	IridumPos    geo.Point
	IridiumCep   int
	HexData      string
	Data         string
}

type Endpoint struct {
	channel chan Message
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		make(chan Message),
	}
}

// Returns a channel contained all recieved messages.
// Getter method is required due to access locking.
func (end *Endpoint) GetChannel() <-chan Message {
	return end.channel
}

func (end *Endpoint) ServeHTTP(http.ResponseWriter, *http.Request) {

}
