// service/model/Kifu.go
// DB上の棋譜データモデルを定義

package model

import (
	"time"
)

// table: `kifus`
type Kifu struct {
	ID              string          `db:"id"`
	AccountID       string          `db:"account_id"`
	Title           string          `db:"title"`
	IsPublic        bool            `db:"is_public"`
	BlackPlayer     *string         `db:"black_player"`     // 先手の名前
	WhitePlayer     *string         `db:"white_player"`     // 後手の名前
	StartedAt       *time.Time      `db:"started_at"`       // 開始日時
	TimeRule        *TimeRuleString `db:"time_rule"`        // 持ち時間
	InitialPosition *SFEN           `db:"initial_position"` // 開始局面（平手初期局面はNULL）
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

// table: `kifu_options`
type KifuOption struct {
	KifuID string `db:"kifu_id"`
	Name   string `db:"name"`
	Value  string `db:"value"`
}

// table: `kifu_tags`
type KifuTag struct {
	KifuID string `db:"kifu_id"`
	Name   string `db:"name"`
}

// table: `kifu_branches`
type KifuBranch struct {
	ID            string      `db:"id"`
	KifuID        string      `db:"kifu_id"`
	RootBranchID  *string     `db:"root_branch_id"` // 分岐元の分岐ID（メインラインはNULL）
	RootNumber    *int64      `db:"root_number"`    // 分岐元の番号（メインラインはNULL、初期局面は0）
	EndingNumber  *int64      `db:"ending_number"`  // 最終手の次の番号
	EndingType    *EndingType `db:"ending_type"`    // 終局の種類（投了／千日手／中断など）
	EndingComment *string     `db:"ending_comment"` // 終局時のコメント
}

// table: `kifu_moves`
type KifuMove struct {
	BranchID    string     `db:"branch_id"`     // 分岐ID
	Number      int64      `db:"number"`        // 何手目か（分岐の場合も初手からカウント）
	Piece       PieceType  `db:"piece"`         // 動いた結果の駒種
	FromPlace   PiecePlace `db:"from_place"`    // 動いた先の場所
	ToPlace     PiecePlace `db:"to_place"`      // 動いた先の場所
	Comment     *string    `db:"comment"`       // コメント
	TimeSpentMs *int64     `db:"time_spent_ms"` // 消費時間（ミリ秒）
}

type KifuBranchWithMoves struct {
	*KifuBranch
	Moves []*KifuMove
}
