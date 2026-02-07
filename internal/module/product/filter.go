package product

import "strings"

var allowedSortFields = map[string]string{
	"id":         "id",
	"name":       "name",
	"price":      "price",
	"stock":      "stock",
	"categoryId": "category_id",
}

func normalizeSort(sort, order string) (string, string) {
	sort = strings.ToLower(sort)
	if _, ok := allowedSortFields[sort]; !ok {
		sort = "id"
	}

	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	return sort, order
}
