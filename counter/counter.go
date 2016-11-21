package counter

import "sync"

type Args struct {
	Key   string
	Value uint64
}

type GCounter struct {
	sync.RWMutex

	Counter map[string]uint64
}

func NewGCounter() *GCounter {
	return &GCounter{Counter: make(map[string]uint64)}
}

func (g *GCounter) AddUint64(key string, delta uint64) uint64 {
	g.Lock()
	defer g.Unlock()

	g.Counter[key] += delta
	return g.Counter[key]
}

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

func (g *GCounter) LoadUint64(key string, val *uint64) error {
	g.RLock()
	defer g.RUnlock()

	*val = g.Counter[key]
	return nil
}
