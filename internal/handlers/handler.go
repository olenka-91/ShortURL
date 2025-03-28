package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/olenka-91/shorturl/internal/compressMiddleware"
	"github.com/olenka-91/shorturl/internal/logger"
	"github.com/olenka-91/shorturl/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{services: serv}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compressMiddleware.GzipMiddleware)

	r.Post(`/`, (h.PostShortURL))
	r.Post(`/api/shorten`, (h.PostShortURLJSON))
	r.Get(`/{id}`, (h.GetUnShortURL))
	r.Get(`/ping`, (h.GetDBPing))

	return r

}
