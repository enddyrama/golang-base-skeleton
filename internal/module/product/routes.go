package product

import (
	"base-skeleton/internal/shared/middleware"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func Register(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/v1/products", middleware.Wrap(h.product))
	mux.HandleFunc("/api/v1/products/", middleware.Wrap(h.productByID))
}
