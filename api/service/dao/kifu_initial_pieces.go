// service/dao/kifu_initial_pieces.go

package dao

import (
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropKifuInitialPieceTable() error {
	query := `DROP TABLE IF EXISTS kifu_initial_pieces`
	_, err := db.Exec(query)
	return err
}

func CreateKifuInitialPieceTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kifu_initial_pieces (
			kifu_id TEXT NOT NULL,
			piece INTEGER NOT NULL,
			is_black BOOLEAN NOT NULL,
			file INTEGER,
			rank INTEGER,
			FOREIGN KEY (kifu_id) REFERENCES kifus(id) ON DELETE CASCADE,
			CHECK (piece > 0),
			CHECK (
				(file IS NULL AND rank IS NULL) OR
				(file IS NOT NULL AND rank IS NOT NULL AND
				 file >= 1 AND file <= 9 AND
				 rank >= 1 AND rank <= 9)
			)
		);
		CREATE INDEX IF NOT EXISTS idx_kifu_initial_pieces_kifu_id ON kifu_initial_pieces(kifu_id)
	`
	_, err := db.Exec(query)
	return err
}

func InsertKifuInitialPieces(pieces []*model.KifuInitialPiece) error {
	if len(pieces) == 0 {
		return nil
	}

	query := `
		INSERT INTO kifu_initial_pieces (
			kifu_id, piece, is_black, file, rank
		) VALUES (?, ?, ?, ?, ?)
	`
	for _, p := range pieces {
		_, err := db.Exec(
			query,
			p.KifuID, p.Piece, p.IsBlack, p.File, p.Rank,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteKifuInitialPiecesByKifuID(kifuID string) error {
	query := `DELETE FROM kifu_initial_pieces WHERE kifu_id = ?`
	_, err := db.Exec(query, kifuID)
	return err
}

func ListKifuInitialPiecesByKifuID(kifuID string) ([]*model.KifuInitialPiece, error) {
	query := `
		SELECT * FROM kifu_initial_pieces
		WHERE kifu_id = ?
		ORDER BY 
			CASE 
				WHEN file IS NOT NULL THEN 0                -- 盤上の駒を最初に
				WHEN is_black = true THEN 1                 -- 次に先手の持ち駒
				ELSE 2                                      -- 最後に後手の持ち駒
			END,
			CASE
				WHEN file IS NOT NULL THEN rank * 10 + file -- 盤上の駒は段の昇順->筋の昇順
				ELSE piece                                  -- 持ち駒は駒種の昇順
			END
	`
	rows, err := db.Query(query, kifuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pieces := []*model.KifuInitialPiece{}
	for rows.Next() {
		p := &model.KifuInitialPiece{}
		err := rows.Scan(
			&p.KifuID, &p.Piece, &p.IsBlack, &p.File, &p.Rank,
		)
		if err != nil {
			return nil, err
		}
		pieces = append(pieces, p)
	}
	return pieces, nil
}
