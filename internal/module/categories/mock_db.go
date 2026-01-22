package categories

import "sync"

type mockDB struct {
	mu     sync.Mutex
	data   map[int64]Category
	autoID int64
}

func newMockDB() *mockDB {
	return &mockDB{
		data: map[int64]Category{
			1: {ID: 1, Name: "Electronics", Description: "Electronic items"},
			2: {ID: 2, Name: "Books", Description: "All kinds of books"},
		},
		autoID: 3,
	}
}
