package transactiondetail

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	InsertTx(tx *sql.Tx, d TransactionDetail) error
	InsertManyTx(tx *sql.Tx, details []TransactionDetail) error
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

func (r *repository) InsertManyTx(
	tx *sql.Tx,
	details []TransactionDetail,
) error {
	if len(details) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(details))
	valueArgs := make([]interface{}, 0, len(details)*4)

	for i, d := range details {
		// PostgreSQL parameter placeholders: $1, $2, ...
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs,
			d.TransactionID,
			d.ProductID,
			d.Quantity,
			d.Subtotal,
		)
	}

	stmt := fmt.Sprintf(`
		INSERT INTO transaction_details
		(transaction_id, product_id, quantity, subtotal)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := tx.Exec(stmt, valueArgs...)
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
