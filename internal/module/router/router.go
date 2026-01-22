package router

import (
	"net/http"

	"base-skeleton/internal/module/categories"
	"base-skeleton/internal/module/health"
	"base-skeleton/internal/shared/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// Health
	healthService := health.NewService()
	healthHandler := health.NewHandler(healthService)
	health.Register(mux, healthHandler)

	// Categories
	categoryService := categories.NewService()
	categoryHandler := categories.NewHandler(categoryService)
	categories.Register(mux, categoryHandler)

	return middleware.Recover(mux)
}
