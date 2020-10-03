package api

type ErrorResponse struct {
	Message string
	Error   string
}

type StatusResponse struct {
	Success bool
}
