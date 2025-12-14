package response

type APIResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
	Errors  any         `json:"errors,omitempty"`
	Meta    any         `json:"meta,omitempty"`
}
