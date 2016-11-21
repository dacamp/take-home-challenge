// Package counter contains the necessary types used to create a
// grow only distributed counter.
package counter

import "sync"

// Args is a type used to pass a key value pair to during RPC calls
type Args struct {
	Key   string
	Value uint64
}

// GCounter is a grow-only counter
type GCounter struct {
	sync.RWMutex

	Counter map[string]uint64
}

// NewGCounter returns a new GCounter object
func NewGCounter() *GCounter {
	return &GCounter{Counter: make(map[string]uint64)}
}

// AddUint64 adds to the existing (*GCounter).Counter valued as
// defined by the key parameter.  It's not useful as an RPC call and
// should only be used by the local counter object
func (g *GCounter) AddUint64(key string, delta uint64) uint64 {
	g.Lock()
	defer g.Unlock()

	g.Counter[key] += delta
	return g.Counter[key]
}

// LoadUint64 loads the value of (*GCounter).Counter[key] and assigns
// that value to the *uint64 parameter.  An error must be specified to
// satisfy rpc.Call params.
func (g *GCounter) LoadUint64(key string, val *uint64) error {
	g.RLock()
	defer g.RUnlock()

	*val = g.Counter[key]
	return nil
}

// SetUint64 sets the value at (*GCounter).Counter[(*Args).Key] if it
// does not exist or it's value is less than (*Args).Value.  An
// optional *uint64 may be passed to satisfy the rpc.Call params.
func (g *GCounter) SetUint64(a *Args, val *uint64) error {
	g.Lock()
	defer g.Unlock()

	if v, ok := g.Counter[a.Key]; !ok || v < a.Value {
		g.Counter[a.Key] = a.Value
	}

	if val != nil {
		*val = g.Counter[a.Key]
	}
	return nil
}
