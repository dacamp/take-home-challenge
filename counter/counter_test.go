package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	c := NewGCounter()
	assert.Zero(t, c.Counter["foobar"])
}

func TestAddUint64(t *testing.T) {
	c := NewGCounter()

	assert.Equal(t, c.AddUint64("foobar", 5), uint64(5))
	assert.Equal(t, c.AddUint64("foobar", 5), uint64(10))
}

func TestLoadUint64(t *testing.T) {
	c := NewGCounter()

	assert.Nil(t, c.SetUint64(&Args{Key: "foobar", Value: 100}, nil))

	var val uint64
	assert.Nil(t, c.LoadUint64("foobar", &val))
	assert.Equal(t, uint64(100), val)
}

func TestSetUint64(t *testing.T) {
	c := NewGCounter()

	assert.Nil(t, c.SetUint64(&Args{Key: "foobar", Value: 100}, nil))

	var val uint64
	assert.Nil(t, c.SetUint64(&Args{Key: "foobar", Value: 5}, &val))
	assert.Equal(t, uint64(100), val)
}
