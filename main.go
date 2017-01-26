package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/dustin/go-coap"
)

// check is a simple wrapper around the verbose
// `if err != nil` checks.
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Metric is the layout of the metric payload from a client.
type Metric struct {
	Fqdn    uint16
	Service uint8
	Key     uint16
	Value   string
}

func main() {

	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		if VersionPrerelease != "" {
			fmt.Println("Version PreRelease:", VersionPrerelease)
		}
		return
	}

	coap.ListenAndServe("udp", ":5683",
		coap.FuncHandler(func(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
			fmt.Printf("Go message path=%q: %#v from %v\n", m.Path(), m, a)
			fmt.Printf("Payload: %#v\n", m.Payload)
			if m.IsConfirmable() {
				res := &coap.Message{
					Type:      coap.Acknowledgement,
					Code:      coap.Content,
					MessageID: m.MessageID,
					Token:     m.Token,
					Payload:   []byte("Halp"),
				}
				res.SetOption(coap.ContentFormat, coap.TextPlain)

				return res
			}
			return nil
		}),
	)
}
