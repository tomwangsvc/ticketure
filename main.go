package main

import (
	"math/rand"
	"net/http"
	"time"
)

// this file should not need to be edited

func main() {
	cache, shutdownFunc := NewCache(30 * time.Second)
	defer shutdownFunc()

	h := http.HandlerFunc(doExpensiveWork)
	h = CachedHandler(cache, h)

	http.ListenAndServe(":9000", h)
}

func doExpensiveWork(w http.ResponseWriter, r *http.Request) {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	out := make([]byte, 20)
	for i := range out {
		out[i] = letters[rand.Intn(len(letters))]
	}

	time.Sleep(4 * time.Second) // this makes it expensive
	w.WriteHeader(200)
	w.Write(out)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
