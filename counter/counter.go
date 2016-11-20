package counter

import (
	"crypto/rand"
	"fmt"
	"sync"
)

type GCounter struct {
	sync.RWMutex

	ident   string
	Counter map[string]uint64
}

func NewGCounter() *GCounter {
	return &GCounter{
		ident:   pseudo_uuid(),
		Counter: make(map[string]uint64),
	}
}

func (g *GCounter) AddUint64(key string, delta uint64) uint64 {
	g.Lock()
	defer g.Unlock()

	g.Counter[key] += delta
	return g.Counter[key]
}

type Args struct {
	Key   string
	Value uint64
}

func (g *GCounter) SetUint64(a *Args, val *uint64) error {
	g.Lock()
	defer g.Unlock()

	if v, ok := g.Counter[a.Key]; !ok || v < a.Value {
		g.Counter[a.Key] = a.Value
	}

	*val = g.Counter[a.Key]
	return nil
}

func (g *GCounter) Set(key string, val uint64) {
	g.Lock()
	defer g.Unlock()

	if v, ok := g.Counter[key]; !ok || v < val {
		g.Counter[key] = val
	}
}

func (g *GCounter) Merge(a *GCounter, val *int64) error {
	g.RLock()
	defer g.RUnlock()

	for i, rv := range a.Counter {
		if lv, ok := g.Counter[i]; !ok || lv < rv {
			g.Counter[i] = rv
		}
	}
	return nil
}

func (g *GCounter) LoadUint64(key string, val *uint64) error {
	g.RLock()
	defer g.RUnlock()

	*val = g.Counter[key]
	return nil
}

func pseudo_uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
