// service/dao/kifu_tags.go

package dao

import (
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropKifuTagTable() error {
	query := `DROP TABLE IF EXISTS kifu_tags`
	_, err := db.Exec(query)
	return err
}

func CreateKifuTagTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kifu_tags (
			kifu_id TEXT NOT NULL,
			name TEXT NOT NULL,
			PRIMARY KEY (kifu_id, name),
			FOREIGN KEY (kifu_id) REFERENCES kifus(id) ON DELETE CASCADE,
			CHECK (LENGTH(name) >= 1 AND LENGTH(name) <= 50)
		);
		CREATE INDEX IF NOT EXISTS idx_kifu_tags_kifu_id ON kifu_tags(kifu_id);
		CREATE INDEX IF NOT EXISTS idx_kifu_tags_name ON kifu_tags(name)
	`
	_, err := db.Exec(query)
	return err
}

func InsertKifuTags(tags []*model.KifuTag) error {
	if len(tags) == 0 {
		return nil
	}

	query := `
		INSERT INTO kifu_tags (kifu_id, name)
		VALUES (?, ?)
	`
	for _, tag := range tags {
		_, err := db.Exec(query, tag.KifuID, tag.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteKifuTag(kifuID string, name string) error {
	query := `DELETE FROM kifu_tags WHERE kifu_id = ? AND name = ?`
	res, err := db.Exec(query, kifuID, name)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

func ListKifuTagsByKifuID(kifuID string) ([]*model.KifuTag, error) {
	query := `SELECT * FROM kifu_tags WHERE kifu_id = ?`
	rows, err := db.Query(query, kifuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*model.KifuTag, 0)
	for rows.Next() {
		tag := &model.KifuTag{}
		err := rows.Scan(&tag.KifuID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// 不使用：タグによる棋譜の検索用
// func ListPublicKifusByKifuTag(name string, limit int, offset int) ([]*model.Kifu, error) {
// 	query := `
// 		SELECT k.* FROM kifus k
// 		INNER JOIN kifu_tags t ON k.id = t.kifu_id
// 		WHERE t.name = ? AND k.is_public = true
// 		ORDER BY k.updated_at DESC
// 		LIMIT ? OFFSET ?
// 	`
// 	rows, err := db.Query(query, name, limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	kifus := []*model.Kifu{}
// 	for rows.Next() {
// 		kifu := &model.Kifu{}
// 		err := rows.Scan(
// 			&kifu.ID, &kifu.AccountID, &kifu.Title, &kifu.IsPublic,
// 			&kifu.BlackPlayer, &kifu.WhitePlayer, &kifu.StartedAt, &kifu.EndedAt,
// 			&kifu.TimeRule, &kifu.StartsWithBlack, &kifu.OriginalFormat,
// 			&kifu.CreatedAt, &kifu.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		kifus = append(kifus, kifu)
// 	}
// 	return kifus, nil
// }
