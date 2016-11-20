package peer

import (
	"log"
	"net/rpc"
)

type Peer struct {
	*rpc.Client
	Target string `json:"peer,omitempty"`
	Port   string `json:"port,omitempty"`

	// reserved for future use
	health int
}

func NewPeer(host, port string) (*Peer, error) {
	client, err := rpc.DialHTTP("tcp", host+":"+port)
	if err != nil {
		log.Printf("[ERROR] rpc.DialHTTP failed: %v", err)
		return nil, err
	}

	return &Peer{
		Client: client,
		Target: host,
		Port:   port,
	}, nil
}
