//Copyright (c) 2017 Phil

package apollo

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

type namespaceCache struct {
	lock   sync.RWMutex
	caches map[string]*cache
}

func newNamespaceCache() *namespaceCache {
	return &namespaceCache{
		caches: map[string]*cache{},
	}
}

func (n *namespaceCache) mustGetCache(namespace string) *cache {
	n.lock.RLock()
	if ret, ok := n.caches[namespace]; ok {
		n.lock.RUnlock()
		return ret
	}
	n.lock.RUnlock()

	n.lock.Lock()
	defer n.lock.Unlock()

	cache := newCache()
	n.caches[namespace] = cache
	return cache
}

func (n *namespaceCache) drain() {
	for namespace := range n.caches {
		delete(n.caches, namespace)
	}
}

func (n *namespaceCache) dump(name string) error {

	var dumps = map[string]map[string][]byte{}

	for namespace, cache := range n.caches {
		dumps[namespace] = cache.dump()
	}

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewEncoder(f).Encode(&dumps)
}

func (n *namespaceCache) load(name string) error {
	n.drain()

	f, err := os.OpenFile(name, os.O_RDONLY, 0755)
	if err != nil {
		return errors.New("open cache file error")
	}
	defer f.Close()

	var dumps = map[string]map[string][]byte{}

	if err := gob.NewDecoder(f).Decode(&dumps); err != nil {
		return errors.New("cache file decoder error")
	}

	for namespace, kv := range dumps {
		cache := n.mustGetCache(namespace)
		for k, v := range kv {
			cache.set(k, v)
		}
	}

	return nil
}

type cache struct {
	kv sync.Map
}

func newCache() *cache {
	return &cache{
		kv: sync.Map{},
	}
}

func (c *cache) set(key string, val []byte) {
	c.kv.Store(key, val)
}

func (c *cache) get(key string) ([]byte, bool) {
	if val, ok := c.kv.Load(key); ok {
		if ret, ok := val.([]byte); ok {
			return ret, true
		}
	}
	return nil, false
}

func (c *cache) delete(key string) {
	c.kv.Delete(key)
}

func (c *cache) dump() map[string][]byte {
	var ret = map[string][]byte{}
	c.kv.Range(func(key, val interface{}) bool {
		k, _ := key.(string)
		v, _ := val.([]byte)
		ret[k] = v

		return true
	})
	return ret
}
