package pokecache

import (
	"testing"
	"time"
)

func TestCreateCache(t *testing.T) {
	cache := NewCache(5)
	if cache.cache == nil {
		t.Error("cache is nil")
	}
}
func TestAddGetCache(t *testing.T) {
	cache := NewCache(5)
	cache.Add("key1", []byte("val1"))
	actual, ok := cache.Get("key1")
	if !ok {
		t.Error("key 1 not found")
	}
	if string(actual) != "val1" {
		t.Error("value doesn't match")
	}
}

func TestReap(t *testing.T) {
	cache := NewCache(5)

	keyOne := "key1"
	cache.Add(keyOne, []byte("val1"))
	time.Sleep(time.Second*5 + time.Second)

	_, ok := cache.Get(keyOne)
	if ok {
		t.Errorf("not reaped")
	}
}
