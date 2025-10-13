package dto

type ErrorDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorDTO(code int, msg string) *ErrorDTO {
	return &ErrorDTO{Code: code, Message: msg}
}
