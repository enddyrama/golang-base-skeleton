package category

import (
	"base-skeleton/internal/shared/errors"
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	GetAll(
		size, offset int,
		search, sort, order string,
	) ([]Category, int64, error)

	FindByID(id int64) (Category, error)
	Create(c Category) (Category, error)
	Update(id int64, c Category) (Category, error)
	Delete(id int64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll(
	size, offset int,
	search, sort, order string,
) ([]Category, int64, error) {

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

	baseQuery := ` FROM categories `
	var where []string
	var args []interface{}

	if search != "" {
		where = append(where, "LOWER(name) LIKE LOWER($1)")
		args = append(args, "%"+search+"%")
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

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

	rows, err := r.db.Query(listQuery, append(args, size, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, 0, err
		}
		result = append(result, c)
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) %s`, baseQuery)

	var total int64
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *repository) FindByID(id int64) (Category, error) {
	var c Category
	err := r.db.QueryRow(`SELECT id, name, description FROM categories WHERE id=$1`, id).
		Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return Category{}, errors.ErrNotFound
	}
	return c, err
}

func (r *repository) Create(c Category) (Category, error) {
	err := r.db.QueryRow(`
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
	`, c.Name, c.Description).Scan(&c.ID)

	return c, err
}

func (r *repository) Update(id int64, c Category) (Category, error) {
	res, err := r.db.Exec(`
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
	`, c.Name, c.Description, id)

	if err != nil {
		return Category{}, err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return Category{}, sql.ErrNoRows
	}

	c.ID = id
	return c, nil
}

func (r *repository) Delete(id int64) error {
	res, err := r.db.Exec(`
		DELETE FROM categories
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
