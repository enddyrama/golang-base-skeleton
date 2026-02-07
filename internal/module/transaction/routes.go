package transaction

import (
	"base-skeleton/internal/shared/middleware"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func Register(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/v1/checkout", middleware.Wrap(h.HandleCheckout))
	mux.HandleFunc("/api/v1/report", middleware.Wrap(h.HandleReport))
}
