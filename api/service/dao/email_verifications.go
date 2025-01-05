// service/dao/email_verifications.go

package dao

import (
	"time"

	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropEmailVerificationTable() error {
	query := `
		DROP TABLE IF EXISTS email_verifications
	`
	_, err := db.Exec(query)
	return err
}

func CreateEmailVerificationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS email_verifications (
			email TEXT PRIMARY KEY,
			verification_code TEXT NOT NULL,
			expired_at TIMESTAMP NOT NULL,
			is_used BOOLEAN NOT NULL DEFAULT false,
			CHECK (LENGTH(email) <= 255),
			CHECK (LENGTH(verification_code) = 6)
		);
		CREATE INDEX IF NOT EXISTS idx_email_verifications_email ON email_verifications(email);
		CREATE INDEX IF NOT EXISTS idx_email_verifications_expired_at ON email_verifications(expired_at)
	`
	_, err := db.Exec(query)
	return err
}

func InsertEmailVerification(ev *model.EmailVerification) error {
	query := `
		INSERT INTO email_verifications (
			email, verification_code, expired_at, is_useed
		) VALUES (?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			verification_code = excluded.verification_code,
			expired_at = excluded.expired_at,
			is_used = excluded.is_used
	`
	_, err := db.Exec(
		query,
		ev.Email, ev.VerificationCode, ev.ExpiredAt, ev.IsUsed,
	)
	return err
}

func GetEmailVerification(email string) (*model.EmailVerification, error) {
	query := `SELECT * FROM email_verifications WHERE email = ?`
	ev := &model.EmailVerification{}
	err := db.QueryRow(query, email).Scan(
		&ev.Email, &ev.VerificationCode, &ev.ExpiredAt, &ev.IsUsed,
	)
	if err != nil {
		return nil, err
	}
	return ev, nil
}

func MarkEmailVerificationAsUsed(email string) error {
	query := `UPDATE email_verifications SET is_used = true WHERE email = ?`
	res, err := db.Exec(query, email)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func DeleteExpiredVerifications() error {
	query := `DELETE FROM email_verifications WHERE expired_at < ?`
	_, err := db.Exec(query, time.Now())
	return err
}
