package models

// ErrorResponse
type ResponseData struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}
