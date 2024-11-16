// service/dao/kifu_branches.go

package dao

import (
	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/common/db"
	"github.com/jcytp/kifup-api/service/model"
)

func DropKifuBranchTable() error {
	query := `DROP TABLE IF EXISTS kifu_branches`
	_, err := db.Exec(query)
	return err
}

func CreateKifuBranchTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS kifu_branches (
			id TEXT PRIMARY KEY,
			kifu_id TEXT NOT NULL,
			root_branch_id TEXT,
			root_number INTEGER,
			ending_number INTEGER,
			ending_type INTEGER,
			ending_comment TEXT,
			FOREIGN KEY (kifu_id) REFERENCES kifus(id) ON DELETE CASCADE,
			FOREIGN KEY (root_branch_id) REFERENCES kifu_branches(id) ON DELETE CASCADE,
			CHECK (root_number IS NULL OR root_number >= 0),
			CHECK (ending_number IS NULL OR ending_number > 0),
			CHECK (ending_type IS NULL OR ending_type >= 0),
			CHECK (LENGTH(ending_comment) <= 1000),
			CHECK (
				(root_branch_id IS NULL AND root_number IS NULL) OR
				(root_branch_id IS NOT NULL AND root_number IS NOT NULL)
			)
		);
		CREATE INDEX IF NOT EXISTS idx_kifu_branches_kifu_id ON kifu_branches(kifu_id);
		CREATE INDEX IF NOT EXISTS idx_kifu_branches_root ON kifu_branches(root_branch_id, root_number)
	`
	_, err := db.Exec(query)
	return err
}

func InsertKifuBranch(branch *model.KifuBranch) (string, error) {
	branch.ID = auxi.NewULID()

	query := `
		INSERT INTO kifu_branches (
			id, kifu_id, root_branch_id, root_number,
			ending_number, ending_type, ending_comment
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(
		query,
		branch.ID, branch.KifuID, branch.RootBranchID, branch.RootNumber,
		branch.EndingNumber, branch.EndingType, branch.EndingComment,
	)
	return branch.ID, err
}

// 不使用
// func UpdateKifuBranch(branch *model.KifuBranch) error {
// 	query := `
// 		UPDATE kifu_branches SET
// 			ending_number = ?,
// 			ending_type = ?,
// 			ending_comment = ?
// 		WHERE id = ? AND kifu_id = ?
// 	`
// 	res, err := db.Exec(
// 		query,
// 		branch.EndingNumber, branch.EndingType, branch.EndingComment,
// 		branch.ID, branch.KifuID,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	return db.CheckAffectedRows(res, 1)
// }

func DeleteKifuBranch(branchID string, kifuID string) error {
	query := `DELETE FROM kifu_branches WHERE id = ? AND kifu_id = ?`
	res, err := db.Exec(query, branchID, kifuID)
	if err != nil {
		return err
	}
	return db.CheckAffectedRows(res, 1)
}

// 不使用
// func DeleteKifuBranchAfter(rootBranchID string, rootNumber int64) error {
// 	query := `DELETE FROM kifu_branches WHERE root_branch_id = ? AND root_number >= ?`
// 	_, err := db.Exec(query, rootBranchID, rootNumber)
// 	return err
// }

func ListKifuBranchesByKifuID(kifuID string) ([]*model.KifuBranch, error) {
	query := `
		SELECT * FROM kifu_branches 
		WHERE kifu_id = ? 
		ORDER BY 
			CASE WHEN root_branch_id IS NULL THEN 0 ELSE 1 END, -- メインラインを最初に
			COALESCE(root_number, 0),                           -- 分岐発生手数で並べる
			id                                                  -- 同一手数からの分岐は作成順
	`
	rows, err := db.Query(query, kifuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branches := []*model.KifuBranch{}
	for rows.Next() {
		branch := &model.KifuBranch{}
		err := rows.Scan(
			&branch.ID, &branch.KifuID,
			&branch.RootBranchID, &branch.RootNumber,
			&branch.EndingNumber, &branch.EndingType,
			&branch.EndingComment,
		)
		if err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	return branches, nil
}
