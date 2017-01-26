package main

import (
	"fmt"
	"log"

	"github.com/dustin/go-coap"
)

// Metric is the layout of the metric payload from a client.
type Metric struct {
	Fqdn    uint16
	Service uint8
	Key     uint16
	Value   string
}

// Pack converts the struct into a byte slice for sending on the wire
func (m Metric) Pack() []byte {
	return []byte(fmt.Sprintf("%v%v%v%v", m.Fqdn, m.Service, m.Key, m.Value))
}

func main() {
	// req := coap.Message{
	// 	Type:      coap.Confirmable,
	// 	Code:      coap.GET,
	// 	MessageID: 12345,
	// 	Payload:   []byte("hello, world!"),
	// }

	// req.SetOption(coap.ETag, "weetag")
	// req.SetOption(coap.MaxAge, 3)
	// req.SetPathString("/some/path")

	// c, err := coap.Dial("udp", "localhost:5683")
	// if err != nil {
	// 	log.Fatalf("Error dialing: %v", err)
	// }

	// rv, err := c.Send(req)
	// if err != nil {
	// 	log.Fatalf("Error sending request: %v", err)
	// }

	// if rv != nil {
	// 	log.Printf("Response payload: %s", rv.Payload)
	// }

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
		Payload:   m.Pack(),
	}

	c, err := coap.Dial("udp", "localhost:5683")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if rv != nil {
		log.Printf("Response payload: %s", rv.Payload)
	}

}
