package main

import (
	"fmt"
	"log"

	"github.com/dustin/go-coap"
)

// Metric is the layout of the metric payload from a client.
type Metric struct {
	Fqdn    uint8
	Service uint8
	Key     uint8
	Value   string
}

// Pack converts the struct into a byte slice for sending on the wire
func (m Metric) Pack() []byte {
	return []byte(fmt.Sprintf("%v%v%v%v", m.Fqdn, m.Service, m.Key, m.Value))
}

func main() {
	m := Metric{
		Fqdn:    0x00000001,
		Service: 0x00000001,
		Key:     0x00000001,
		Value:   "Metric",
	}

	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 1,
		Payload:   []byte("com.example.api my_super_service endpoint_name.method.status_code"),
	}

	req.SetPathString("register_key")

	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s :: %#v", rv.Payload, rv.Payload)
	}

	m = Metric{
		Fqdn:    rv.Payload[0],
		Service: rv.Payload[1],
		Key:     rv.Payload[2],
		Value:   "T",
	}

	req = coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 1,
		Payload:   m.Pack(),
	}

	rv, err = c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s", rv.Payload)
	}

}
