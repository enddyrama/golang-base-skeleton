package transactiondetail

type TransactionDetail struct {
	ID            int64
	TransactionID int64
	ProductID     int64
	Quantity      int64
	Subtotal      int64
}

type TransactionDetailResponse struct {
	ID            int64  `json:"id"`
	TransactionID int64  `json:"transaction_id"`
	ProductID     int64  `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int64  `json:"quantity"`
	Subtotal      int64  `json:"subtotal"`
}
