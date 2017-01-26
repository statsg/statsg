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

// KeyValue states the index of each part of the tree
type KeyValue struct {
	FqdnIndex    uint8
	ServiceIndex uint8
	KeyIndex     uint8
}

var reverseMap = make(map[KeyValue]string) // A reverse mapping of the index keys to the metric key name

// F is the FQDN struct
type F struct {
	Index    uint8
	Services map[string]S
}

// S is the Service struct
type S struct {
	Index uint8
	Keys  map[string]K
}

// K is the Key struct
type K struct {
	Index uint8
}

var keyspace = make(map[string]F) // The root of the FQDN->Services->Keys mapping

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
			fmt.Printf("Payload: %#v\n", m.Payload)
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			if "register_key" == m.PathString() {
				metricKey := registerKeyspace(m.Payload)
				fmt.Printf("MetricKey: %#v\n\n", metricKey)

				res.Payload = metricKey
			} else {
				res.Payload = []byte("Halp")
			}
			if m.IsConfirmable() {
				return res
			}
			return nil
		}),
	)
}

// registerKeyspace takes the incoming request payload
// and ascertains if it already exists in the mapping,
// else it will create the FQDN->Services->Keys mapping
// as well as the reverse lookup of
// {FQDN, Service, Key} => Full Metric Name
func registerKeyspace(key []byte) []byte {
	names := strings.Split(string(key), " ")
	fmt.Println(names)

	f, ok := keyspace[names[0]]
	if ok == true {
		s, ok := f.Services[names[1]]
		if ok == true {
			_, ok := s.Keys[names[2]]
			if ok == false {
				// Only the metric `Key` does not currently exist, so create it
				keyspace[names[0]].Services[names[1]].Keys[names[2]] = K{uint8(len(s.Keys) + 1)}
			}
		} else {
			// The `Service` and `Key` do not exist, so create them
			keyspace[names[0]].Services[names[1]] = S{
				uint8(len(f.Services) + 1),
				map[string]K{
					names[2]: K{
						uint8(1),
					},
				},
			}
		}
	} else {
		// The `Fqdn`, `Service`, and `Key` do not exist, so create them
		keyspace[names[0]] = F{
			uint8(len(keyspace) + 1),
			map[string]S{
				names[1]: S{
					uint8(1),
					map[string]K{
						names[2]: K{
							uint8(1),
						},
					},
				},
			},
		}
	}

	kv := KeyValue{
		FqdnIndex:    keyspace[names[0]].Index,
		ServiceIndex: keyspace[names[0]].Services[names[1]].Index,
		KeyIndex:     keyspace[names[0]].Services[names[1]].Keys[names[2]].Index,
	}
	reverseMap[kv] = strings.Join(names, ".")

	return []byte(fmt.Sprintf("%v%v%v", kv.FqdnIndex, kv.ServiceIndex, kv.KeyIndex))
}
