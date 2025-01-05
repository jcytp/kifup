// service/model/EmailVerification.go

package model

import (
	"fmt"
	"math/rand"
	"time"
)

// table: `email_verifications`
type EmailVerification struct {
	Email            string    `db:"email"`
	VerificationCode string    `db:"verification_code"`
	ExpiredAt        time.Time `db:"expired_at"`
	IsUsed           bool      `db:"is_used"`
}

func GenerateVerificationCode() string {
	n := rand.Intn(1000000)
	return fmt.Sprintf("%06d", n)
}

func NewEmailVerification(email string) *EmailVerification {
	const timeLimitPeriod = 30 * time.Minute
	code := GenerateVerificationCode()
	expiredAt := time.Now().Add(timeLimitPeriod)
	ev := &EmailVerification{
		Email:            email,
		VerificationCode: code,
		ExpiredAt:        expiredAt,
		IsUsed:           false,
	}
	return ev
}

func (v *EmailVerification) IsValid(code string) bool {
	return !v.IsUsed && time.Now().Before(v.ExpiredAt) && v.VerificationCode == code
}
