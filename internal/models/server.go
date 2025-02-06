package models

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, mux *http.ServeMux) error {
	fmt.Println("Сервер запускается")
	return http.ListenAndServe(port, mux)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
