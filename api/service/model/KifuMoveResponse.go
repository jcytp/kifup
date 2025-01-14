// service/model/KifuMoveResponse.go

package model

import (
	"log/slog"
)

// --------------------------------------------------------------------------------
type EndingType int64

const (
	ENDING_TORYO EndingType = 0x0 + iota
	ENDING_CHUDAN
	ENDING_SENNICHITE
	ENDING_TIME_UP
	ENDING_ILLEGAL_MOVE
	ENDING_BLACK_ILLEGAL_ACTION
	ENDING_WHITE_ILLEGAL_ACTION
	ENDING_JISHOGI
	ENDING_KACHI
	ENDING_HIKIWAKE
	ENDING_MAX_MOVES
	ENDING_MATTA
	ENDING_TSUMI
	ENDING_FUZUMI
	ENDING_ERROR
)

var EndingTypeName = map[EndingType]string{
	ENDING_TORYO:                "TORYO",
	ENDING_CHUDAN:               "CHUDAN",
	ENDING_SENNICHITE:           "SENNICHITE",
	ENDING_TIME_UP:              "TIME_UP",
	ENDING_ILLEGAL_MOVE:         "ILLEGAL_MOVE",
	ENDING_BLACK_ILLEGAL_ACTION: "+ILLEGAL_ACTION",
	ENDING_WHITE_ILLEGAL_ACTION: "-ILLEGAL_ACTION",
	ENDING_JISHOGI:              "JISHOGI",
	ENDING_KACHI:                "KACHI",
	ENDING_HIKIWAKE:             "HIKIWAKE",
	ENDING_MAX_MOVES:            "MAX_MOVES",
	ENDING_MATTA:                "MATTA",
	ENDING_TSUMI:                "TSUMI",
	ENDING_FUZUMI:               "FUZUMI",
	ENDING_ERROR:                "ERROR",
}

var EndingNameToEndingTypeKIF = map[string]EndingType{
	"中断":   ENDING_CHUDAN,
	"投了":   ENDING_TORYO,
	"持将棋":  ENDING_JISHOGI,
	"千日手":  ENDING_SENNICHITE,
	"切れ負け": ENDING_TIME_UP,
	"反則勝ち": ENDING_KACHI,
	"反則負け": ENDING_ILLEGAL_MOVE,
	"入玉勝ち": ENDING_KACHI,
	"不戦勝":  ENDING_WHITE_ILLEGAL_ACTION,
	"不戦敗":  ENDING_BLACK_ILLEGAL_ACTION,
	"詰み":   ENDING_TSUMI,
	"不詰":   ENDING_FUZUMI,
}

var EndingNameToEndingTypeCSA = map[string]EndingType{
	"TORYO":           ENDING_TORYO,
	"CHUDAN":          ENDING_CHUDAN,
	"SENNICHITE":      ENDING_SENNICHITE,
	"TIME_UP":         ENDING_TIME_UP,
	"ILLEGAL_MOVE":    ENDING_ILLEGAL_MOVE,
	"+ILLEGAL_ACTION": ENDING_BLACK_ILLEGAL_ACTION,
	"-ILLEGAL_ACTION": ENDING_WHITE_ILLEGAL_ACTION,
	"JISHOGI":         ENDING_JISHOGI,
	"KACHI":           ENDING_KACHI,
	"HIKIWAKE":        ENDING_HIKIWAKE,
	"MAX_MOVES":       ENDING_MAX_MOVES,
	"MATTA":           ENDING_MATTA,
	"TSUMI":           ENDING_TSUMI,
	"FUZUMI":          ENDING_FUZUMI,
	"ERROR":           ENDING_ERROR,
}

// --------------------------------------------------------------------------------
type KifuMoveResponse struct {
	Number        int64                   `json:"number"`                   // 何手目か（分岐の場合も初手からカウント）
	Piece         PieceType               `json:"piece"`                    // 動いた元の駒種
	FromPlace     PiecePlace              `json:"from_place"`               // 移動元の場所
	ToPlace       PiecePlace              `json:"to_place"`                 // 移動先の場所
	Promote       *bool                   `json:"promote,omitempty"`        // 成ったか（成らなければNULL）
	CatchPiece    *PieceType              `json:"catch_piece,omitempty"`    // 取った駒種（取ってなければNULL）
	DirectionSign *string                 `json:"direction_sign,omitempty"` // 方向の符号（無ければNULL）
	Variations    *[]KifuMoveLineResponse `json:"variations,omitempty"`     // この手に変わる分岐
	Comment       *string                 `json:"comment"`                  // コメント
	TimeSpentMs   *int64                  `json:"time_spent_ms"`            // 消費時間（ミリ秒）
}

type KifuMoveLineResponse []*KifuMoveResponse

func (t *Kifu) buildMoves(branches []*KifuBranchWithMoves) KifuMoveLineResponse {
	// 開始局面を生成
	position, err := NewBoardPosition(t.InitialPosition)
	if err != nil {
		slog.Error("failed to create board position from initial position", "kifuID", t.ID)
		return KifuMoveLineResponse{}
	}

	// メインラインを特定（RootBranchIDがnullのもの）
	var mainBranch *KifuBranchWithMoves
	for _, branch := range branches {
		if branch.RootBranchID == nil {
			mainBranch = branch
		}
	}
	if mainBranch == nil {
		slog.Error("failed to find main branch", "kifuID", t.ID)
		return KifuMoveLineResponse{}
	}

	// メインラインから再帰的にKifuMoveResponseを作成
	moveResponses := make([]*KifuMoveResponse, len(mainBranch.Moves))
	for i, branchMove := range mainBranch.Moves {
		moveResponses[i] = branchMove.ToResponse(mainBranch.ID, branches, position)
	}
	return moveResponses
}

func (t *KifuMove) ToResponse(branchID string, allBranches []*KifuBranchWithMoves, position *BoardPosition) *KifuMoveResponse {
	slog.Debug("KifuMove.ToResponse", "move", *t)
	resp := &KifuMoveResponse{
		Number:        t.Number,
		Piece:         t.Piece,
		FromPlace:     t.FromPlace,
		ToPlace:       t.ToPlace,
		Promote:       position.IsPromote(t),     // 成判定
		CatchPiece:    position.CatchPiece(t),    // 取った駒
		DirectionSign: position.DirectionSign(t), // 方向を表す符号
		Comment:       t.Comment,
		TimeSpentMs:   t.TimeSpentMs,
	}

	// 局面を進める
	if err := position.Move(t); err != nil {
		return nil
	}

	// 分岐の追加
	variations := []KifuMoveLineResponse{}
	for _, branch := range allBranches {
		if branch.RootBranchID != nil && *branch.RootBranchID == branchID && *branch.RootNumber == t.Number {
			// 分岐ラインから再帰的にKifuMoveResponseを作成
			moveResponses := make([]*KifuMoveResponse, len(branch.Moves))
			for i, branchMove := range branch.Moves {
				moveResponses[i] = branchMove.ToResponse(branch.ID, allBranches, position.Copy())
			}
			variations = append(variations, moveResponses)
		}
	}
	if len(variations) > 0 {
		resp.Variations = &variations
	}

	return resp
}
