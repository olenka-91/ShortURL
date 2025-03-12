package handlers

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHandler_InitRoutes(t *testing.T) {
	tests := []struct {
		name string
		h    *Handler
		want *chi.Mux
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.InitRoutes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.InitRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}
