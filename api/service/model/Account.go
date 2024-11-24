// service/model/Account.go

package model

import (
	"time"
)

// table: `accounts`
type Account struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PassHash     string    `db:"pass_hash"`
	IconID       string    `db:"icon_id"`      // アイコン画像のID
	Introduction string    `db:"introduction"` // 自己紹介
	CreatedAt    time.Time `db:"created_at"`
	LastLoginAt  time.Time `db:"last_login_at"`
}

// ------------------------------------------------------------

type AccountResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	IconID       string    `json:"icon_id"`
	Introduction string    `json:"introduction"`
	CreatedAt    time.Time `json:"created_at"`
	LastLoginAt  time.Time `json:"last_login_at"`
}

func (t *Account) ToResponse() *AccountResponse {
	resp := &AccountResponse{
		ID:           t.ID,
		Name:         t.Name,
		IconID:       t.IconID,
		Introduction: t.Introduction,
		CreatedAt:    t.CreatedAt,
		LastLoginAt:  t.LastLoginAt,
	}
	return resp
}
