package gpsd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// Package comment in here; loosely modelled on the http pacakge

// GpsdResponse represents the name of a valid, supported gpsd response
type GpsdResponse string

// HandlerFunc defines the handler for a given gpsd response
type HandlerFunc func(interface{})

const (
	// Currently supported gpsd responses
	Tpv     GpsdResponse = "TPV"
	Skye    GpsdResponse = "SKY"
	Version GpsdResponse = "VERSION"
	Watch   GpsdResponse = "WATCH"
)

var (
	socket   net.Conn
	reader   *bufio.Reader
	handlers = make(map[GpsdResponse]HandlerFunc)
	mutex    sync.Mutex // For safe access to the handlers map
)

// ConnectAndListen initiates the connection to the gpsd instance specified by address, and
// sends a gpsd WATCH command. Listening continues until the client explicitly calls Close().
func ConnectAndListen(address string) error {
	socket, err := net.Dial("tcp4", address)
	if err != nil {
		return err
	}
	fmt.Printf("%v", socket)
	reader = bufio.NewReader(socket)
	reader.ReadString('\n')

	go listen(socket, reader)

	return nil // Hmmm.
}

// Close stops the listener go routine, and closes the connection to gpsd.
func Close() {
	// TODO: Tell the listen go routine to stop, via channel

	socket.Close()
}

// HandleFunc registers the supplied HandlerFunc for the given GpsdResponse. Multiple
// calls to HandleFunc with the same GpsdResposne will overwrite the previous value.
func HandleFunc(response GpsdResponse, handler HandlerFunc) {
	if response == "" {
		panic("gpsd: nil response")
	}
	if handler == nil {
		panic("gpsd: nil handler")
	}
	mutex.Lock()
	handlers[response] = handler
	mutex.Unlock()
}

// listen is run in a go routine; needs local access to the socket and reader
func listen(socket net.Conn, reader *bufio.Reader) {
	fmt.Println("Listening on socket...")
	// TODO: Hard-coded at the moment
	fmt.Fprintf(socket, "?WATCH={\"enable\":true, \"json\":true}")

	// We'll go off and await any messages that come through until told to stop
	for {
		// Read response (comes as a single line of JSON)
		line, err := reader.ReadString('\n')
		if err != nil {
			// Handle stream reader error
			fmt.Printf("gpsd: error reading from socket: %v\n", err)
		}

		// Check whether we're interested in this response (i.e. whether
		// the client has specified a handler for it)
		var check Report
		lineBytes := []byte(line)
		if err = json.Unmarshal(lineBytes, &check); err != nil {
			// Handle initial parsing error
		}
		if handlers[check.Class] == nil {
			continue
		}
		// If handler map != empty, check response has a handler
		if response, err := unmarshalResponse(check.Class, lineBytes); err == nil {
			handlers[check.Class](response)
		}
	}
}

func unmarshalResponse(class GpsdResponse, bytes []byte) (interface{}, error) {
	var err error

	switch class {
	case Watch:
		var rsp *WatchResponse
		fmt.Printf("%s\n", bytes)
		err = json.Unmarshal(bytes, &rsp)
		return rsp, err
	case Tpv:
		var rsp *TimePositionVelocityResponse
		err = json.Unmarshal(bytes, &rsp)
		return rsp, err
	}
	// Others as required
	return nil, err
}
