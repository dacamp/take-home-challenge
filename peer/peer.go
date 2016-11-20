// Package peer provides access to remote peers
package peer

import (
	"log"
	"net/rpc"
)

// Peer reflects all the remote peers, which may include itself
type Peer struct {
	*rpc.Client
	Target string `json:"peer,omitempty"`
	Port   string `json:"port,omitempty"`

	// reserved for future use
	health int
}

// NewPeer attempts to establish a new connection to a peer.  If that
// connection fails, an error is returned.
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
		health: 1,
	}, nil
}
