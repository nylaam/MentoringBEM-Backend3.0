package dto

type ErrorDetail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type APIResponse struct {
	Success bool         `json:"success"`
	Data    any          `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

func Success(data any) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func Fail(code, message, requestID string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}
}
