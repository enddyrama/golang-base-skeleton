package transaction

import (
	"base-skeleton/internal/shared/response"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
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

func (h *Handler) HandleReportToday(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	tz := r.URL.Query().Get("timezone")
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	startLocal := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endLocal := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, loc)

	startStrUTC := startLocal.UTC().Format(time.RFC3339)
	endStrUTC := endLocal.UTC().Format(time.RFC3339)

	log.Println(startStrUTC)
	log.Println(endStrUTC)
	report, appErr := h.service.GetReport(startStrUTC, endStrUTC)
	if appErr != nil {
		return response.JSON(w, appErr.Code, appErr.Message, nil)
	}

	return response.JSON(w, http.StatusOK, "report fetched successfully", report)
}

func (h *Handler) HandleReport(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	// Read query params
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")
	tz := r.URL.Query().Get("timezone")
	if tz == "" {
		tz = "UTC"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}

	if startStr == "" || endStr == "" {
		return response.JSON(w, http.StatusBadRequest, "start_date and end_date are required", nil)
	}

	// Remove Z if user accidentally sends UTC timestamp
	startStr = strings.TrimSuffix(startStr, "Z")
	endStr = strings.TrimSuffix(endStr, "Z")

	// Parse input as **user local time**
	startLocal, err := time.ParseInLocation("2006-01-02T15:04:05", startStr, loc)
	if err != nil {
		return response.JSON(w, http.StatusBadRequest, "invalid start_date format, must be YYYY-MM-DDTHH:MM:SS", nil)
	}
	endLocal, err := time.ParseInLocation("2006-01-02T15:04:05", endStr, loc)
	if err != nil {
		return response.JSON(w, http.StatusBadRequest, "invalid end_date format, must be YYYY-MM-DDTHH:MM:SS", nil)
	}

	if endLocal.Before(startLocal) {
		return response.JSON(w, http.StatusBadRequest, "end_date cannot be before start_date", nil)
	}

	// Convert to UTC for database query
	startUTC := startLocal.UTC().Format(time.RFC3339)
	endUTC := endLocal.UTC().Format(time.RFC3339)

	log.Println("User timezone:", tz)
	log.Println("StartLocal:", startLocal, "StartUTC:", startUTC)
	log.Println("EndLocal:", endLocal, "EndUTC:", endUTC)

	report, appErr := h.service.GetReport(startUTC, endUTC)
	if appErr != nil {
		return response.JSON(w, appErr.Code, appErr.Message, nil)
	}

	return response.JSON(w, http.StatusOK, "report fetched successfully", report)
}
