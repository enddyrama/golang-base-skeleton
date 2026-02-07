package router

import (
	"database/sql"
	"net/http"

	"base-skeleton/internal/module/category"
	"base-skeleton/internal/module/health"
	"base-skeleton/internal/module/product"
	"base-skeleton/internal/module/transaction"
	transactiondetail "base-skeleton/internal/module/transaction_detail"
	"base-skeleton/internal/shared/middleware"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// =========================
	// Repository
	// =========================
	categoryRepo := category.NewRepository(db)
	productRepo := product.NewRepository(db)
	transactionDetailRepo := transactiondetail.NewRepository(db)
	transactionRepo := transaction.NewRepository(db, productRepo, transactionDetailRepo)
	// =========================
	// Health
	// =========================
	healthService := health.NewService()
	healthHandler := health.NewHandler(healthService)
	health.Register(mux, healthHandler)

	// =========================
	// Category
	// =========================
	categoryService := category.NewService(categoryRepo)
	categoryHandler := category.NewHandler(categoryService)
	category.Register(mux, categoryHandler)

	// =========================
	// Product
	// =========================
	productService := product.NewService(productRepo, categoryRepo)
	productHandler := product.NewHandler(productService)
	product.Register(mux, productHandler)

	// =========================
	// Transaction
	// =========================
	transactionService := transaction.NewService(transactionRepo, productRepo, transactionDetailRepo)
	transactionHandler := transaction.NewHandler(transactionService)
	transaction.Register(mux, transactionHandler)

	return middleware.Recover(mux)
}
