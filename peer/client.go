package peer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/trace"
)

type peerList struct {
	Actors []string `json:"actors"`
}

type Client struct {
	sync.RWMutex

	Peers []*Peer
}

func NewClient() *Client {
	return new(Client)
}

func (c *Client) ReceivePeers(s []string, i *int) error {
	var peers []*Peer
	for _, a := range s {
		log.Printf("[DEBUG] actor %v", a)
		p, err := NewPeer(a, "7777")
		if err != nil {
			log.Printf("[WARN] peer contact failed: %v", err)
			continue
		}
		peers = append(peers, p)
	}

	log.Printf("[INFO] %v new peers added", len(peers))
	c.Lock()
	c.Peers = peers
	c.Unlock()

	return nil
}

func (c *Client) HasPeers() bool {
	c.RLock()
	defer c.RUnlock()

	return len(c.Peers) > 0
}

func (c *Client) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	tr := trace.New("peer.ConfigHandler %q", r.URL.Path)
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
