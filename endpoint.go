package rock7

import (
	"net/http"
)

type Message struct {
}

type Endpoint struct {
	channel chan Message
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

func (end *Endpoint) ServeHTTP(http.ResponseWriter, *http.Request) {

}
