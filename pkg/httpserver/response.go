package httpserver

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

var _ http.ResponseWriter = (*Response)(nil)
var _ http.Flusher = (*Response)(nil)
var _ http.Hijacker = (*Response)(nil)

// CustomWriter
type Response struct {
	rw         http.ResponseWriter
	StatusCode int
}

func (r *Response) Header() http.Header {
	return r.rw.Header()
}

func (r *Response) Write(b []byte) (int, error) {
	return r.rw.Write(b)
}

func (r *Response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.rw.WriteHeader(statusCode)
}

func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// Проверяем, реализует ли исходный ResponseWriter интерфейс Hijacker
	if hijacker, ok := r.rw.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
}

// Flush реализует http.Flusher если исходный ResponseWriter его поддерживает
func (r *Response) Flush() {
	if flusher, ok := r.rw.(http.Flusher); ok {
		flusher.Flush()
	}
}
