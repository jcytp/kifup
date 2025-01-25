// service/dao/kifu_comments.go

package dao

import (
	"time"

	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func CreateKifuCommentTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS kifu_comments (
            id TEXT PRIMARY KEY,
            kifu_id TEXT NOT NULL,
            account_id TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (kifu_id) REFERENCES kifus(id) ON DELETE CASCADE,
            FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
            CHECK (LENGTH(content) >= 1 AND LENGTH(content) <= 1000)
        );
        CREATE INDEX IF NOT EXISTS idx_kifu_comments_kifu_id ON kifu_comments(kifu_id);
        CREATE INDEX IF NOT EXISTS idx_kifu_comments_account_id ON kifu_comments(account_id)
    `
	_, err := db.Exec(query)
	return err
}

func InsertKifuComment(comment *model.KifuComment) (string, error) {
	comment.ID = auxi.NewULID()
	comment.CreatedAt = time.Now()

	query := `
        INSERT INTO kifu_comments (
			id, kifu_id, account_id,
			content, created_at
		) VALUES (?, ?, ?, ?, ?)
    `
	_, err := db.Exec(
		query,
		comment.ID, comment.KifuID, comment.AccountID,
		comment.Content, comment.CreatedAt,
	)
	return comment.ID, err
}

func ListKifuComments(kifuID string) ([]*model.KifuComment, error) {
	query := `
        SELECT * FROM kifu_comments
        WHERE kifu_id = ?
        ORDER BY created_at ASC
    `
	rows, err := db.Query(query, kifuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.KifuComment{}
	for rows.Next() {
		comment := &model.KifuComment{}
		err := rows.Scan(
			&comment.ID, &comment.KifuID, &comment.AccountID,
			&comment.Content, &comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
