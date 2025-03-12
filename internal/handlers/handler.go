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
	r.Post(`/`, logger.WithLogging(compressMiddleware.GzipMiddleware(h.PostShortURL)))
	r.Post(`/api/shorten`, logger.WithLogging(compressMiddleware.GzipMiddleware(h.PostShortURLJSON)))
	r.Get(`/{id}`, logger.WithLogging(compressMiddleware.GzipMiddleware(h.GetUnShortURL)))

	return r

}
