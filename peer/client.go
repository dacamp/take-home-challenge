package peer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/trace"
)

// Client handles remote peer coordination
type Client struct {
	sync.RWMutex

	Peers []*Peer
}

// peerList is a type made specifically for this challenge to accept
// an incoming JSON POST
type peerList struct {
	Actors []string `json:"actors"`
}

// NewClient returns a new Client object
func NewClient() *Client {
	return new(Client)
}

// ReceivePeers is an RPC function takes an incoming slice of peers
func (c *Client) ReceivePeers(s []string, i *int) error {
	var peers []*Peer
	for _, a := range s {
		tr := trace.New("peer.ReceivePeers", a)
		defer tr.Finish()

		log.Printf("[DEBUG] actor %v", a)
		p, err := NewPeer(a, "7777")
		if err != nil {
			tr.LazyPrintf("peer %q contact failed: %v", a, err)
			log.Printf("[WARN] peer %q contact failed: %v", a, err)
			continue
		}
		tr.LazyPrintf("peer %v added", p.Target)
		peers = append(peers, p)
	}

	c.Lock()
	c.Peers = peers
	c.Unlock()
	log.Printf("[INFO] added %v new peers", len(peers))
	if i != nil {
		*i = len(peers)
	}

	return nil
}

// HasPeers is a helper function testing for the presents of available
// remote peers
//
// TODO: validate the health of those peers
func (c *Client) HasPeers() bool {
	c.RLock()
	defer c.RUnlock()

	return len(c.Peers) > 0
}

// ConfigHandler is an http.HandlerFunc that accepts a JSON blob of
// actors, validates the request and on success, broadcasts the
// configs to peers.
//
// NOTE: this is last write wins, so conflicting configs sent to
// different nodes may produce unexpected results.
// TODO: fix this^
func (c *Client) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	tr := trace.New("peer.ConfigHandler", r.URL.Path)
	defer tr.Finish()

	if r.Method != "POST" {
		tr.LazyPrintf("method not allowed %v %v", r.Method, r.URL.Path)
		tr.SetError()
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var peerInput peerList
	if err = json.Unmarshal(body, &peerInput); err != nil {
		tr.LazyPrintf("json.Unmarshal failed: %v", err)
		tr.SetError()
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// if we care, we can return the results of this data
	// to the user, maybe with a param?
	go c.pushConfig(peerInput.Actors)

	// return 200
}

// should be async
func (c *Client) pushConfig(s []string) {
	log.Println("[DEBUG] Within PushConfig")
	for _, a := range s {
		log.Printf("[DEBUG] actor %v", a)
		p, err := NewPeer(a, "7777")
		if err != nil {
			log.Printf("[WARN] peer contact failed: %v", err)
			continue
		}
		log.Printf("[DEBUG] Peer: %+v", p)
		var i int
		if err := p.Call("Client.ReceivePeers", s, &i); err != nil {
			log.Fatal("Client.ReceivePeers [%d] error:", i, err)
		}
	}
}
