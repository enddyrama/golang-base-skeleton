package response

type Pagination struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"totalPage"`
}
