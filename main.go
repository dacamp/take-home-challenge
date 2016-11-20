package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"regexp"
	"sync"
	"time"

	"golang.org/x/net/trace"

	"github.com/dacamp/challenge/counter"
)

var (
	routeRegexp = regexp.MustCompile("^/counter/([a-zA-Z0-9]+)/?((?:consistent_)?value)?$")
)

type peerList struct {
	Actors []string `json:"actors"`
}

func main() {
	gCounter := counter.NewGCounter()
	rpc.Register(gCounter)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":7777")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	var peers peerList
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err = json.Unmarshal(body, &peers); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		// return 200
	})

	http.HandleFunc("/counter/", func(w http.ResponseWriter, r *http.Request) {
		tr := trace.New("main.counterHandler", r.URL.Path)
		defer tr.Finish()
		// tr.LazyPrintf("some event %q happened", str)
		// if err := somethingImportant(); err != nil {
		// 	tr.LazyPrintf("somethingImportant failed: %v", err)
		// 	tr.SetError()
		// }

		path := routeRegexp.FindStringSubmatch(r.URL.Path)
		if path == nil {
			http.Error(w, `{"method": "`+r.Method+`", "error": "No route found for '`+r.URL.Path+`'"}`, http.StatusBadRequest)
			tr.LazyPrintf("no route found %v %q", r.Method, r.URL.Path)
			return
		}

		switch r.Method {
		case "GET":
			var buf bytes.Buffer
			var val uint64

			switch path[2] {
			case "value":
				gCounter.LoadUint64(path[1], &val)
			case "consistent_value":
				if peers.Actors == nil {
					gCounter.LoadUint64(path[1], &val)
					break
				}

				gCounter.LoadUint64(path[1], &val)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

				// Even though ctx should have expired already, it is good
				// practice to call its cancelation function in any case.
				// Failure to do so may keep the context and its parent alive
				// longer than necessary.
				defer cancel()

				var wg sync.WaitGroup
				for _, a := range peers.Actors {
					wg.Add(1)

					go func(a string) {
						defer wg.Done()
						tr.LazyPrintf("contacting remote peer %q", a)

						client, err := rpc.DialHTTP("tcp", a+":7777")
						if err != nil {
							tr.LazyPrintf("rpc.DialHTTP failed: %v", err)
							tr.SetError()

							log.Printf("[ERROR] rpc.DialHTTP failed: %v", err)
						}

						if err = client.Call("GCounter.SetUint64", &counter.Args{
							Key:   path[1],
							Value: val,
						}, &val); err != nil {
							tr.LazyPrintf("counter.Setuint64 rpc failed: %v", err)
							tr.SetError()

							log.Printf("[ERROR] counter.Setuint64 rpc failed: %v:", err)
						}
					}(a)
				}
				wg.Wait()

				d := make(chan struct{})
				go func() {
					defer close(d)
					wg.Wait()
				}()

				select {
				case <-ctx.Done():
					log.Println("[WARN] timeout before operation completed")
					http.Error(w, `{"method": "`+r.Method+`", "error": "`+ctx.Err().Error()+`"}`, http.StatusInternalServerError)
					return
				case <-d:
					// fall through
				}
				// gCounter.SetUint64(&counter.Args{
				//	Key:   path[1],
				//	Value: val,
				//}, &val)
			default:
				http.Error(w, `{"method": "`+r.Method+`", "error": "No route found for '`+r.URL.Path+`'"}`, http.StatusBadRequest)
				return
			}
			fmt.Fprintln(&buf, val)
			w.Write(buf.Bytes())
		case "POST":
			var c uint64
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err = json.Unmarshal(body, &c); err != nil {
				http.Error(w, `{"method": "`+r.Method+`", "error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			newVal := gCounter.AddUint64(path[1], c)

			args := &counter.Args{
				Key:   path[1],
				Value: newVal,
			}
			var wg sync.WaitGroup
			for _, a := range peers.Actors {
				wg.Add(1)

				go func(a string) {
					defer wg.Done()

					client, err := rpc.DialHTTP("tcp", a+":7777")
					if err != nil {
						log.Fatal("dialing:", err)
					}

					var v uint64
					if err = client.Call("GCounter.SetUint64", args, &v); err != nil {
						log.Fatal("counter.LoadUint64 error:", err)
					}
				}(a)
			}
			wg.Wait()
		}
	})

	http.Serve(l, nil)
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
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
