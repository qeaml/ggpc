package ggpc_test

import (
	"io"
	"os"
	"testing"

	"github.com/qeaml/ggpc"
)

func assert(t *testing.T, b bool, m string) {
	if !b {
		t.Fatal(m)
	}
}

func assertNonerror(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestMemory(t *testing.T) {
	c := ggpc.NewMemory[string, int]()
	c.Set("a", 1)
	c.Set("B", 2)
	v, ok := c.Get("a")
	assert(t, ok, "value not present")
	assert(t, v == 1, "value incorrect")
	v = c.GetOrDefault("c", -1)
	assert(t, v == -1, "value incorrect")

	f, err := os.CreateTemp("", "ggpc_TestMemory")
	assertNonerror(t, err)
	_, err = f.WriteString(`{"a":1,"b":2}`)
	assertNonerror(t, err)
	_, err = f.Seek(0, 0)
	assertNonerror(t, err)
	c, err = ggpc.LoadMemory[string, int](f)
	assertNonerror(t, err)
	v, ok = c.Get("a")
	assert(t, ok, "value not present")
	assert(t, v == 1, "value incorrect")
	v = c.GetOrDefault("c", -1)
	assert(t, v == -1, "value incorrect")
	solid := c.Solidify()
	v, ok = solid["a"]
	assert(t, ok, "solid copy value not present")
	assert(t, v == 1, "solid copy value incorrect")
	c.Set("a", 99)
	v, ok = solid["a"]
	assert(t, ok, "solid copy value not present")
	assert(t, v == 1, "solid copy value incorrect")

	assertNonerror(t, f.Close())
}

func TestStored(t *testing.T) {
	f, err := os.CreateTemp("", "ggpc_TestStored")
	assertNonerror(t, err)

	c := ggpc.NewStored[string, int](f)
	c.Set("a", 1)
	c.Set("B", 2)
	v, ok := c.Get("a")
	assert(t, ok, "value not present")
	assert(t, v == 1, "value incorrect")
	v = c.GetOrDefault("c", -1)
	assert(t, v == -1, "value incorrect")
	assertNonerror(t, c.Save())

	_, err = f.Seek(0, 0)
	assertNonerror(t, err)
	data, err := io.ReadAll(f)
	assertNonerror(t, err)
	assert(t, string(data) == "{\"B\":2,\"a\":1}\n", "incorrect saved data")

	_, err = f.Seek(0, 0)
	assertNonerror(t, err)
	c, err = ggpc.LoadStored[string, int](f)
	assertNonerror(t, err)
	v, ok = c.Get("a")
	assert(t, ok, "value not present")
	assert(t, v == 1, "value incorrect")
	v = c.GetOrDefault("c", -1)
	assert(t, v == -1, "value incorrect")

	assertNonerror(t, f.Close())
}
