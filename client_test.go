package rock7

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	defaultUser = "this_is_me"
	defaultPass = "not_my_pass"
)

func createTestServer(user, pass, imei string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.FormValue("username") != user || r.FormValue("password") != pass {
			fmt.Fprintf(w, "FAILED,10,invalid login credentials")
			return
		}

		if r.FormValue("imei") != imei {
			fmt.Fprintf(w, "FAILED,11,no RockBLOCK with this IMEI found on your account")
			return
		}

		fmt.Fprint(w, "OK,RANDOMCODE124afio")
	}))
}

func TestSuccessParsingResponse(t *testing.T) {
	resp := bytes.NewBufferString("OK,12345678")
	if id, err := parseResponse(resp); id != "12345678" || err != nil {
		t.Fatalf("Expected 12345678, got %q", id)
	}
}

func TestFailedParsingResponse(t *testing.T) {
	resp := bytes.NewBufferString("FAILED,15,some random text")
	if id, err := parseResponse(resp); id != "" || err != ErrLongData {
		t.Fatalf("Expected nothing and ErrLongData, got %q and %q", id, err)
	}
}

func TestFailedUnknownParsingResponse(t *testing.T) {
	resp := bytes.NewBufferString("FAILED,70,some random text")
	if id, err := parseResponse(resp); id != "" || err == nil {
		t.Fatalf("Expected nothing and and dynamically created error, got %q and %q", id, err)
	}
}

func TestIntFromSlice(t *testing.T) {
	if num := intFromSlice([]byte("12")); num != 12 {
		t.Fatalf("Expected 12, got %v", num)
	}
}

func TestBasicSend(t *testing.T) {
	serv := createTestServer(defaultUser, defaultPass, "123456789")
	defer serv.Close()

	cl := NewClient(defaultUser, defaultPass)
	cl.address = serv.URL

	code, err := cl.Send("123456789", []byte("1234abcdefg"))

	if err != nil || code != "RANDOMCODE124afio" {
		t.Fatalf("Expected nil error and code 'RANDOMCODE124afio', got %v and %q", err, code)
	}
}

func TestBasicDefaultSend(t *testing.T) {
	serv := createTestServer(defaultUser, defaultPass, "123456789")
	defer serv.Close()

	cl := NewClient(defaultUser, defaultPass)
	cl.address = serv.URL

	cl.SetDefaultIMEI("123456789")
	code, err := cl.SendToDefault([]byte("1234abcdefg"))

	if err != nil || code != "RANDOMCODE124afio" {
		t.Fatalf("Expected nil error and code 'RANDOMCODE124afio', got %v and %q", err, code)
	}
}

func TestBasicStringSend(t *testing.T) {
	serv := createTestServer(defaultUser, defaultPass, "123456789")
	defer serv.Close()

	cl := NewClient(defaultUser, defaultPass)
	cl.address = serv.URL

	code, err := cl.SendString("123456789", "1234abcdefg")

	if err != nil || code != "RANDOMCODE124afio" {
		t.Fatalf("Expected nil error and code 'RANDOMCODE124afio', got %v and %q", err, code)
	}
}

func TestBasicStringDefaultSend(t *testing.T) {
	serv := createTestServer(defaultUser, defaultPass, "123456789")
	defer serv.Close()

	cl := NewClient(defaultUser, defaultPass)
	cl.address = serv.URL

	cl.SetDefaultIMEI("123456789")
	code, err := cl.SendStringToDefault("1234abcdefg")

	if err != nil || code != "RANDOMCODE124afio" {
		t.Fatalf("Expected nil error and code 'RANDOMCODE124afio', got %v and %q", err, code)
	}
}
