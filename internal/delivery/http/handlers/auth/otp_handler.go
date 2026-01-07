package auth

import (
	"net/http"

	mail "github.com/dhanarrizky/Golang-template/internal/usecase/email"
)

type OTPHandler struct {
	usecase *mail.OTPUsecase
}

func NewOTPHandler(u *mail.OTPUsecase) *OTPHandler {
	return &OTPHandler{u}
}

func (h *OTPHandler) Request(w http.ResponseWriter, r *http.Request) {
	// parse request → call usecase
}

func (h *OTPHandler) Verify(w http.ResponseWriter, r *http.Request) {
	// parse request → call usecase
}
