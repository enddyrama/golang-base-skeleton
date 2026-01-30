package product

import (
	appErr "base-skeleton/internal/shared/errors"
	"database/sql"
	"fmt"
	"strings"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetAll(
	size, offset int,
	search, sort, order string,
) ([]ProductResponse, int64, *appErr.AppError) {

	// ======================
	// Defaults
	// ======================
	if sort == "" {
		sort = "id"
	}
	if order == "" {
		order = "asc"
	}

	sortColumn, ok := allowedSortFields[sort]
	if !ok {
		sortColumn = "id"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// ======================
	// Build WHERE
	// ======================
	baseQuery := `
		FROM product
	`
	var where []string
	var args []interface{}

	if search != "" {
		where = append(where, "LOWER(name) LIKE LOWER($1)")
		args = append(args, "%"+search+"%")
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	// ======================
	// Query list
	// ======================
	listQuery := fmt.Sprintf(`
		SELECT id, name, price, stock, category_id
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`,
		baseQuery,
		sortColumn,
		order,
		len(args)+1,
		len(args)+2,
	)

	listArgs := append(args, size, offset)

	rows, err := s.db.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, appErr.Internal("failed to query products")
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
			// &p.CategoryID,
		); err != nil {
			return nil, 0, appErr.Internal("failed to scan product")
		}
		result = append(result, p)
	}

	// ======================
	// Query total
	// ======================
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
	`, baseQuery)

	var total int64
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, appErr.Internal("failed to count products")
	}

	return result, total, nil
}

func (s *Service) GetByID(id int64) (ProductDetailResponse, *appErr.AppError) {
	var res ProductDetailResponse

	err := s.db.QueryRow(`
		SELECT 
			p.id,
			p.name,
			p.price,
			p.stock,
			c.id,
			c.name
		FROM product p
		JOIN category c ON c.id = p.category_id
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
		return ProductDetailResponse{}, appErr.Custom(404, "product not found")
	}
	if err != nil {
		return ProductDetailResponse{}, appErr.Internal("failed to query product")
	}

	return res, nil
}

func (s *Service) Create(p Product) (Product, *appErr.AppError) {
	if strings.TrimSpace(p.Name) == "" {
		return Product{}, appErr.BadRequest("name is required")
	}

	if p.Price <= 0 {
		return Product{}, appErr.BadRequest("price is required and must be greater than 0")
	}

	if p.Stock < 0 {
		return Product{}, appErr.BadRequest("stock is required and cannot be negative")
	}

	if p.CategoryID <= 0 {
		return Product{}, appErr.BadRequest("category_id is required")
	}

	err := s.db.QueryRow(`
		INSERT INTO product (name, price, stock, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, p.Name, p.Price, p.Stock, p.CategoryID).Scan(&p.ID)

	if err != nil {
		return Product{}, appErr.Internal("failed to create product")
	}

	return p, nil
}

func (s *Service) Update(id int64, p Product) (Product, *appErr.AppError) {
	res, err := s.db.Exec(`
		UPDATE product
		SET name = $1,
		    price = $2,
		    stock = $3,
		    category_id = $4
		WHERE id = $5
	`, p.Name, p.Price, p.Stock, p.CategoryID, id)

	if err != nil {
		return Product{}, appErr.Internal("failed to update product")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return Product{}, appErr.Custom(404, "product not found")
	}

	p.ID = id
	return p, nil
}

func (s *Service) Delete(id int64) *appErr.AppError {
	res, err := s.db.Exec(`
		DELETE FROM product
		WHERE id = $1
	`, id)

	if err != nil {
		return appErr.Internal("failed to delete product")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return appErr.Custom(404, "product not found")
	}

	return nil
}
