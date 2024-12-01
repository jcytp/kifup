// service/dao/kifu_moves.go

package dao

import (
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropKifuMoveTable() error {
	query := `DROP TABLE IF EXISTS kifu_moves`
	_, err := db.Exec(query)
	return err
}

func CreateKifuMoveTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kifu_moves (
			branch_id TEXT NOT NULL,
			number INTEGER NOT NULL,
			piece INTEGER NOT NULL,
			from_place INTEGER NOT NULL,
			to_place INTEGER NOT NULL,
			comment TEXT,
			time_spent_ms INTEGER,
			PRIMARY KEY (branch_id, number),
			FOREIGN KEY (branch_id) REFERENCES kifu_branches(id) ON DELETE CASCADE,
			CHECK (number > 0),
			CHECK (piece >= 0 AND piece <= 15),
			CHECK (from_place >= 0 AND from_place <= 255),
			CHECK (to_place >= 0 AND to_place <= 255),
			CHECK (LENGTH(comment) <= 1000),
			CHECK (time_spent_ms IS NULL OR time_spent_ms >= 0)
		);
		CREATE INDEX IF NOT EXISTS idx_kifu_moves_branch_number ON kifu_moves(branch_id, number)
	`
	_, err := db.Exec(query)
	return err
}

func InsertKifuMoves(moves []*model.KifuMove) error {
	if len(moves) == 0 {
		return nil
	}

	query := `
		INSERT INTO kifu_moves (
			branch_id, number, piece,
			from_place, to_place,
			comment, time_spent_ms
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, move := range moves {
		_, err := db.Exec(
			query,
			move.BranchID, move.Number, move.Piece,
			move.FromPlace, move.ToPlace,
			move.Comment, move.TimeSpentMs,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListKifuMovesByBranchID(branchID string) ([]*model.KifuMove, error) {
	query := `
		SELECT * FROM kifu_moves 
		WHERE branch_id = ? 
		ORDER BY number
	`
	rows, err := db.Query(query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moves := []*model.KifuMove{}
	for rows.Next() {
		move := &model.KifuMove{}
		err := rows.Scan(
			&move.BranchID, &move.Number, &move.Piece,
			&move.FromPlace, &move.ToPlace,
			&move.Comment, &move.TimeSpentMs,
		)
		if err != nil {
			return nil, err
		}
		moves = append(moves, move)
	}
	return moves, nil
}
