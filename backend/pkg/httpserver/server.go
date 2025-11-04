package httpserver

import "net/http"

type Server struct {
	mv  []func(next http.Handler) http.Handler
	mux *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		mv:  make([]func(next http.Handler) http.Handler, 0),
		mux: http.NewServeMux(),
	}
}

func (s *Server) AddMiddleware(m func(next http.Handler) http.Handler) {
	s.mv = append(s.mv, m)
}

func (s *Server) AddHandler(pat string, h http.Handler, mvs ...func(next http.Handler) http.Handler) {
	for _, m := range s.mv {
		h = m(h)
	}
	for _, m := range mvs {
		h = m(h)
	}
	s.mux.Handle(pat, h)
}

func (s *Server) Run(addr string) error {
	server := &http.Server{Addr: addr, Handler: s.mux}
	return server.ListenAndServe()
}
