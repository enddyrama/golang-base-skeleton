package router

import (
	"database/sql"
	"net/http"

	"base-skeleton/internal/module/category"
	"base-skeleton/internal/module/health"
	"base-skeleton/internal/module/product"
	"base-skeleton/internal/shared/middleware"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	healthService := health.NewService()
	healthHandler := health.NewHandler(healthService)
	health.Register(mux, healthHandler)

	categoryService := category.NewService(db)
	categoryHandler := category.NewHandler(categoryService)
	category.Register(mux, categoryHandler)

	productService := product.NewService(db)
	productHandler := product.NewHandler(productService)
	product.Register(mux, productHandler)

	return middleware.Recover(mux)
}
