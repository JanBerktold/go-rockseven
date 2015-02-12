package rock7

import (
	"net/http"
	"testing"
)

func TestInterface(t *testing.T) {
	http.Handle("/recieve", NewEndpoint())
}
