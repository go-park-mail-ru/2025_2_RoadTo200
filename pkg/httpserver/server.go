package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/internal/config"
	"github.com/go-park-mail-ru/2025_2_RoadTo200/backend/pkg/logger"
)

type Server struct {
	mv  []func(next http.Handler) http.Handler
	mux *http.ServeMux
	smv SpecificMV
}

type SpecificMV struct {
	init func(next http.Handler) http.Handler
	auth func(next http.Handler) http.Handler
	mtrc func(next http.Handler) http.Handler
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		mv:  make([]func(next http.Handler) http.Handler, 0),
		mux: http.NewServeMux(),
		smv: SpecificMV{
			init: InitMiddleware(logger.NewFactory(&cfg.Logger)),
			auth: func(next http.Handler) http.Handler { return next },
			mtrc: func(next http.Handler) http.Handler { return next },
		},
	}
}

func (s *Server) AddMiddleware(m func(next http.Handler) http.Handler) {
	s.mv = append(s.mv, m)
}

func (s *Server) SetAuthMiddleware(m func(next http.Handler) http.Handler) {
	s.smv.auth = m
}

func (s *Server) SetMetricMiddleware(m func(next http.Handler) http.Handler) {
	s.smv.mtrc = m
}

func (s *Server) AddHandler(pat string, h http.Handler, mvs ...func(next http.Handler) http.Handler) {
	for _, m := range s.mv {
		h = m(h)
	}
	for _, m := range mvs {
		h = m(h)
	}
	h = s.smv.auth(h)
	h = s.smv.mtrc(h)
	h = s.smv.init(h)
	s.mux.Handle(pat, h)
}

func (s *Server) GET(pat string, h func(w http.ResponseWriter, r *http.Request), mvs ...func(next http.Handler) http.Handler) {
	s.AddHandler(fmt.Sprintf("GET %s", pat), http.HandlerFunc(h), mvs...)
}

func (s *Server) POST(pat string, h func(w http.ResponseWriter, r *http.Request), mvs ...func(next http.Handler) http.Handler) {
	s.AddHandler(fmt.Sprintf("POST %s", pat), http.HandlerFunc(h), mvs...)
}

func (s *Server) PUT(pat string, h func(w http.ResponseWriter, r *http.Request), mvs ...func(next http.Handler) http.Handler) {
	s.AddHandler(fmt.Sprintf("PUT %s", pat), http.HandlerFunc(h), mvs...)
}

func (s *Server) DELETE(pat string, h func(w http.ResponseWriter, r *http.Request), mvs ...func(next http.Handler) http.Handler) {
	s.AddHandler(fmt.Sprintf("DELETE %s", pat), http.HandlerFunc(h), mvs...)
}

func (s *Server) PATCH(pat string, h func(w http.ResponseWriter, r *http.Request), mvs ...func(next http.Handler) http.Handler) {
	s.AddHandler(fmt.Sprintf("PATCH %s", pat), http.HandlerFunc(h), mvs...)
}

func (s *Server) Run(addr string) error {
	server := &http.Server{Addr: addr, Handler: s.mux}
	return server.ListenAndServe()
}
