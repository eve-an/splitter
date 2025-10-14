package handler

type ErrorResponse struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func NewErrorResponse(message string, details ...any) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
		Details: details,
	}
}
