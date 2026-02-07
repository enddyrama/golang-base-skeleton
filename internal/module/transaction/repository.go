package transaction

import (
	"base-skeleton/internal/module/product"
	transactiondetail "base-skeleton/internal/module/transaction_detail"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	CreateTransaction(items []CheckoutItem) (*Transaction, error)
	GetReport(start, end time.Time) (ReportResponse, error)
}

type repository struct {
	db          *sql.DB
	productRepo product.Repository
	detailRepo  transactiondetail.Repository
}

func NewRepository(
	db *sql.DB,
	productRepo product.Repository,
	detailRepo transactiondetail.Repository,
) Repository {
	return &repository{
		db:          db,
		productRepo: productRepo,
		detailRepo:  detailRepo,
	}
}

func (r *repository) CreateTransaction(
	items []CheckoutItem,
) (*Transaction, error) {

	if len(items) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var (
		totalAmount int64
		details     []transactiondetail.TransactionDetailResponse
	)

	// 1️⃣ lock products, calculate, update stock
	for _, item := range items {

		p, err := r.productRepo.FindByIDForUpdateTx(
			tx,
			item.ProductID,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if p.Stock < item.Quantity {
			return nil, fmt.Errorf(
				"insufficient stock for product %d",
				item.ProductID,
			)
		}

		subtotal := p.Price * item.Quantity
		totalAmount += subtotal

		if err := r.productRepo.DecreaseStockTx(
			tx,
			item.ProductID,
			item.Quantity,
		); err != nil {
			return nil, err
		}

		details = append(details, transactiondetail.TransactionDetailResponse{
			ProductID:   item.ProductID,
			ProductName: p.Name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// 2️⃣ create transaction master
	var transactionID int64
	err = tx.QueryRow(`
		INSERT INTO transactions (total_amount)
		VALUES ($1)
		RETURNING id
	`, totalAmount).Scan(&transactionID)
	if err != nil {
		fmt.Println("ERROR scanning transaction id:", err)
		return nil, err
	}
	fmt.Println("transactionID:", transactionID)

	// 3️⃣ insert details
	for i := range details {
		if err := r.detailRepo.InsertTx(
			tx,
			transactiondetail.TransactionDetail{
				TransactionID: transactionID,
				ProductID:     details[i].ProductID,
				Quantity:      details[i].Quantity,
				Subtotal:      details[i].Subtotal,
			},
		); err != nil {
			return nil, err
		}

		details[i].TransactionID = transactionID
	}

	// 4️⃣ commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// 5️⃣ return response
	return &Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (r *repository) GetReport(start, end time.Time) (ReportResponse, error) {
	var resp ReportResponse

	// ======================
	// Total revenue from transactions
	// ======================
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount),0)
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2
	`, start, end).Scan(&resp.TotalRevenue)
	if err != nil {
		log.Println("resp total revenue", resp)
		return resp, err
	}

	// ======================
	// Total transaction (count of unique transaction_id in transaction_details)
	// ======================
	err = r.db.QueryRow(`
		SELECT COUNT(DISTINCT transaction_id)
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= $1 AND t.created_at <= $2
	`, start, end).Scan(&resp.TotalTransaction)
	if err != nil {
		log.Println("resp total transaction", resp)
		return resp, err
	}

	// ======================
	// Best-selling product
	// ======================
	err = r.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) as sold
		FROM transaction_details td
		JOIN transactions t ON t.id = td.transaction_id
		JOIN products p ON p.id = td.product_id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.name
		ORDER BY sold DESC
		LIMIT 1
	`, start, end).Scan(&resp.BestSellingProduct.Name, &resp.BestSellingProduct.Sold)

	// If no transactions exist in the range, return empty
	if err == sql.ErrNoRows {
		log.Println("no rows", err)
		resp.BestSellingProduct.Name = ""
		resp.BestSellingProduct.Sold = 0
		err = nil
	} else if err != nil {
		log.Println("resp total transaction", resp)
		return resp, err
	}

	return resp, nil
}
