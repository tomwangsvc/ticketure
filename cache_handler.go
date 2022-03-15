package main

import (
	"bytes"
	"encoding/gob"
	"net/http"
)

// CachedHandler wraps a handler to cache responses for short periods of time
func CachedHandler(c CacheProvider, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + r.URL.RequestURI()

		found, cachedBytes, err := c.Read(key)
		if err != nil {
			panic(err)
		}

		if found {
			cachedResponse := deserialiseBufferedResponseWriter(cachedBytes)
			cachedResponse.WriteToResponseWriter(w)
			return
		}

		buffer := &bufferedResponseWriter{buffer: &bytes.Buffer{}}
		h(buffer, r)

		err = c.Write(key, buffer.Serialise())
		if err != nil {
			panic(err)
		}

		buffer.WriteToResponseWriter(w)
		return
	}
}

// bufferedResponseWriter packages up an outgoing response for saving and rehydrating
type bufferedResponseWriter struct {
	buffer     *bytes.Buffer
	headers    http.Header
	statusCode int
}

func (b *bufferedResponseWriter) Header() http.Header            { return b.headers }
func (b *bufferedResponseWriter) Write(data []byte) (int, error) { return b.buffer.Write(data) }
func (b *bufferedResponseWriter) WriteHeader(code int)           { b.statusCode = code }

func deserialiseBufferedResponseWriter(buf []byte) *bufferedResponseWriter {
	out := &bufferedResponseWriter{}
	dec := gob.NewDecoder(bytes.NewReader(buf))
	dec.Decode(&out.headers)
	dec.Decode(&out.statusCode)
	b := []byte{}
	dec.Decode(&b)
	out.buffer = bytes.NewBuffer(b)
	return out
}

func (b *bufferedResponseWriter) Serialise() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.headers)
	enc.Encode(b.statusCode)
	enc.Encode(b.buffer.Bytes())
	return buf.Bytes()
}

func (b *bufferedResponseWriter) WriteToResponseWriter(w http.ResponseWriter) error {
	headers := w.Header()
	for k, v := range b.headers {
		headers[k] = v
	}
	w.WriteHeader(b.statusCode)
	_, err := b.buffer.WriteTo(w)
	return err
}
