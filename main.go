package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

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

var tree []Fqdn

// Fqdn is the base of the tree which links to `Service`s.
// ie. com.example.api
type Fqdn struct {
	Name     string
	Index    uint16
	Services []Service
}

// Service is the name of the service, which links to `Key`s.
// ie. my_super_service
type Service struct {
	Name  string
	Index uint8
	Keys  []Key
}

// Key is the metric name.
// ie. endpoint_name.method.status_code
type Key struct {
	Name  string
	Index uint16
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
			fmt.Printf("Go message path=%q: %#v from %v\n", m.PathString(), m, a)
			fmt.Printf("Payload: %#v\n\n", m.Payload)
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			if "register_key" == m.PathString() {
				metricKey := registerKeyspace(m.Payload)

				res.Payload = metricKey
				// res := &coap.Message{
				// 	Type:      coap.Acknowledgement,
				// 	Code:      coap.Content,
				// 	MessageID: m.MessageID,
				// 	Token:     m.Token,
				// 	Payload:   metricKey,
				// }

			} else {
				res.Payload = []byte("Halp")
				// res := &coap.Message{
				// 	Type:      coap.Acknowledgement,
				// 	Code:      coap.Content,
				// 	MessageID: m.MessageID,
				// 	Token:     m.Token,
				// 	Payload:   []byte("Halp"),
				// }
				// res.SetOption(coap.ContentFormat, coap.TextPlain)
			}
			if m.IsConfirmable() {
				return res
			}
			return nil
		}),
	)
}

// KeyValue states the index of each part of the tree
type KeyValue struct {
	FqdnIndex    uint16
	ServiceIndex uint8
	KeyIndex     uint16
}

var reverseMap = make(map[KeyValue]string)

type F struct {
	Index    uint16
	Services map[string]S
}

type S struct {
	Index uint8
	Keys  map[string]K
}

type K struct {
	Index uint16
}

var keyspace = make(map[string]F)

// Keyspace := make(map[string]map[string]map[string]KeyValue)
// var Fqdnspace map[]
// rs map[string]map[time.Time]int
// map[struct{s string, t time.Time}]int

func registerKeyspace(key []byte) []byte {
	names := strings.Split(string(key), " ")
	fmt.Println(names)

	f, ok := keyspace[names[0]]
	if ok == true {
		s, ok := f.Services[names[1]]
		if ok == true {
			_, ok := s.Keys[names[2]]
			if ok == false {
				keyspace[names[0]].Services[names[1]].Keys[names[2]] = K{uint16(len(s.Keys) + 1)}
			}
		} else {
			keyspace[names[0]].Services[names[1]] = S{
				uint8(len(f.Services) + 1),
				map[string]K{
					names[2]: K{
						uint16(1),
					},
				},
			}
		}
	} else {
		keyspace[names[0]] = F{
			uint16(len(keyspace) + 1),
			map[string]S{
				names[1]: S{
					uint8(1),
					map[string]K{
						names[2]: K{
							uint16(1),
						},
					},
				},
			},
		}
	}

	return []byte(fmt.Sprintf("%v%v%v", keyspace[names[0]].Index,
		keyspace[names[0]].Services[names[1]].Index,
		keyspace[names[0]].Services[names[1]].Keys[names[2]].Index))

	// _, ok := Keyspace[names[0]]
	// if ok == nil {
	// 	_, ok = Keyspace[names[0]][names[1]]
	// 	if ok == nil {
	// 		key, ok = Keyspace[names[0]][names[1]][names[2]]
	// 		if ok == nil {
	// 			return []byte(fmt.Sprintf("%v%v%v", key.FqdnIndex, key.ServiceIndex, key.KeyIndex)
	// 		} else {
	// 			Keyspace[names[0]][names[1]][names[2]] = KeyValue{
	// 				FqdnIndex:
	// 			}
	// 		}
	// 	}
	// }

	// var fqdnIndex, keyIndex uint16
	// var serviceIndex uint8
	// for fqdn := range tree {
	// 	if fqdn.Name == names[0] {
	// 		fqdnIndex = fqdn.Index

	// 		for service := range fqdn.Services {
	// 			if service.Name == names[1] {
	// 				serviceIndex = service.Index

	// 				for key := range service.Keys {
	// 					if key.Name == names[2] {
	// 						keyIndex = key.Index
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// f := Fqdn{
	// 	Name:  "com.example.api",
	// 	Index: 0x0000001,
	// 	Services: []Service{
	// 		Service{
	// 			Name:  "my_super_service",
	// 			Index: 0x00000001,
	// 			Keys: []Key{
	// 				Key{
	// 					Name:  "endpoint_name.method.status_code",
	// 					Index: 0x00000001,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// tree = append(tree, f)

	// return []byte(fmt.Sprintf("%v.%v.%v", f.Name, f.Services[0].Name, f.Services[0].Keys[0].Name))
}
