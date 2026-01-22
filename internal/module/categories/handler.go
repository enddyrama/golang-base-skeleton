package categories

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	appErr "base-skeleton/internal/shared/errors"
	"base-skeleton/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func getIDFromPath(r *http.Request) (int64, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]

	return strconv.ParseInt(idStr, 10, 64)
}

func (h *Handler) categories(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {

	case http.MethodGet:
		data, err := h.service.GetAll()
		if err != nil {
			return err
		}
		return response.JSON(w, http.StatusOK, "success", data)

	case http.MethodPost:
		var req Category
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		res, appErr := h.service.Create(req)
		if appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusCreated, "category created", res)

	default:
		return nil
	}
}

func (h *Handler) categoryByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromPath(r)
	if err != nil {
		return appErr.BadRequest("invalid category id")
	}

	switch r.Method {

	case http.MethodGet:
		res, appErr := h.service.GetByID(id)
		if appErr != nil {
			return appErr
		}
		return response.JSON(w, http.StatusOK, "success", res)

	case http.MethodPut:
		var req Category
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		res, appErr := h.service.Update(id, req)
		if appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusOK, "category updated", res)

	case http.MethodDelete:
		if appErr := h.service.Delete(id); appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusOK, "category deleted", nil)

	default:
		return nil
	}
}
