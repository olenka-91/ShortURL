package handlers

import (
	"net/http"

	"github.com/olenka-91/shorturl/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{services: serv}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, h.PostShortURL)
	mux.HandleFunc(`/{id}`, h.GetUnShortURL)

	return mux

}
