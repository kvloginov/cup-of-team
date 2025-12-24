package model

// ErrorResponse for API errors
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// HealthResponse for health check
type HealthResponse struct {
	Status string `json:"status"`
}
