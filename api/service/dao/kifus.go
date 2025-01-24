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
			time_rule TEXT,
			initial_position TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            like_count INTEGER NOT NULL DEFAULT 0,
            comment_count INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
			CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
			CHECK (LENGTH(black_player) <= 100),
			CHECK (LENGTH(white_player) <= 100),
			CHECK (LENGTH(time_rule) <= 100),
			CHECK (LENGTH(initial_position) <= 200)
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
	kifu.LikeCount = 0
	kifu.CommentCount = 0

	query := `
		INSERT INTO kifus (
			id, account_id, title, is_public,
			black_player, white_player, started_at,
			time_rule, initial_position,
			created_at, updated_at,
			like_count, comment_count
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		kifu.ID, kifu.AccountID, kifu.Title, kifu.IsPublic,
		kifu.BlackPlayer, kifu.WhitePlayer, kifu.StartedAt,
		kifu.TimeRule, kifu.InitialPosition,
		kifu.CreatedAt, kifu.UpdatedAt,
		kifu.LikeCount, kifu.CommentCount,
	)
	return kifu.ID, err
}

func UpdateKifu(kifu *model.Kifu) error {
	kifu.UpdatedAt = time.Now()

	// 初期局面は変更不可（新規作成が必要）
	query := `
		UPDATE kifus SET
			title = ?, is_public = ?,
			black_player = ?, white_player = ?, started_at = ?,
			time_rule = ?,
			updated_at = ?
		WHERE id = ? AND account_id = ?
	`
	res, err := db.Exec(
		query,
		kifu.Title, kifu.IsPublic,
		kifu.BlackPlayer, kifu.WhitePlayer, kifu.StartedAt,
		kifu.TimeRule,
		kifu.UpdatedAt,
		kifu.ID, kifu.AccountID,
	)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func DeleteKifu(kifuID string, accountID string) error {
	query := `
		DELETE FROM kifus
		WHERE id = ? AND account_id = ?
	`
	res, err := db.Exec(query, kifuID, accountID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func GetKifu(kifuID string) (*model.Kifu, error) {
	query := `
		SELECT * FROM kifus
		WHERE id = ?
	`
	kifu := &model.Kifu{}
	err := db.QueryRow(query, kifuID).Scan(
		&kifu.ID, &kifu.AccountID, &kifu.Title, &kifu.IsPublic,
		&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt,
		&kifu.TimeRule, &kifu.InitialPosition,
		&kifu.CreatedAt, &kifu.UpdatedAt,
		&kifu.LikeCount, &kifu.CommentCount,
	)
	if err != nil {
		return nil, err
	}
	return kifu, nil
}

func CountKifusByAccountID(accountID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM kifus
		WHERE account_id = ?
	`
	var count int
	if err := db.QueryRow(query, accountID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func ListKifusByAccountID(accountID string, limit int, offset int) ([]*model.Kifu, error) {
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
			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt,
			&kifu.TimeRule, &kifu.InitialPosition,
			&kifu.CreatedAt, &kifu.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		kifus = append(kifus, kifu)
	}
	return kifus, nil
}

func CountPublicKifus() (int, error) {
	query := `
		SELECT COUNT(*) FROM kifus
		WHERE is_public = true 
	`
	var count int
	if err := db.QueryRow(query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
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
			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt,
			&kifu.TimeRule, &kifu.InitialPosition,
			&kifu.CreatedAt, &kifu.UpdatedAt,
			&kifu.LikeCount, &kifu.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		kifus = append(kifus, kifu)
	}
	return kifus, nil
}

func CountPublicKifusByAccountID(accountID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM kifus
		WHERE is_public = true AND account_id = ?
	`
	var count int
	if err := db.QueryRow(query, accountID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
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
			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt,
			&kifu.TimeRule, &kifu.InitialPosition,
			&kifu.CreatedAt, &kifu.UpdatedAt,
			&kifu.LikeCount, &kifu.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		kifus = append(kifus, kifu)
	}
	return kifus, nil
}

func IncrementKifuLikeCount(kifuID string) error {
	query := `
		UPDATE kifus
		SET like_count = like_count + 1
		WHERE id = ?
	`
	res, err := db.Exec(query, kifuID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func DecrementKifuLikeCount(kifuID string) error {
	query := `
		UPDATE kifus
		SET like_count = like_count - 1
		WHERE id = ? AND like_count > 0
	`
	res, err := db.Exec(query, kifuID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func IncrementKifuCommentCount(kifuID string) error {
	query := `
		UPDATE kifus
		SET comment_count = comment_count + 1
		WHERE id = ?
	`
	res, err := db.Exec(query, kifuID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}
