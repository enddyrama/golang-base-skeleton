package category

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
) ([]Category, int64, *appErr.AppError) {

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
		FROM category
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
	// Query data
	// ======================
	listQuery := fmt.Sprintf(`
		SELECT id, name, description
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
		return nil, 0, appErr.Internal("failed to query categories")
	}
	defer rows.Close()

	var result []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, 0, appErr.Internal("failed to scan category")
		}
		result = append(result, c)
	}

	// ======================
	// Query total (SAFE)
	// ======================
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		%s
	`, baseQuery)

	var total int64
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, appErr.Internal("failed to count categories")
	}

	return result, total, nil
}

func (s *Service) GetByID(id int64) (Category, *appErr.AppError) {
	var c Category

	err := s.db.QueryRow(`
		SELECT id, name, description
		FROM categories
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Description)

	if err == sql.ErrNoRows {
		return Category{}, appErr.Custom(404, "category not found")
	}
	if err != nil {
		return Category{}, appErr.Internal("failed to query category")
	}

	return c, nil
}

func (s *Service) Create(c Category) (Category, *appErr.AppError) {
	if c.Name == "" {
		return Category{}, appErr.BadRequest("name is required")
	}

	err := s.db.QueryRow(`
		INSERT INTO category (name, description)
		VALUES ($1, $2)
		RETURNING id
	`, c.Name, c.Description).Scan(&c.ID)

	if err != nil {
		return Category{}, appErr.Internal("failed to create category")
	}

	return c, nil
}

func (s *Service) Update(id int64, c Category) (Category, *appErr.AppError) {
	res, err := s.db.Exec(`
		UPDATE category
		SET name = $1, description = $2
		WHERE id = $3
	`, c.Name, c.Description, id)

	if err != nil {
		return Category{}, appErr.Internal("failed to update category")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return Category{}, appErr.Custom(404, "category not found")
	}

	c.ID = id
	return c, nil
}

func (s *Service) Delete(id int64) *appErr.AppError {
	res, err := s.db.Exec(`
		DELETE FROM category
		WHERE id = $1
	`, id)

	if err != nil {
		return appErr.Internal("failed to delete category")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return appErr.Custom(404, "category not found")
	}

	return nil
}
