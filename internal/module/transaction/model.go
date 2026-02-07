package transaction

import (
	transactiondetail "base-skeleton/internal/module/transaction_detail"
	"time"
)

type Transaction struct {
	ID          int64                                         `json:"id"`
	TotalAmount int64                                         `json:"total_amount"`
	CreatedAt   time.Time                                     `json:"created_at"`
	Details     []transactiondetail.TransactionDetailResponse `json:"details"`
}

type CheckoutItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// /REPORT
type ReportRequest struct {
	StartDate string `json:"start_date"` // ISO8601 "2026-02-01T00:00:00"
	EndDate   string `json:"end_date"`   // ISO8601 "2026-02-07T23:59:59"
}

type BestSellingProduct struct {
	Name string `json:"name"`
	Sold int64  `json:"sold"`
}

type ReportResponse struct {
	TotalRevenue       int64              `json:"total_revenue"`
	TotalTransaction   int64              `json:"total_transaction"`
	BestSellingProduct BestSellingProduct `json:"best_selling_product"`
}
