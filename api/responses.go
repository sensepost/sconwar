package api

// ErrorResponse is an error response
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error_detail"`
}

// StatusResponse is a status response
type StatusResponse struct {
	Success bool `json:"success"`
}
