package rock7

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func fakeRequest(handler http.Handler, method, msg string) (code int, returnBody string) {
	req := constructRequest(method, msg)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}

func constructRequest(method, msg string) *http.Request {
	param := url.Values{}
	param.Add("imei", "123456789")
	param.Add("momsn", "12345")
	param.Add("transmit_time", time.Now().UTC().Format("06-02-01 15:04:05"))
	param.Add("iridium_latitude", "54.123")
	param.Add("iridium_longitude", "23.987")
	param.Add("iridium_cep", "2")
	param.Add("data", hex.EncodeToString([]byte(msg)))
	req, _ := http.NewRequest(method, "http://localhost/recieve", bytes.NewBufferString(param.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func compareMessage(msg Message, body string) bool {
	diff := time.Now().UTC().Sub(msg.TransmitTime)
	return msg.Data == body &&
		msg.MOMSN == 12345 &&
		msg.IMEI == "123456789" &&
		msg.IridiumCep == 2 &&
		*msg.IridumPos == *geo.NewPoint(54.123, 23.987) &&
		diff.Seconds() < 1
}

func TestInterface(t *testing.T) {
	http.Handle("/recieve", NewEndpoint())
}

func TestSimpleMessage(t *testing.T) {
	endpoint := NewEndpoint()

	go func() {
		for i := 0; i < 5; i++ {
			code, _ := fakeRequest(endpoint, "POST", fmt.Sprintf("Request %v", i))
			if code != 200 {
				t.Fatalf("Recieved non-OK status %v", code)
			}
		}
	}()

	for i := 0; i < 5; i++ {
		message := <-endpoint.GetChannel()
		if !compareMessage(message, fmt.Sprintf("Request %v", i)) {
			t.Fatalf("Failed")
		}
	}

}

func TestWrongMethod(t *testing.T) {

	endpoint := NewEndpoint()

	go func() {
		<-endpoint.GetChannel()
		t.Fatal("Recieved a message, even though none should have been sent.")
	}()

	if code, _ := fakeRequest(endpoint, "GET", "RequestData"); code == 200 {
		t.Fatalf("Non-OK code expected, got %v", code)
	}

}
