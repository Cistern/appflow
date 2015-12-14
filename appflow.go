// Package appflow provides middleware to export HTTP application flow data.
package appflow

import (
	"encoding/json"
	"net"
	"net/http"
)

// Destination represents a flow destination.
type Destination struct {
	conn *net.UDPConn
}

// HTTPFlowData represents flow data that will be
// serialized and sent to a collector.
type HTTPFlowData struct {
	Method        string              `json:"method"`
	URL           string              `json:"url"`
	Proto         string              `json:"proto"`
	Header        map[string][]string `json:"header"`
	ContentLength int                 `json:"contentLength"`
	Host          string              `json:"host"`
	RemoteAddr    string              `json:"remoteAddr"`
}

// NewDestination creates a new Destination that sends
// flow data to the given address.
func NewDestination(address string) (*Destination, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return &Destination{
		conn: conn,
	}, nil
}

// Emit serializes an http.Request and sends
// it to the Destination.
func (d *Destination) Emit(r *http.Request) {
	flowData := HTTPFlowData{
		Method:        r.Method,
		URL:           r.URL.String(),
		Proto:         r.Proto,
		Header:        map[string][]string(r.Header),
		ContentLength: int(r.ContentLength),
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
	}
	json.NewEncoder(d.conn).Encode(flowData)
}

// Decode unmarshals JSON-encoded HTTPFlowData from b.
func Decode(b []byte) (*HTTPFlowData, error) {
	var flowData HTTPFlowData
	err := json.Unmarshal(b, &flowData)
	if err != nil {
		return nil, err
	}
	return &flowData, nil
}
