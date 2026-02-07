package product

import (
	"base-skeleton/internal/shared/errors"
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	FindAll(
		size, offset int,
		search, sort, order string,
	) ([]ProductResponse, int64, error)
	FindByID(id int64) (ProductDetailResponse, error)

	Create(p Product) (Product, error)
	Update(id int64, p Product) (Product, error)
	Delete(id int64) error

	FindByIDForUpdateTx(tx *sql.Tx, id int64) (Product, error)
	DecreaseStockTx(tx *sql.Tx, id int64, qty int64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll(
	size, offset int,
	search, sort, order string,
) ([]ProductResponse, int64, error) {

	sort, order = normalizeSort(sort, order)

	baseQuery := `FROM products`
	var where []string
	var args []interface{}

	if search != "" {
		where = append(where, fmt.Sprintf(
			"LOWER(name) LIKE LOWER($%d)", len(args)+1,
		))
		args = append(args, "%"+search+"%")
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	listQuery := fmt.Sprintf(`
		SELECT id, name, price, stock
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`,
		baseQuery,
		sort,
		order,
		len(args)+1,
		len(args)+2,
	)

	rows, err := r.db.Query(listQuery, append(args, size, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []ProductResponse
	for rows.Next() {
		var p ProductResponse
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) %s`, baseQuery)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *repository) FindByID(id int64) (ProductDetailResponse, error) {
	var res ProductDetailResponse

	err := r.db.QueryRow(`
		SELECT 
			p.id,
			p.name,
			p.price,
			p.stock,
			c.id,
			c.name
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`, id).Scan(
		&res.ID,
		&res.Name,
		&res.Price,
		&res.Stock,
		&res.Category.ID,
		&res.Category.Name,
	)
	if err == sql.ErrNoRows {
		return ProductDetailResponse{}, errors.ErrNotFound
	}

	return res, err
}

func (r *repository) Create(p Product) (Product, error) {
	err := r.db.QueryRow(`
		INSERT INTO products (name, price, stock, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, p.Name, p.Price, p.Stock, p.CategoryID).Scan(&p.ID)

	return p, err
}

func (r *repository) Update(id int64, p Product) (Product, error) {
	res, err := r.db.Exec(`
		UPDATE products
		SET name = $1,
		    price = $2,
		    stock = $3,
		    category_id = $4
		WHERE id = $5
	`, p.Name, p.Price, p.Stock, p.CategoryID, id)

	if err != nil {
		return Product{}, err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return Product{}, sql.ErrNoRows
	}

	p.ID = id
	return p, nil
}

func (r *repository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) FindByIDForUpdateTx(
	tx *sql.Tx,
	id int64,
) (Product, error) {

	var p Product

	err := tx.QueryRow(`
		SELECT
			id,
			name,
			price,
			stock,
			category_id
		FROM products
		WHERE id = $1
		FOR UPDATE
	`, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
	)

	return p, err
}

func (r *repository) DecreaseStockTx(
	tx *sql.Tx,
	id int64,
	qty int64,
) error {

	res, err := tx.Exec(`
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2 AND stock >= $1
	`, qty, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
