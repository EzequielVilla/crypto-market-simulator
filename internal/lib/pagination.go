package lib

type Pagination struct {
	Limit  int
	Offset int
}

func GetPaginationLimitOffset(page int) Pagination {
	limit := 5
	offset := (page - 1) * limit
	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}
