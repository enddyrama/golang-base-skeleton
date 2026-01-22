package categories

import appErr "base-skeleton/internal/shared/errors"

type Service struct {
	db *mockDB
}

func NewService() *Service {
	return &Service{
		db: newMockDB(),
	}
}

func (s *Service) GetAll() ([]Category, *appErr.AppError) {
	s.db.mu.Lock()
	defer s.db.mu.Unlock()

	result := make([]Category, 0, len(s.db.data))
	for _, c := range s.db.data {
		result = append(result, c)
	}
	return result, nil
}

func (s *Service) GetByID(id int64) (Category, *appErr.AppError) {
	s.db.mu.Lock()
	defer s.db.mu.Unlock()

	c, ok := s.db.data[id]
	if !ok {
		return Category{}, appErr.Custom(404, "category not found")
	}
	return c, nil
}

func (s *Service) Create(c Category) (Category, *appErr.AppError) {
	if c.Name == "" {
		return Category{}, appErr.BadRequest("name is required")
	}

	s.db.mu.Lock()
	defer s.db.mu.Unlock()

	c.ID = s.db.autoID
	s.db.autoID++
	s.db.data[c.ID] = c

	return c, nil
}

func (s *Service) Update(id int64, c Category) (Category, *appErr.AppError) {
	s.db.mu.Lock()
	defer s.db.mu.Unlock()

	existing, ok := s.db.data[id]
	if !ok {
		return Category{}, appErr.Custom(404, "category not found")
	}

	existing.Name = c.Name
	existing.Description = c.Description
	s.db.data[id] = existing

	return existing, nil
}

func (s *Service) Delete(id int64) *appErr.AppError {
	s.db.mu.Lock()
	defer s.db.mu.Unlock()

	if _, ok := s.db.data[id]; !ok {
		return appErr.Custom(404, "category not found")
	}

	delete(s.db.data, id)
	return nil
}
