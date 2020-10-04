package api

// ErrorResponse is an error response
type ErrorResponse struct {
	Message string
	Error   string
}

// StatusResponse is a status response
type StatusResponse struct {
	Success bool
}
