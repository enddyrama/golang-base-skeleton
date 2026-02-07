package category

import (
	appErr "base-skeleton/internal/shared/errors"
	"database/sql"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll(
	size, offset int,
	search, sort, order string,
) ([]Category, int64, *appErr.AppError) {

	res, total, err := s.repo.GetAll(size, offset, search, sort, order)
	if err != nil {
		return nil, 0, appErr.Internal("failed to query categories")
	}

	return res, total, nil
}

func (s *Service) GetByID(id int64) (Category, *appErr.AppError) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		if err == appErr.ErrNotFound {
			return Category{}, appErr.Custom(404, "category not found")
		}
		return Category{}, appErr.Internal("failed to query category")
	}
	return c, nil
}

func (s *Service) Create(c Category) (Category, *appErr.AppError) {
	if strings.TrimSpace(c.Name) == "" {
		return Category{}, appErr.BadRequest("name is required")
	}

	res, err := s.repo.Create(c)
	if err != nil {
		return Category{}, appErr.Internal("failed to create category")
	}

	return res, nil
}

func (s *Service) Update(id int64, c Category) (Category, *appErr.AppError) {
	res, err := s.repo.Update(id, c)

	if err == sql.ErrNoRows {
		return Category{}, appErr.Custom(404, "category not found")
	}
	if err != nil {
		return Category{}, appErr.Internal("failed to update category")
	}

	return res, nil
}

func (s *Service) Delete(id int64) *appErr.AppError {
	err := s.repo.Delete(id)

	if err == sql.ErrNoRows {
		return appErr.Custom(404, "category not found")
	}
	if err != nil {
		return appErr.Internal("failed to delete category")
	}

	return nil
}
