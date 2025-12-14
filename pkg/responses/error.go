package response

import "net/http"

// AppError adalah error standar aplikasi
// yang akan ditangkap oleh middleware
type AppError struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Errors  any         `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// =========================
// Core Constructor
// =========================

func New(status int, code, message string, errs ...any) *AppError {
	var err any
	if len(errs) > 0 {
		err = errs[0]
	}

	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Errors:  err,
	}
}

// =========================
// 4xx Errors
// =========================

func BadRequest(code, msg string, errs ...any) *AppError {
	return New(http.StatusBadRequest, code, msg, errs...)
}

func Unauthorized(code, msg string) *AppError {
	return New(http.StatusUnauthorized, code, msg)
}

func Forbidden(code, msg string) *AppError {
	return New(http.StatusForbidden, code, msg)
}

func NotFound(code, msg string) *AppError {
	return New(http.StatusNotFound, code, msg)
}

func Conflict(code, msg string) *AppError {
	return New(http.StatusConflict, code, msg)
}

func Unprocessable(code, msg string, errs ...any) *AppError {
	return New(http.StatusUnprocessableEntity, code, msg, errs...)
}

// =========================
// 5xx Errors
// =========================

func Internal(code string) *AppError {
	return New(
		http.StatusInternalServerError,
		code,
		"Internal server error",
	)
}

func ServiceUnavailable(code, msg string) *AppError {
	return New(http.StatusServiceUnavailable, code, msg)
}
