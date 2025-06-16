package domain

type Filter struct {
	UserID      UserID
	Location    string
	PageSize    int64
	CurrentPage int64
}

type PaginatedResponse[T any] struct {
	Items      []T
	TotalCount uint16
	PageNumber int64
}
