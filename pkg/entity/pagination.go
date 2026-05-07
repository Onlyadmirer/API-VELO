package entity

// PaginateMeta menyimpan metadata jumlah halaman dan total data.
type PaginateMeta struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
	Limit       int `json:"limit"`
}

type PaginatedProductResponse struct {
	Data     []Product    `json:"data"`
	Metadata PaginateMeta `json:"metadata"`
}
