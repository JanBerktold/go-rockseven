package rock7

import (
	"bytes"
	"code.google.com/p/intmath/intgr"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	user, pass string
	defIMEI    string
}

// NewClient a new client which is used for sending messages to deployed devices. Please note that the credentials are not checked upon creation of the Client object, but once the first request is triggered.
func NewClient(user, pass string) *Client {
	return &Client{
		user,
		pass,
		"",
	}
}

// SetDefaultIMEI sets IMEI for use with client.SendStringToDefault and client.SendToDefault.
func (cl *Client) SetDefaultIMEI(imei string) {
	cl.defIMEI = imei
}

func (cl *Client) Send(imei string, msg []byte) (string, error) {
	data := make([]byte, hex.EncodedLen(len(msg)))
	hex.Encode(data, msg)

	values := url.Values{}
	values.Add("imei", imei)
	values.Add("username", cl.user)
	values.Add("password", cl.pass)
	values.Add("data", string(data))

	resp, err := http.Post(sendURL, "", bytes.NewBufferString(values.Encode()))

	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
		return parseResponse(resp.Body)
	}
}

func (cl *Client) SendString(imei, msg string) (string, error) {
	return cl.Send(imei, []byte(msg))
}

func (cl *Client) SendStringToDefault(msg string) (string, error) {
	if len(cl.defIMEI) == 0 {
		return "", ErrDefaultSet
	}
	return cl.Send(cl.defIMEI, []byte(msg))
}

func (cl *Client) SendToDefault(msg []byte) (string, error) {
	if len(cl.defIMEI) == 0 {
		return "", ErrDefaultSet
	}
	return cl.Send(cl.defIMEI, msg)
}

func intFromSlice(by []byte) int {
	ret := 0
	for i := 0; i < len(by); i++ {
		ret += int(by[i]-'0') * intgr.Pow(10, len(by)-i-1)
	}
	return ret
}

// Not exactly nice, but yee
func parseResponse(read io.Reader) (string, error) {
	response := make([]byte, 200)
	n, _ := io.ReadFull(read, response)
	response = response[0:n]

	// checks for OK
	if response[0] != byte(79) || response[1] != byte(75) {
		errNum := intFromSlice(response[7:9])
		err, ok := MappedErrNum[errNum]
		if ok {
			return "", err
		} else {
			return "", fmt.Errorf("Unknown error %v: %q", errNum, string(response[10:]))
		}
	}

	return string(response[3:]), nil
}
