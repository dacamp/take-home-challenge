package peer

import (
	"net"
	"net/http/httptest"
	"net/rpc"
	"testing"

	"golang.org/x/net/trace"

	"github.com/stretchr/testify/assert"
)

func TestTraceAuthRequest(t *testing.T) {
	a, b := trace.AuthRequest(nil)
	assert.True(t, a)
	assert.False(t, b)
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.Equal(t, getLocalIP(), c.localIP)
	assert.False(t, c.HasPeers())
}

func TestClientReceivePeers(t *testing.T) {
	ts := httptest.NewUnstartedServer(nil)
	defer ts.Close()

	c := NewClient()
	rpc.Register(c)

	l, e := net.Listen("tcp", ":7777")
	defer l.Close()
	assert.NoError(t, e)

	ts.Listener = l
	ts.Start()

	var i int
	assert.NoError(t, c.ReceivePeers([]string{c.localIP, "abcdefg", "127.0.0.1"}, &i))
	assert.Equal(t, 2, i)

	if assert.True(t, c.HasPeers()) {
		assert.Equal(t, 0, c.pushConfig([]string{}))
		assert.Equal(t, 2, c.pushConfig([]string{c.localIP, "127.0.0.1", "abasc123"}))

		if assert.NotEmpty(t, c.Peers) {
			assert.NoError(t, c.Peers[0].Call("Client.ReceivePeers", []string{}, &i))
			assert.Equal(t, 1, i)

			assert.NotEmpty(t, c.Peers)
			assert.Error(t, c.Peers[0].Call("Client.FakeCall", []string{"127.0.0.1"}, nil))
		}
	}
}

func init() {
	rpc.HandleHTTP()
}
