package category

import (
	"base-skeleton/internal/shared/middleware"
	"net/http"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func Register(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/api/v1/category", middleware.Wrap(h.categories))
	mux.HandleFunc("/api/v1/category/", middleware.Wrap(h.categoryByID))
}
