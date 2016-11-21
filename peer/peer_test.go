package peer

import (
	"net"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPeer(t *testing.T) {
	ts := httptest.NewUnstartedServer(nil)
	defer ts.Close()

	l, e := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, e)

	s := strings.Split(l.Addr().String(), ":")
	port := s[len(s)-1]

	ts.Listener = l
	ts.Start()

	p, err := NewPeer("127.0.0.1", port)
	assert.NoError(t, err)

	if assert.NotNil(t, p) {
		assert.Equal(t, "127.0.0.1", p.Target)
		assert.Equal(t, port, p.Port)
		assert.Equal(t, 1, p.health)
	}
}

func TestNewPeer_Err(t *testing.T) {
	p, err := NewPeer("1.2.3.4", "")
	assert.Nil(t, p)
	assert.Error(t, err)
}
