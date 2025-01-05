// service/dao/accounts.go

package dao

import (
	"time"

	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropAccountTable() error {
	query := `
		DROP TABLE IF EXISTS accounts
	`
	_, err := db.Exec(query)
	return err
}

func CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS accounts (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			pass_hash TEXT NOT NULL,
			icon_id TEXT,
			introduction TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login_at TIMESTAMP,
			CHECK (LENGTH(name) >= 2),
			CHECK (LENGTH(email) <= 255),
			CHECK (LENGTH(introduction) <= 1000)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email)
	`
	_, err := db.Exec(query)
	return err
}

func InsertAccount(account *model.Account) (string, error) {
	account.ID = auxi.NewULID()
	now := time.Now()
	account.CreatedAt = now
	account.LastLoginAt = now

	query := `
		INSERT INTO accounts (
			id, name, email, pass_hash,
			icon_id, introduction,
			created_at, last_login_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		account.ID, account.Name, account.Email, account.PassHash,
		account.IconID, account.Introduction,
		account.CreatedAt, account.LastLoginAt,
	)
	return account.ID, err
}

func UpdateAccount(account *model.Account) error {
	query := `
		UPDATE accounts SET
			name = ?, email = ?, pass_hash = ?,
			icon_id = ?, introduction = ?,
			last_login_at = ?
		WHERE id = ?
	`
	_, err := db.Exec(
		query,
		account.Name, account.Email, account.PassHash,
		account.IconID, account.Introduction,
		account.LastLoginAt,
		account.ID,
	)
	return err
}

func DeleteAccount(id string) error {
	query := `DELETE FROM accounts WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func GetAccountByID(id string) (*model.Account, error) {
	query := `SELECT * FROM accounts WHERE id = ?`
	account := &model.Account{}
	err := db.QueryRow(query, id).Scan(
		&account.ID, &account.Name, &account.Email,
		&account.PassHash, &account.IconID, &account.Introduction,
		&account.CreatedAt, &account.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func GetAccountByEmail(email string) (*model.Account, error) {
	query := `SELECT * FROM accounts WHERE email = ?`
	account := &model.Account{}
	err := db.QueryRow(query, email).Scan(
		&account.ID, &account.Name, &account.Email,
		&account.PassHash, &account.IconID, &account.Introduction,
		&account.CreatedAt, &account.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}
