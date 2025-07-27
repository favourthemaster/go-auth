package utils

import "time"

// Response represents a standard API response structure

type Response struct {
	Success   bool        `json:"success"`           // Status of the response
	Data      interface{} `json:"data,omitempty"`    // Data returned in the response, if any
	Message   string      `json:"message,omitempty"` // Optional message for additional context
	Error     string      `json:"error,omitempty"`   // Optional error message, if any
	TimeStamp string      `json:"timestamp"`         // Timestamp of the response
}

// SuccessResponse creates a successful response
func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Success:   true,
		Data:      data,
		Message:   message,
		TimeStamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// ErrorResponse creates an error response
func ErrorResponse(err error, message string) Response {
	return Response{
		Success:   false,
		Error:     err.Error(),
		Message:   message,
		TimeStamp: time.Now().UTC().Format(time.RFC3339),
	}
}
