package response

// Response represents a standard global API response wrapper
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// PaginatedResponse represents a standard paginated API response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalItems int         `json:"totalItems"`
}
