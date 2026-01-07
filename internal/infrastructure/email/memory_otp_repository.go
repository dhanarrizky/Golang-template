package email

import (
	"sync"

	"github.com/dhanarrizky/Golang-template/internal/domain/valueobjects"
)

type MemoryOTPRepository struct {
	mu   sync.Mutex
	data map[string]valueobjects.OTP
}

func NewMemoryOTPRepository() *MemoryOTPRepository {
	return &MemoryOTPRepository{
		data: make(map[string]valueobjects.OTP),
	}
}

func (r *MemoryOTPRepository) Save(email string, otp valueobjects.OTP) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[email] = otp
	return nil
}

func (r *MemoryOTPRepository) Find(email string) (valueobjects.OTP, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.data[email]
	return v, ok
}

func (r *MemoryOTPRepository) Delete(email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, email)
	return nil
}
