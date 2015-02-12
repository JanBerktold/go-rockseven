package rock7

import (
	"bytes"
	"testing"
)

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
