package utils

import "fmt"

type CustomError struct {
	Inner error
	Msg   string
}

// Implementasi method Error() agar CustomError menjadi error
func (e *CustomError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Inner)
	}
	return e.Msg
}

// Membungkus error
func WrapError(err error, msg string) error {
	return &CustomError{Inner: err, Msg: msg}
}

// Mengembalikan pesan error lengkap
func UnwrapError(err error) string {
	if customErr, ok := err.(*CustomError); ok {
		if customErr.Inner != nil {
			return fmt.Sprintf("%s: %v", customErr.Msg, customErr.Inner)
		}
		return customErr.Msg
	}
	return err.Error()
}
