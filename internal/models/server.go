package models

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, mux *chi.Mux) error {
	fmt.Println("Сервер запускается")
	return http.ListenAndServe(port, mux)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
