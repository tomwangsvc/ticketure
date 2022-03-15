package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test_NewCache_WriteAndRead(t *testing.T) {
	cacheProvider, shutdownFunc := NewCache(time.Second)
	defer shutdownFunc()

	// test read 'key' not found
	f, value, err := cacheProvider.Read("key")
	if err != nil {
		t.Errorf("unexpected error to read: %s", err)
	} else if f {
		t.Errorf("Expected '%t', got '%t'", false, f)
	} else if len(value) != 0 {
		t.Errorf("Expected '%d', got '%d'", 0, len(value))
	}

	if err := cacheProvider.Write("key", []byte("value")); err != nil {
		t.Errorf("unexpected error to write: %s", err)
	}

	// test read 'key' found
	f, value, err = cacheProvider.Read("key")
	if err != nil {
		t.Errorf("unexpected error to read: %s", err)
	} else if !f {
		t.Errorf("Expected '%t', got '%t'", true, f)
	} else if string(value) != "value" {
		t.Errorf("Expected '%s', got '%s'", "value", string(value))
	}
}

func Test_Concurrency_NewCache_WriteAndRead(t *testing.T) {
	cacheProvider, shutdownFunc := NewCache(time.Hour) // set this not to expire to test concurrency
	defer shutdownFunc()

	numOfRun := 1000
	wg := sync.WaitGroup{}
	wg.Add(numOfRun)
	for i := 0; i < numOfRun; i++ {
		go func(i int) {
			wg.Done()
			if err := cacheProvider.Write(fmt.Sprintf("%d", i), []byte(fmt.Sprintf("value_%d", i))); err != nil {
				fmt.Printf("failed to write %s \n", err)
				return
			}
		}(i)
	}

	wg.Wait()

	for i := 0; i < numOfRun; i++ {
		if f, v, err := cacheProvider.Read(fmt.Sprintf("%d", i)); err != nil {
			t.Errorf("unexpected error to read: %s", err)
		} else if !f {
			t.Errorf("Expected '%t', got '%t'", true, f)
		} else if string(v) != fmt.Sprintf("value_%d", i) {
			t.Errorf("Expected '%s', got '%s'", string(v), fmt.Sprintf("value_%d", i))
		}
	}
}

func Test_NewCache_shutdownFunc(t *testing.T) {
	cacheProvider, shutdownFunc := NewCache(time.Second)

	cacheProvider.Write("key", []byte("value"))
	time.Sleep(2 * time.Second) // greater than cacheFor and value in cache will be cleaned
	f, _, _ := cacheProvider.Read("key")
	if f {
		t.Errorf("Expected key in cache to be removed")
	}

	shutdownFunc()
	cacheProvider.Write("new_key", []byte("value"))
	time.Sleep(2 * time.Second)  // greater than cacheFor and value in cache will not be cleaned as func shutdownFunc has been run before
	f, _, _ = cacheProvider.Read("new_key")
	if !f {
		t.Errorf("Expected key in cache not to be removed after running shutdown func")
	}
}

func Benchmark_NewCache_Write(b *testing.B) {
	cacheProvider, shutdownFunc := NewCache(time.Hour) // set this not to expire for benchmark test
	defer shutdownFunc()

	for i := 0; i < b.N; i++ {
		cacheProvider.Write("key", []byte("value"))
		cacheProvider.Read("key")
	}
}
