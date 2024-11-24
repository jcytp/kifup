// service/dao/kifus.go

package dao

import (
	"time"

	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropKifuTable() error {
	query := `
		DROP TABLE IF EXISTS kifus
	`
	_, err := db.Exec(query)
	return err
}

func CreateKifuTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kifus (
			id TEXT PRIMARY KEY,
			account_id TEXT NOT NULL,
			title TEXT NOT NULL,
			is_public BOOLEAN NOT NULL DEFAULT false,
			black_player TEXT,
			white_player TEXT,
			started_at TIMESTAMP,
			ended_at TIMESTAMP,
			time_rule TEXT,
			starts_with_black BOOLEAN,
			original_format TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
			CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
			CHECK (LENGTH(black_player) <= 100),
			CHECK (LENGTH(white_player) <= 100),
			CHECK (LENGTH(time_rule) <= 100),
			CHECK (ended_at IS NULL OR started_at <= ended_at)
		);
		CREATE INDEX IF NOT EXISTS idx_kifus_account_id ON kifus(account_id);
		CREATE INDEX IF NOT EXISTS idx_kifus_is_public ON kifus(is_public);
		CREATE INDEX IF NOT EXISTS idx_kifus_updated_at ON kifus(updated_at)
	`
	_, err := db.Exec(query)
	return err
}

func InsertKifu(kifu *model.Kifu) (string, error) {
	kifu.ID = auxi.NewULID()
	now := time.Now()
	kifu.CreatedAt = now
	kifu.UpdatedAt = now

	query := `
		INSERT INTO kifus (
			id, account_id, title, is_public,
			black_player, white_player, started_at, ended_at,
			time_rule, starts_with_black, original_format,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		kifu.ID, kifu.AccountID, kifu.Title, kifu.IsPublic,
		kifu.BlackPlayer, kifu.WhitePlayer, kifu.StartedAt, kifu.EndedAt,
		kifu.TimeRule, kifu.StartsWithBlack, kifu.OriginalFormat,
		kifu.CreatedAt, kifu.UpdatedAt,
	)
	return kifu.ID, err
}

func UpdateKifu(kifu *model.Kifu) error {
	kifu.UpdatedAt = time.Now()

	query := `
		UPDATE kifus SET
			title = ?, is_public = ?,
			black_player = ?, white_player = ?, started_at = ?, ended_at = ?,
			time_rule = ?, starts_with_black = ?,
			updated_at = ?
		WHERE id = ? AND account_id = ?
	`
	res, err := db.Exec(
		query,
		kifu.Title, kifu.IsPublic,
		kifu.BlackPlayer, kifu.WhitePlayer, kifu.StartedAt, kifu.EndedAt,
		kifu.TimeRule, kifu.StartsWithBlack,
		kifu.UpdatedAt,
		kifu.ID, kifu.AccountID,
	)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func DeleteKifu(kifuID string, accountID string) error {
	query := `DELETE FROM kifus WHERE id = ? AND account_id = ?`
	res, err := db.Exec(query, kifuID, accountID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func GetKifu(kifuID string) (*model.Kifu, error) {
	query := `SELECT * FROM kifus WHERE id = ?`
	kifu := &model.Kifu{}
	err := db.QueryRow(query, kifuID).Scan(
		&kifu.ID, &kifu.AccountID, &kifu.Title, &kifu.IsPublic,
		&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt, &kifu.EndedAt,
		&kifu.TimeRule, &kifu.StartsWithBlack, &kifu.OriginalFormat,
		&kifu.CreatedAt, &kifu.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return kifu, nil
}

func GetKifusByAccountID(accountID string, limit int, offset int) ([]*model.Kifu, error) {
	query := `
		SELECT * FROM kifus
		WHERE account_id = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kifus := []*model.Kifu{}
	for rows.Next() {
		kifu := &model.Kifu{}
		err := rows.Scan(
			&kifu.ID, &kifu.AccountID, &kifu.Title, &kifu.IsPublic,
			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt, &kifu.EndedAt,
			&kifu.TimeRule, &kifu.StartsWithBlack, &kifu.OriginalFormat,
			&kifu.CreatedAt, &kifu.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		kifus = append(kifus, kifu)
	}
	return kifus, nil
}

func ListPublicKifus(limit int, offset int) ([]*model.Kifu, error) {
	query := `
		SELECT * FROM kifus 
		WHERE is_public = true 
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kifus := []*model.Kifu{}
	for rows.Next() {
		kifu := &model.Kifu{}
		err := rows.Scan(
			&kifu.ID, &kifu.AccountID, &kifu.Title, &kifu.IsPublic,
			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt, &kifu.EndedAt,
			&kifu.TimeRule, &kifu.StartsWithBlack, &kifu.OriginalFormat,
			&kifu.CreatedAt, &kifu.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		kifus = append(kifus, kifu)
	}
	return kifus, nil
}

func ListPublicKifusByAccountID(accountID string, limit int, offset int) ([]*model.Kifu, error) {
	query := `
		SELECT * FROM kifus
		WHERE is_public = true AND account_id = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kifus := []*model.Kifu{}
	for rows.Next() {
		kifu := &model.Kifu{}
		err := rows.Scan(
			&kifu.ID, &kifu.AccountID, &kifu.Title, &kifu.IsPublic,
			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt, &kifu.EndedAt,
			&kifu.TimeRule, &kifu.StartsWithBlack, &kifu.OriginalFormat,
			&kifu.CreatedAt, &kifu.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		kifus = append(kifus, kifu)
	}
	return kifus, nil
}
