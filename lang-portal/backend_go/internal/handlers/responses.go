package handlers

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

type ListResponse struct {
	Items      interface{}        `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
} 