package health

import (
	"base-skeleton/internal/shared/middleware"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func Register(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/health", middleware.Wrap(h.Check))
}
