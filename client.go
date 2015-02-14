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

// Client is used for sending messages to deployed devices. Should only be created using NewClient.
type Client struct {
	user, pass string
	defIMEI    string
	address    string
}

// NewClient creates a new Client which is used for sending messages to deployed devices. Please note that the credentials are not checked upon creation of the Client object, but once the first request is triggered.
func NewClient(user, pass string) *Client {
	return &Client{
		user,
		pass,
		"",
		sendURL,
	}
}

// SetDefaultIMEI sets IMEI for use with Client.SendStringToDefault and Client.SendToDefault.
func (cl *Client) SetDefaultIMEI(imei string) {
	cl.defIMEI = imei
}

// Send performs the task of transmitting a message to a rockseven device.
// Returns the unique id which has been assigned to the message and is also available for the receiving device.
// Errors can occur in the case of connection issue between the client and rockseven's
// servers or in the case of problems on rockseven's site. The returned error object
// gives more detailed information in each case.
func (cl *Client) Send(imei string, msg []byte) (string, error) {
	data := make([]byte, hex.EncodedLen(len(msg)))
	hex.Encode(data, msg)

	values := url.Values{}
	values.Add("imei", imei)
	values.Add("username", cl.user)
	values.Add("password", cl.pass)
	values.Add("data", string(data))

	resp, err := http.Post(cl.address, "application/x-www-form-urlencoded", bytes.NewBufferString(values.Encode()))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	return parseResponse(resp.Body)
}

// SendString is convenience wrapper for Send which takes the string representation
// of a message instead of a byte slice.
func (cl *Client) SendString(imei, msg string) (string, error) {
	return cl.Send(imei, []byte(msg))
}

// SendStringToDefault's behaviour is similar to SendString, however unlike SendString which sends
// its message to a specified IMEI, this method sends the message to the default IMEI number.
func (cl *Client) SendStringToDefault(msg string) (string, error) {
	if len(cl.defIMEI) == 0 {
		return "", ErrDefaultSet
	}
	return cl.Send(cl.defIMEI, []byte(msg))
}

// SendToDefault's behaviour is similar to Send, however unlike Send which targets a specified IMEI,
// this method sends the message to the default IMEI number.
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
	if len(response) < 2 || response[0] != byte(79) || response[1] != byte(75) {
		errNum := intFromSlice(response[7:9])
		err, ok := MappedErrNum[errNum]
		if ok {
			return "", err
		}
		var stringResp string
		if len(response) > 10 {
			stringResp = string(response[10:])
		} else {
			stringResp = ""
		}
		return "", fmt.Errorf("unknown error %v: %q", errNum, stringResp)
	}

	return string(response[3:]), nil
}
