package rock7

import (
	"encoding/hex"
	"github.com/kellydunn/golang-geo"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	IMEI         string
	MOMSN        int
	TransmitTime time.Time
	IridumPos    *geo.Point
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

func (end *Endpoint) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data := req.FormValue("data")
	by, err := hex.DecodeString(data)

	// TODO: Read time
	if err == nil {
		end.channel <- Message{
			req.FormValue("imei"),
			readInt(req.FormValue("momsn")),
			time.Now(),
			geo.NewPoint(
				readFloat64(req.FormValue("iridium_latitude")),
				readFloat64(req.FormValue("iridium_longitude")),
			),
			readInt(req.FormValue("iridium_cep")),
			data,
			string(by),
		}
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
	}

}
