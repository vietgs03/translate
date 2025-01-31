package types

// APIError represents an error response
type APIError struct {
	Error string `json:"error" example:"error message"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Status  string      `json:"status" example:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
}

// LoginResponse represents the response for login
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Page     int `json:"page" example:"1"`
		PageSize int `json:"page_size" example:"10"`
		Total    int `json:"total" example:"100"`
	} `json:"pagination"`
} 