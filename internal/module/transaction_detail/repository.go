package transactiondetail

import "database/sql"

type Repository interface {
	InsertTx(tx *sql.Tx, d TransactionDetail) error
	FindByTransactionID(transactionID int64) ([]TransactionDetailResponse, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) InsertTx(
	tx *sql.Tx,
	d TransactionDetail,
) error {

	_, err := tx.Exec(`
		INSERT INTO transaction_details
			(transaction_id, product_id, quantity, subtotal)
		VALUES ($1, $2, $3, $4)
	`,
		d.TransactionID,
		d.ProductID,
		d.Quantity,
		d.Subtotal,
	)

	return err
}

func (r *repository) FindByTransactionID(
	transactionID int64,
) ([]TransactionDetailResponse, error) {

	rows, err := r.db.Query(`
		SELECT
			td.id,
			td.transaction_id,
			td.product_id,
			p.name,
			td.quantity,
			td.subtotal
		FROM transaction_details td
		JOIN product p ON p.id = td.product_id
		WHERE td.transaction_id = $1
	`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TransactionDetailResponse
	for rows.Next() {
		var d TransactionDetailResponse
		if err := rows.Scan(
			&d.ID,
			&d.TransactionID,
			&d.ProductID,
			&d.ProductName,
			&d.Quantity,
			&d.Subtotal,
		); err != nil {
			return nil, err
		}
		result = append(result, d)
	}

	return result, nil
}
