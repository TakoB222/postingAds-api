package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer http.Server
}

type Config struct {
	Host           string
	Port           string
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func NewServer(cfg Config, handler http.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:           cfg.Host + ":" + cfg.Port,
			Handler:        handler,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
