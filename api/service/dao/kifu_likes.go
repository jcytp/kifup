// service/dao/kifu_likes.go

package dao

import (
	"github.com/jcytp/kifup-api/common/db"
)

func CreateKifuLikeTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS kifu_likes (
            kifu_id TEXT NOT NULL,
            account_id TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (kifu_id, account_id),
            FOREIGN KEY (kifu_id) REFERENCES kifus(id) ON DELETE CASCADE,
            FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
        );
        CREATE INDEX IF NOT EXISTS idx_kifu_likes_kifu_id ON kifu_likes(kifu_id);
        CREATE INDEX IF NOT EXISTS idx_kifu_likes_account_id ON kifu_likes(account_id)
    `
	_, err := db.Exec(query)
	return err
}

func InsertKifuLike(kifuID string, accountID string) error {
	query := `
		INSERT INTO kifu_likes (kifu_id, account_id)
		VALUES (?, ?)
	`
	_, err := db.Exec(query, kifuID, accountID)
	return err
}

func DeleteKifuLike(kifuID string, accountID string) error {
	query := `
		DELETE FROM kifu_likes
		WHERE kifu_id = ? AND account_id = ?
	`
	res, err := db.Exec(query, kifuID, accountID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func HasKifuLike(kifuID string, accountID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM kifu_likes
		WHERE kifu_id = ? AND account_id = ?
	`
	var count int
	if err := db.QueryRow(query, kifuID, accountID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
