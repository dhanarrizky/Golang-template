package response

import "net/http"

// Success adalah response standar untuk HTTP 2xx
type Success struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
	Meta    any         `json:"meta,omitempty"` // pagination, extra info
}

// =========================
// Core Constructor
// =========================

func NewSuccess(status int, code, msg string, data any, meta ...any) *Success {
	var m any
	if len(meta) > 0 {
		m = meta[0]
	}

	return &Success{
		Status:  status,
		Code:    code,
		Message: msg,
		Data:    data,
		Meta:    m,
	}
}

// =========================
// 2xx Responses
// =========================

func OK(code, msg string, data any) *Success {
	return NewSuccess(http.StatusOK, code, msg, data)
}

func Created(code, msg string, data any) *Success {
	return NewSuccess(http.StatusCreated, code, msg, data)
}

func Accepted(code, msg string, data any) *Success {
	return NewSuccess(http.StatusAccepted, code, msg, data)
}

func NoContent(code, msg string) *Success {
	return NewSuccess(http.StatusNoContent, code, msg, nil)
}

func PartialContent(code, msg string, data any, meta any) *Success {
	return NewSuccess(http.StatusPartialContent, code, msg, data, meta)
}
