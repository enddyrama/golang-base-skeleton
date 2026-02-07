package product

import (
	"base-skeleton/internal/module/category"
	"base-skeleton/internal/shared/errors"
	appErr "base-skeleton/internal/shared/errors"
	"database/sql"
	"log"
	"strings"
)

type Service struct {
	productRepo  Repository
	categoryRepo category.Repository
}

func NewService(productRepo Repository, categoryRepo category.Repository) *Service {
	return &Service{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *Service) GetAll(
	size, offset int,
	search, sort, order string,
) ([]ProductResponse, int64, *appErr.AppError) {

	if size <= 0 {
		size = 10
	}
	if offset < 0 {
		offset = 0
	}

	search = strings.TrimSpace(search)

	// ✅ normalize sort & order
	sort, order = normalizeSort(sort, order)

	data, total, err := s.productRepo.FindAll(
		size,
		offset,
		search,
		sort,
		order,
	)
	if err != nil {
		return nil, 0, appErr.Internal("failed to query products")
	}

	return data, total, nil
}

func (s *Service) GetByID(id int64) (ProductDetailResponse, *appErr.AppError) {
	if id <= 0 {
		return ProductDetailResponse{}, appErr.BadRequest("invalid product id")
	}

	res, err := s.productRepo.FindByID(id)
	if err == errors.ErrNotFound {
		return ProductDetailResponse{}, appErr.Custom(404, "product not found")
	}
	if err != nil {
		return ProductDetailResponse{}, appErr.Internal("Failed to get product id:" + err.Error())
	}

	return res, nil
}

func (s *Service) Create(p Product) (Product, *appErr.AppError) {
	if strings.TrimSpace(p.Name) == "" {
		return Product{}, appErr.BadRequest("name is required")
	}
	if p.Price <= 0 {
		return Product{}, appErr.BadRequest("price must be greater than 0")
	}
	if p.Stock < 0 {
		return Product{}, appErr.BadRequest("stock cannot be negative")
	}
	if p.CategoryID <= 0 {
		return Product{}, appErr.BadRequest("category_id is required")
	}

	// ✅ Validate category
	_, err := s.categoryRepo.FindByID(p.CategoryID)
	if err != nil {
		log.Println(err) // prints with date and time
		if err == appErr.ErrNotFound {
			return Product{}, appErr.Custom(404, "category not found")
		}
		return Product{}, appErr.Internal("Failed to query category:" + err.Error())
	}

	// Create product
	res, err := s.productRepo.Create(p)
	if err != nil {
		return Product{}, appErr.Internal("failed to create product: " + err.Error())
	}

	return res, nil
}

func (s *Service) Update(id int64, p Product) (Product, *appErr.AppError) {
	if id <= 0 {
		return Product{}, appErr.BadRequest("invalid product id")
	}

	if strings.TrimSpace(p.Name) == "" {
		return Product{}, appErr.BadRequest("name is required")
	}
	if p.Price <= 0 {
		return Product{}, appErr.BadRequest("price must be greater than 0")
	}
	if p.Stock < 0 {
		return Product{}, appErr.BadRequest("stock cannot be negative")
	}
	if p.CategoryID <= 0 {
		return Product{}, appErr.BadRequest("category_id is required")
	}
	// ✅ Validate category
	_, err := s.categoryRepo.FindByID(p.CategoryID)
	if err != nil {
		log.Println(err) // prints with date and time
		if err == appErr.ErrNotFound {
			return Product{}, appErr.Custom(404, "category not found")
		}
		return Product{}, appErr.Internal("Failed to query category:" + err.Error())
	}

	res, err := s.productRepo.Update(id, p)
	log.Println(err)
	if err == sql.ErrNoRows {
		return Product{}, appErr.Custom(404, "product not found")
	}
	if err != nil {
		return Product{}, appErr.Internal("failed to update product" + err.Error())
	}

	return res, nil
}

func (s *Service) Delete(id int64) *appErr.AppError {
	if id <= 0 {
		return appErr.BadRequest("invalid product id")
	}

	err := s.productRepo.Delete(id)
	if err == sql.ErrNoRows {
		return appErr.Custom(404, "product not found")
	}
	if err != nil {
		return appErr.Internal("failed to delete product")
	}

	return nil
}
