package model

type PaginationQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalData  int64 `json:"total_data"`
	TotalPages int   `json:"total_pages"`
}
