package common

type ErrorResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
