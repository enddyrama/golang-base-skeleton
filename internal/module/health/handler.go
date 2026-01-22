package health

import (
	"net/http"

	"base-skeleton/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) error {
	return response.JSON(
		w,
		http.StatusOK,
		"success",
		h.service.Check(),
	)
}
