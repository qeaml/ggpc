package ggpc

import (
	"encoding/json"
	"io"
	"sync"
)

// Cache represents any cache, in-memory or persistent via JSON
type Cache[K comparable, V any] struct {
	// hold the actual values of the cache
	values map[K]V
	// whether or not this cache should be written to storage
	// if memory == true, Save() does nothing
	memory bool
	// storage is a WriteSeeker to use while saving to persistent storage, like
	// a file
	storage io.ReadWriteSeeker
	// ensure thread-safety
	mx *sync.RWMutex
}

// NewMemory creates an empty in-memory cache
func NewMemory[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		values:  make(map[K]V),
		memory:  true,
		storage: nil,
		mx:      &sync.RWMutex{},
	}
}

// LoadMemory creates an empty in-memory cache and populates it with the
// values decoded from the given Reader
func LoadMemory[K comparable, V any](src io.Reader) (*Cache[K, V], error) {
	c := NewMemory[K, V]()
	c.memory = false
	return c, c.LoadFrom(src)
}

// NewStored creates an empty persistent cache with the given storage
func NewStored[K comparable, V any](storage io.ReadWriteSeeker) *Cache[K, V] {
	return &Cache[K, V]{
		values:  make(map[K]V),
		memory:  false,
		storage: storage,
		mx:      &sync.RWMutex{},
	}
}

// LoadStored creates an empty persistent cache and populates it with the
// values decoded from the underlying storage
func LoadStored[K comparable, V any](storage io.ReadWriteSeeker) (*Cache[K, V], error) {
	c := NewStored[K, V](storage)
	c.memory = false
	return c, c.Load()
}

// Load loads a persistent cache from it's underlying storage. For in-memory
// caches this does nothing.
func (c *Cache[K, V]) Load() error {
	if c.memory {
		return nil
	}
	return c.LoadFrom(c.storage)
}

// LoadFrom loads a cache from the given reader
func (c *Cache[K, V]) LoadFrom(src io.Reader) error {
	c.mx.RLock()
	defer c.mx.RUnlock()
	dec := json.NewDecoder(src)
	return dec.Decode(&c.values)
}

// Save saves a persistent cache's current state into it's underlying storage.
// For in-memory caches this does nothing.
func (c *Cache[K, V]) Save() error {
	if c.memory {
		return nil
	}
	c.mx.Lock()
	defer c.mx.Unlock()
	_, err := c.storage.Seek(0, 0)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(c.storage)
	return enc.Encode(c.values)
}

// Get retrieves a value from the cache, outputs work the same way as in
// accessing a map
func (c *Cache[K, V]) Get(key K) (v V, ok bool) {
	v, ok = c.values[key]
	return
}

// GetOrDefault retrieves a value from the cache, returning the given default
// if it is not found
func (c *Cache[K, V]) GetOrDefault(key K, def V) V {
	v, ok := c.Get(key)
	if !ok {
		return def
	}
	return v
}

// GetOrLoad retrieves a value from the cache, reloading said cache if it isn't
// found. This function will give up after 1 reload.
func (c *Cache[K, V]) GetOrLoad(key K) (v V, ok bool, err error) {
	v, ok = c.Get(key)
	if ok {
		return
	}
	err = c.Load()
	v, ok = c.Get(key)
	return
}

// GetLoadOrDefault retireves a value from the cache, reload said cache if it
// isn't found. This will return the default if the value is still not found
// after 1 reload.
func (c *Cache[K, V]) GetLoadOrDefault(key K, def V) (v V, ok bool, err error) {
	v, ok, err = c.GetOrLoad(key)
	if ok || err != nil {
		return
	}
	return def, false, nil
}

// Has checks if the cache contains a certain key
func (c *Cache[K, V]) Has(key K) (ok bool) {
	_, ok = c.Get(key)
	return
}

// Set sets the value in the cache
func (c *Cache[K, V]) Set(key K, val V) {
	c.values[key] = val
}

// SetAndSave sets a value in the cache and then saves it. For in-memory caches
// this only sets without saving. The save is done in a separate goroutine and
// errors are ignored.
func (c *Cache[K, V]) SetAndSave(key K, val V) {
	c.Set(key, val)
	go c.Save()
}
