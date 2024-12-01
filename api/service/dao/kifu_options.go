// service/dao/kifu_options.go

package dao

import (
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropKifuOptionTable() error {
	query := `DROP TABLE IF EXISTS kifu_options`
	_, err := db.Exec(query)
	return err
}

func CreateKifuOptionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kifu_options (
			kifu_id TEXT NOT NULL,
			name TEXT NOT NULL,
			value TEXT NOT NULL,
			PRIMARY KEY (kifu_id, name),
			FOREIGN KEY (kifu_id) REFERENCES kifus(id) ON DELETE CASCADE,
			CHECK (LENGTH(name) >= 1 AND LENGTH(name) <= 50),
			CHECK (LENGTH(value) <= 100)
		);
		CREATE INDEX IF NOT EXISTS idx_kifu_options_kifu_id ON kifu_options(kifu_id)
	`
	_, err := db.Exec(query)
	return err
}

func InsertKifuOptions(options []*model.KifuOption) error {
	if len(options) == 0 {
		return nil
	}

	query := `
		INSERT INTO kifu_options (kifu_id, name, value)
		VALUES (?, ?, ?)
	`
	for _, opt := range options {
		_, err := db.Exec(query, opt.KifuID, opt.Name, opt.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func ClearKifuOptionsByKifuID(kifuID string) error {
	query := `DELETE FROM kifu_options WHERE kifu_id = ?`
	_, err := db.Exec(query, kifuID)
	return err
}

func ListKifuOptionsByKifuID(kifuID string) ([]*model.KifuOption, error) {
	query := `SELECT * FROM kifu_options WHERE kifu_id = ?`
	rows, err := db.Query(query, kifuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := []*model.KifuOption{}
	for rows.Next() {
		opt := &model.KifuOption{}
		err := rows.Scan(&opt.KifuID, &opt.Name, &opt.Value)
		if err != nil {
			return nil, err
		}
		options = append(options, opt)
	}
	return options, nil
}
