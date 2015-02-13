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

func fakeRequest(handler http.Handler, req *http.Request) (code int, returnBody string) {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	return w.Code, w.Body.String()
}

func constructRequest(method, msg string, params ...url.Values) *http.Request {
	var param url.Values
	if len(params) == 0 {
		param = url.Values{}
		param.Add("imei", "123456789")
		param.Add("momsn", "12345")
		param.Add("transmit_time", time.Now().UTC().Format("06-02-01 15:04:05"))
		param.Add("iridium_latitude", "54.123")
		param.Add("iridium_longitude", "23.987")
		param.Add("iridium_cep", "2")
		param.Add("data", hex.EncodeToString([]byte(msg)))
	} else {
		param = params[0]
	}
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
			code, _ := fakeRequest(endpoint, constructRequest("POST", fmt.Sprintf("Request %v", i)))
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

	if code, _ := fakeRequest(endpoint, constructRequest("GET", "RequestData")); code == 200 {
		t.Fatalf("Non-OK code expected, got %v", code)
	}
}

func TestWrongMomsn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Call did not fail\n")
		}
	}()

	endpoint := NewEndpoint()

	param := url.Values{}
	param.Add("imei", "123456789")
	param.Add("momsn", "abc")
	req := constructRequest("POST", "RequestData", param)

	fakeRequest(endpoint, req)
	t.Fatal("This point should not have been reached.")
}

func TestWrongTime(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Call did not fail\n")
		}
	}()

	endpoint := NewEndpoint()

	param := url.Values{}
	param.Add("imei", "123456789")
	param.Add("momsn", "12")
	param.Add("transmit_time", "06/02/01")
	req := constructRequest("POST", "RequestData", param)

	fakeRequest(endpoint, req)
	t.Fatal("This point should not have been reached.")
}

func TestWrongGeoPos(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Call did not fail\n")
		}
	}()

	endpoint := NewEndpoint()

	param := url.Values{}
	param.Add("imei", "123456789")
	param.Add("momsn", "12")
	param.Add("transmit_time", "06-02-01 15:04:05")
	param.Add("iridium_latitude", "ab25ad8g.12")
	req := constructRequest("POST", "RequestData", param)

	fakeRequest(endpoint, req)
	t.Fatal("This point should not have been reached.")
}

func TestWrongHex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Call did not fail\n")
		}
	}()

	endpoint := NewEndpoint()

	param := url.Values{}
	param.Add("imei", "123456789")
	param.Add("momsn", "12345")
	param.Add("transmit_time", time.Now().UTC().Format("06-02-01 15:04:05"))
	param.Add("iridium_latitude", "54.123")
	param.Add("iridium_longitude", "23.987")
	param.Add("iridium_cep", "2")
	param.Add("data", "afadfpoi3a5oudf")
	req := constructRequest("POST", "RequestData", param)

	fakeRequest(endpoint, req)
	t.Fatal("This point should not have been reached.")
}
