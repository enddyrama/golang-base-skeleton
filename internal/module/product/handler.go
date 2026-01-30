package product

import (
	"encoding/json"
	"math"
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

func (h *Handler) product(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {

	case http.MethodGet:
		// ========================
		// Pagination
		// ========================
		page := 1
		size := 10

		if v := r.URL.Query().Get("page"); v != "" {
			page, _ = strconv.Atoi(v)
		}
		if v := r.URL.Query().Get("size"); v != "" {
			size, _ = strconv.Atoi(v)
		}

		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 10
		}
		if size > 100 {
			size = 100
		}

		offset := (page - 1) * size

		// ========================
		// Search & Sort
		// ========================
		search := r.URL.Query().Get("search")
		sort := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")

		// ========================
		// Service call
		// ========================
		data, total, err := h.service.GetAll(
			size,
			offset,
			search,
			sort,
			order,
		)
		if err != nil {
			return err
		}

		totalPage := int(math.Ceil(float64(total) / float64(size)))

		result := response.ListResult[ProductResponse]{
			Data: data,
			Pagination: response.Pagination{
				Page:      page,
				Size:      size,
				Total:     total,
				TotalPage: totalPage,
			},
		}

		return response.JSON(w, http.StatusOK, "success", result)

	case http.MethodPost:
		var req Product
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		res, appErr := h.service.Create(req)
		if appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusCreated, "product created", res)

	default:
		return nil
	}
}

func (h *Handler) productByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromPath(r)
	if err != nil {
		return appErr.BadRequest("invalid product id")
	}

	switch r.Method {

	case http.MethodGet:
		res, appErr := h.service.GetByID(id)
		if appErr != nil {
			return appErr
		}
		return response.JSON(w, http.StatusOK, "success", res)

	case http.MethodPut:
		var req Product
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		res, appErr := h.service.Update(id, req)
		if appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusOK, "product updated", res)

	case http.MethodDelete:
		if appErr := h.service.Delete(id); appErr != nil {
			return appErr
		}

		return response.JSON(w, http.StatusOK, "product deleted", nil)

	default:
		return nil
	}
}
