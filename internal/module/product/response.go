package product

type ProductResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int64  `json:"stock"`
}

type ProductDetailResponse struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	Price    int64       `json:"price"`
	Stock    int64       `json:"stock"`
	Category CategoryDTO `json:"category"`
}

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
