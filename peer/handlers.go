package peer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"sync"
	"time"

	"github.com/dacamp/challenge/counter"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

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

// CounterHandler is an http.HandlerFunc with three valid endpoints
//   - GET /counter/:name:/value
//   - GET /counter/:name:/consistent_value
//   - POST /counter/:name:
func (c *Client) CounterHandler() func(w http.ResponseWriter, r *http.Request) {
	gCounter := counter.NewGCounter()
	rpc.Register(gCounter)
	return func(w http.ResponseWriter, r *http.Request) {
		tr := trace.New("peer.CounterHandler", r.URL.Path)
		ctx := trace.NewContext(context.Background(), tr)
		defer tr.Finish()

		path := routeRegexp.FindStringSubmatch(r.URL.Path)
		if path == nil {
			routeError(ctx, w, r)
			return
		}

		var buf bytes.Buffer
		var val uint64
		key := path[1]

		switch r.Method {
		case "GET":
			gCounter.LoadUint64(key, &val)

			switch path[2] {
			case "consistent_value":
				if err := c.setPeerValues(ctx, key, &val); err != nil {
					fmt.Fprintf(&buf, `{"method": "%v", "error": "%v", "best-guess": "%d"}`, r.Method, err, val)
					http.Error(w, buf.String(), http.StatusInternalServerError)
					return
				}
			case "value":
				// fall through
			default:
				routeError(ctx, w, r)
				return
			}

			fmt.Fprintln(&buf, val)
			w.Write(buf.Bytes())
		case "POST":
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err = json.Unmarshal(body, &val); err != nil {
				http.Error(w, `{"method": "`+r.Method+`", "error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}

			val = gCounter.AddUint64(key, val)
			c.setPeerValues(ctx, key, &val)
		}
	}

	// returns 200
}

func (c *Client) setPeerValues(ctx context.Context, key string, val *uint64) error {
	if !c.HasPeers() {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout*time.Second)

	// Even though ctx should have expired already, it is good
	// practice to call its cancelation function in any case.
	// Failure to do so may keep the context and its parent alive
	// longer than necessary.
	defer cancel()

	var wg sync.WaitGroup
	for _, p := range c.Peers {
		wg.Add(1)

		go func(p *Peer) {
			defer wg.Done()
			p.Call("GCounter.SetUint64", &counter.Args{
				Key:   key,
				Value: *val,
			}, val)
		}(p)
	}

	d := make(chan struct{})
	go func() {
		defer close(d)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		log.Println("[WARN] timeout before operation completed")
		return ctx.Err()
	case <-d:
		// fall through
	}

	return nil
}

func routeError(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("no route found %v %q", r.Method, r.URL.Path)
	}
	http.Error(w, `{"method": "`+r.Method+`", "error": "No route found for '`+r.URL.Path+`'"}`, http.StatusBadRequest)
}
