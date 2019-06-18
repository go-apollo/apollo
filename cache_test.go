//Copyright (c) 2017 Phil
package apollo

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

func (s *CacheTestSuite) TestCache() {
	log.Println("test cache")
	log.Printf("s: %+v", s)
	cache := newCache()

	cache.set("key", []byte("val"))
	val, ok := cache.get("key")

	s.True(ok)
	s.Equal("val", string(val))

	cache.set("key", []byte("val2"))
	val1, ok1 := cache.get("key")

	s.True(ok1)
	s.Equal("val2", string(val1))

	kv := cache.dump()

	s.Equal(1, len(kv))
	s.Equal("val2", string(kv["key"]))

	cache.delete("key")
	_, ok2 := cache.get("key")
	s.False(ok2)
}

func (s *CacheTestSuite) TestCacheDump() {
	var caches = newNamespaceCache()
	defer caches.drain()

	caches.mustGetCache("namespace").set("key", []byte("val"))

	f, err := ioutil.TempFile(".", "apollo")
	s.NoError(err)
	f.Close()
	defer os.Remove(f.Name())

	s.NoError(caches.dump(f.Name()))

	var restore = newNamespaceCache()
	defer restore.drain()

	s.NoError(restore.load(f.Name()))

	val, _ := restore.mustGetCache("namespace").get("key")

	s.Equal("val", string(val))

	s.Error(restore.load("null"))

	s.Error(restore.load("./testdata/app.yml"))
}
