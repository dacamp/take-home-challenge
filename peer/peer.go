// Package peer provides access to remote peers
package peer

import (
	"log"
	"net/rpc"

	"golang.org/x/net/trace"
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

// Call is a wrapper around rpc.Call, allowing for peer health to be
// evaluated
func (p *Peer) Call(serviceMethod string, args interface{}, reply interface{}) error {
	tr := trace.New("peer.Call", p.Target)
	defer tr.Finish()

	if err := p.Client.Call(serviceMethod, args, reply); err != nil {
		tr.LazyPrintf("%v rpc failed: %v", serviceMethod, err)
		tr.SetError()

		log.Printf("[ERROR] %v rpc failed: %v:", serviceMethod, err)

		return err
	}
	tr.LazyPrintf("%v rpc success", serviceMethod)
	return nil
}
