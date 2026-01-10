package dto

type Pagination struct {
	CurrentPage  *int `json:"current_page,omitempty"`
	Limit        *int `json:"limit,omitempty"`
	TotalPages   *int `json:"total_pages,omitempty"`
	TotalRecords int  `json:"total_records,omitempty"`
}
