package transaction

import (
	"base-skeleton/internal/module/product"
	transactiondetail "base-skeleton/internal/module/transaction_detail"
	"base-skeleton/internal/shared/errors"
	appErr "base-skeleton/internal/shared/errors"
	"log"
	"time"
)

type Service struct {
	transactionRepo       Repository
	productRepo           product.Repository
	transactionDetailRepo transactiondetail.Repository
}

func NewService(transactionRepo Repository, productRepo product.Repository, transactionDetailRepo transactiondetail.Repository) *Service {
	return &Service{
		transactionRepo:       transactionRepo,
		productRepo:           productRepo,
		transactionDetailRepo: transactionDetailRepo,
	}
}

func (s *Service) Checkout(
	req CheckoutRequest,
) (*Transaction, *appErr.AppError) {

	// validation
	if len(req.Items) == 0 {
		return nil, appErr.BadRequest("items cannot be empty")
	}

	for _, item := range req.Items {
		if item.ProductID <= 0 {
			return nil, appErr.BadRequest("invalid product_id")
		}
		if item.Quantity <= 0 {
			return nil, appErr.BadRequest("quantity must be greater than 0")
		}
		_, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			log.Println(err) // prints with date and time
			if err == appErr.ErrNotFound {
				return nil, appErr.Custom(404, "Product not found: id -%d", item.ProductID)
			}
			return nil, appErr.Internal("Failed to query product:" + err.Error())
		}
	}

	// business logic (repo handles transaction & stock)
	res, err := s.transactionRepo.CreateTransaction(req.Items)
	if err != nil {

		// map known domain errors
		switch err.Error() {

		case "insufficient stock":
			return nil, appErr.BadRequest("insufficient stock")

		default:
			return nil, appErr.Internal("failed to create transaction")
		}
	}

	return res, nil
}

func (s *Service) GetReport(startStr, endStr string) (*ReportResponse, *errors.AppError) {
	// ======================
	// Validate dates
	// ======================
	if startStr == "" || endStr == "" {
		return nil, errors.BadRequest("start_date and end_date are required")
	}

	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return nil, errors.BadRequest("invalid start_date format, must be RFC3339")
	}

	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return nil, errors.BadRequest("invalid end_date format, must be RFC3339")
	}

	if endDate.Before(startDate) {
		return nil, errors.BadRequest("end_date cannot be before start_date")
	}

	// ======================
	// Call repository (already returns ReportResponse)
	// ======================
	report, repoErr := s.transactionRepo.GetReport(startDate, endDate)
	if repoErr != nil {
		return nil, errors.Internal("failed to fetch report data")
	}

	return &report, nil
}
