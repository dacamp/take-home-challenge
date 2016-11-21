package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/dacamp/challenge/peer"
)

func main() {
	client := peer.NewClient()
	rpc.Register(client)

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":7777")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	http.HandleFunc("/config", client.ConfigHandler)
	http.HandleFunc("/counter/", client.CounterHandler())

	http.Serve(l, nil)
}

func init() {
	if t := os.Getenv("CHALLENGE_TIMEOUT"); t != "" {
		i, err := strconv.Atoi(t)
		if err != nil {
			log.Printf("[WARN] unable to convert %v: %v", t, err)
			return
		}
		peer.DefaultTimeout = time.Duration(i)
	}
}
