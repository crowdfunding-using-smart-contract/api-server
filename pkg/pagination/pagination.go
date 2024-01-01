package pagination

import "math"

type PaginateOptions struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type PaginateFindOptions struct {
	Limit int
	Skip  int
}

type PaginateResult[T any] struct {
	Total       int64 `json:"total"`
	From        int   `json:"from"`
	To          int   `json:"to"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	PerPage     int   `json:"per_page"`
	Data        []T   `json:"data"`
} // @name PaginateResult

type MakePaginateContextParameters[T any] struct {
	PaginateOptions PaginateOptions
	CountDocuments  func() int64
	FindDocuments   func(findOptions PaginateFindOptions) []T
}

func MakePaginateResult[T any](parameters MakePaginateContextParameters[T]) PaginateResult[T] {
	paginateOptions := parameters.PaginateOptions
	countDocuments := parameters.CountDocuments
	findDocuments := parameters.FindDocuments

	size := paginateOptions.Size
	if size == 0 {
		size = 20
	}

	currentPage := paginateOptions.Page
	if currentPage == 0 {
		currentPage = 1
	}

	skip := (currentPage - 1) * size
	findOptions := PaginateFindOptions{
		Skip:  skip,
		Limit: size,
	}

	total := countDocuments()
	data := findDocuments(findOptions)

	lastPage := int(math.Ceil(float64(total) / float64(size)))
	to := skip + len(data)
	from := int(math.Min(float64(skip+1), float64(to)))

	return PaginateResult[T]{
		Total:       total,
		From:        from,
		To:          to,
		CurrentPage: currentPage,
		LastPage:    lastPage,
		PerPage:     size,
		Data:        data,
	}
}
