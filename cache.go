package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// CacheProvider describes a cache, in a real system this might be implemented using Redis
type CacheProvider interface {
	Read(key string) (found bool, value []byte, err error)
	Write(string, []byte) error
}

type cacheItem struct {
	value   []byte
	created time.Time
}

type cache struct {
	// TODO: maybe implement sync.Mutex if its benchmark test is better than sync.Map
	values   sync.Map
	cacheFor time.Duration
}

// NewCache creates a cache that keeps items for the duration cacheFor
func NewCache(cacheFor time.Duration) (cacheProvider CacheProvider, shutdownFunc func()) {
	c := &cache{
		values:   sync.Map{},
		cacheFor: cacheFor,
	}

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(cacheFor)
		for {
			select {
			case <-ticker.C:
				c.values.Range(func(k, v interface{}) bool {
					value, ok := v.(cacheItem)
					if !ok {
						log.Printf("cacheItem %T is wrong type", v)
						return false
					}

					if time.Now().After(value.created.Add(cacheFor)) {
						c.values.Delete(k)
					}

					return true
				})

			case <-done:
				return
			}
		}
	}()

	// TODO: return a working shutdown function that stops the go-routine
	return c, func() {
		done <- struct{}{}
	}
}

func (c *cache) Read(key string) (found bool, value []byte, err error) {
	time.Sleep(5 * time.Millisecond) // emulate some cache layer delay

	v, found := c.values.Load(key)
	if !found{
		return
	}
	cacheItem, ok := v.(cacheItem)
	if !ok {
		err = fmt.Errorf("value %T is wrong type", v)
		return
	}
	return found, cacheItem.value, nil
}

func (c *cache) Write(k string, value []byte) error {
	time.Sleep(5 * time.Millisecond) // emulate some cache layer delay

	c.values.Store(k, cacheItem{created: time.Now(), value: value})
	return nil
}
