package handlers

type PaginatedResponse struct {
    Items      interface{} `json:"items"`
    Pagination struct {
        CurrentPage  int `json:"current_page"`
        TotalPages   int `json:"total_pages"`
        TotalItems   int `json:"total_items"`
        ItemsPerPage int `json:"items_per_page"`
    } `json:"pagination"`
}

func EmptyPaginatedResponse() PaginatedResponse {
    return PaginatedResponse{
        Items: []interface{}{},  // Empty array instead of null
        Pagination: struct {
            CurrentPage  int `json:"current_page"`
            TotalPages   int `json:"total_pages"`
            TotalItems   int `json:"total_items"`
            ItemsPerPage int `json:"items_per_page"`
        }{
            CurrentPage:  1,
            TotalPages:   0,
            TotalItems:   0,
            ItemsPerPage: 100,
        },
    }
} 