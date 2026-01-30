package product

// DB / internal
type Product struct {
	ID         int64  `json:"id,omitempty"`
	Name       string `json:"name"`
	Price      int64  `json:"price"`
	Stock      int64  `json:"stock"`
	CategoryID int64  `json:"category_id"`
}
