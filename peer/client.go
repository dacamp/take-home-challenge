package peer

import (
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

var (
	// DefaultTimeout is the allowed timeout while making RPC calls to peers
	DefaultTimeout = time.Duration(5)
	routeRegexp    = regexp.MustCompile("^/counter/([a-zA-Z0-9]+)/?((?:consistent_)?value)?$")
)

// Client handles remote peer coordination
type Client struct {
	sync.RWMutex

	localIP string
	Peers   []*Peer
}

// peerList is a type made specifically for this challenge to accept
// an incoming JSON POST
type peerList struct {
	Actors []string `json:"actors"`
}

// NewClient returns a new Client object
func NewClient() *Client {
	return &Client{
		localIP: getLocalIP(),
	}
}

// ReceivePeers is an RPC function takes an incoming slice of peers
func (c *Client) ReceivePeers(s []string, i *int) error {
	var peers []*Peer

	for _, a := range s {
		if a == c.localIP {
			continue
		}

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

	p, err := NewPeer(c.localIP, "7777")
	if err != nil {
		log.Printf("[WARN] local RPC endpoint %q contact failed: %v", c.localIP, err)
	} else {
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

// TODO use async rpc.Client.Go
func (c *Client) pushConfig(s []string) int {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout*time.Second)

	// Even though ctx should have expired already, it is good
	// practice to call its cancelation function in any case.
	// Failure to do so may keep the context and its parent alive
	// longer than necessary.
	defer cancel()

	log.Println("[DEBUG] Within PushConfig")

	intCh := make(chan int)

	var wg sync.WaitGroup
	for _, a := range s {
		wg.Add(1)

		log.Printf("[DEBUG] actor %v", a)
		go func(a string) {
			defer wg.Done()

			p, err := NewPeer(a, "7777")
			if err != nil {
				log.Printf("[WARN] peer contact failed: %v", err)
				return
			}
			log.Printf("[DEBUG] Peer: %+v", p)

			var i int
			if err := p.Call("Client.ReceivePeers", s, &i); err != nil {
				log.Printf("[ERROR] Client.ReceivePeers [%d] error: %v", i, err)
			}
			intCh <- i
		}(a)
	}

	d := make(chan struct{})
	go func() {
		defer close(d)
		defer close(intCh)
		wg.Wait()
	}()

	var i int
	go func() {
		for x := range intCh {
			if i < x {
				i = x
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("[WARN] timeout before operation completed")
	case <-d:
		// fall through
	}

	return i
}

// getLocalIP returns the non loopback local IP of the host
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func init() {
	// This is due to ACLs limitations within the trace package
	// requiring that authorized requests come from localhost
	// only.  I'm sure there's some docker magic to properly mask
	// requests, but this was an easier hack.
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, false
	}
}
