package transaction

import (
	"base-skeleton/internal/shared/response"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleCheckout(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		var req CheckoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		res, appErr := h.service.Checkout(req)
		if appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusCreated, "product created", res)
	default:
		return nil
	}
}

func (h *Handler) HandleReport(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	// Read query params
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	report, appErr := h.service.GetReport(startDate, endDate)
	if appErr != nil {
		return response.JSON(w, appErr.Code, appErr.Message, nil)
	}

	return response.JSON(w, http.StatusOK, "report fetched successfully", report)
}
