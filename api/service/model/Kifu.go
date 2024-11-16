// service/model/Kifu.go

package model

import (
	"fmt"
	"strings"
	"time"
)

// table: `kifus`
type Kifu struct {
	ID              string     `db:"id"`
	AccountID       string     `db:"account_id"`
	Title           string     `db:"title"`
	IsPublic        bool       `db:"is_public"`
	BlackPlayer     *string    `db:"black_player"`      // 先手の名前
	WhitePlayer     *string    `db:"white_player"`      // 後手の名前
	StartedAt       *time.Time `db:"started_at"`        // 開始日時
	EndedAt         *time.Time `db:"ended_at"`          // 終了日時
	TimeRule        *string    `db:"time_rule"`         // 持ち時間
	StartsWithBlack *bool      `db:"starts_with_black"` // 初手の手番
	OriginalFormat  string     `db:"original_format"`   // 新規作成時の入力フォーマット
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
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

// table: `kifu_initial_pieces`
type KifuInitialPiece struct {
	KifuID  string    `db:"kifu_id"`
	Piece   PieceType `db:"piece"`    // 駒種
	IsBlack bool      `db:"is_black"` // 先後
	File    *int64    `db:"file"`     // 筋（持ち駒ならNULL）
	Rank    *int64    `db:"rank"`     // 段（持ち駒ならNULL）
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

// 分岐についての補足
// - 分岐の数や深さに論理的な制限は設けない
// - 実質的な制限は以下の要因による
//   - APIリクエストの最大サイズ
//   - データベースの容量制限
//   - クライアント側での表示・操作の制約
// - 分岐の表現
//   - メインライン：RootBranchID=null, RootNumber=null
//   - 派生ライン：RootBranchID=分岐元のID, RootNumber=分岐が発生した手数

// table: `kifu_moves`
type KifuMove struct {
	BranchID    string    `db:"branch_id"`     // 分岐ID
	Number      int64     `db:"number"`        // 何手目か（分岐の場合も初手からカウント）
	Piece       PieceType `db:"piece"`         // 動いた結果の駒種
	ToFile      int64     `db:"to_file"`       // 移動先の筋
	ToRank      int64     `db:"to_rank"`       // 移動先の段
	FromFile    *int64    `db:"from_file"`     // 移動元の筋（打った場合はNULL）
	FromRank    *int64    `db:"from_rank"`     // 移動元の段（打った場合はNULL）
	Comment     *string   `db:"comment"`       // コメント
	TimeSpentMs *int64    `db:"time_spent_ms"` // 消費時間（ミリ秒）
}

type KifuBranchWithMoves struct {
	KifuBranch
	Moves []*KifuMove
}

// ------------------------------------------------------------
// リスト表示用のレスポンス

type KifuSummaryResponse struct {
	ID        string            `json:"id"`
	AccountID string            `json:"account_id"`
	Title     string            `json:"title"`
	IsPublic  bool              `json:"is_public"`
	UpdatedAt time.Time         `json:"updated_at"`
	GameInfo  map[string]string `json:"game_info"` // 対局情報
	Tags      []string          `json:"tags"`      // タグリスト
}

func (t *Kifu) ToSummaryResponse(kifuTags []*KifuTag) *KifuSummaryResponse {
	resp := &KifuSummaryResponse{
		ID:        t.ID,
		AccountID: t.AccountID,
		Title:     t.Title,
		IsPublic:  t.IsPublic,
		UpdatedAt: t.UpdatedAt,
		GameInfo:  t.buildSummaryGameInfo(),
		Tags:      t.buildTags(kifuTags),
	}
	return resp
}

// ------------------------------------------------------------
// 詳細表示用のレスポンス

type PiecePositionResponse struct {
	Piece    PieceType `json:"piece"`          // 駒種
	IsBlack  *bool     `json:"is_black"`       // 先後（駒箱ならNULL）
	IsInHand *bool     `json:"is_in_hand"`     // 持ち駒（駒箱ならNULL）
	File     *int64    `json:"file,omitempty"` // 筋（持ち駒や駒箱ならNULL）
	Rank     *int64    `json:"rank,omitempty"` // 段（持ち駒や駒箱ならNULL）
}

type InitialPositionResponse []*PiecePositionResponse

type KifuMoveResponse struct {
	Number        int64                   `db:"number"`                     // 何手目か（分岐の場合も初手からカウント）
	Piece         PieceType               `db:"piece"`                      // 動いた元の駒種
	ToFile        int64                   `db:"to_file"`                    // 移動先の筋
	ToRank        int64                   `db:"to_rank"`                    // 移動先の段
	FromFile      *int64                  `db:"from_file"`                  // 移動元の筋（打った場合はNULL）
	FromRank      *int64                  `db:"from_rank"`                  // 移動元の段（打った場合はNULL）
	Promote       *bool                   `json:"promote,omitempty"`        // 成ったか（成らなければNULL）
	CatchPiece    *PieceType              `json:"catch_piece,omitempty"`    // 取った駒種（取ってなければNULL）
	DirectionSign *string                 `json:"direction_sign,omitempty"` // 方向の符号（無ければNULL）
	Variations    *[]KifuMoveLineResponse `json:"variations,omitempty"`     // この手に変わる分岐
	Comment       *string                 `db:"comment"`                    // コメント
	TimeSpentMs   *int64                  `db:"time_spent_ms"`              // 消費時間（ミリ秒）
}

type KifuMoveLineResponse []*KifuMoveResponse

type KifuDetailResponse struct {
	ID              string                  `json:"id"`
	AccountID       string                  `json:"account_id"`
	Title           string                  `json:"title"`
	IsPublic        bool                    `json:"is_public"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	GameInfo        map[string]string       `json:"game_info"`        // 対局情報
	Tags            []string                `json:"tags"`             // タグリスト
	InitialPosition InitialPositionResponse `json:"initial_position"` // 初期局面
	Moves           KifuMoveLineResponse    `json:"moves"`            // 指し手（分岐を含む）
}

func (t *Kifu) ToDetailResponse(options []*KifuOption, kifuTags []*KifuTag, pieces []*KifuInitialPiece, branches []*KifuBranchWithMoves) *KifuDetailResponse {
	resp := &KifuDetailResponse{
		ID:              t.ID,
		AccountID:       t.AccountID,
		Title:           t.Title,
		IsPublic:        t.IsPublic,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		GameInfo:        t.buildGameInfo(options),
		Tags:            t.buildTags(kifuTags),
		InitialPosition: t.buildInitialPosition(pieces),
		Moves:           t.buildMoves(branches),
	}
	return resp
}

// ------------------------------------------------------------
// レスポンス変換の補助関数

func (t *Kifu) buildSummaryGameInfo() map[string]string {
	gameInfo := map[string]string{}
	if t.BlackPlayer != nil {
		gameInfo["先手"] = *t.BlackPlayer
	}
	if t.WhitePlayer != nil {
		gameInfo["後手"] = *t.WhitePlayer
	}
	if t.StartedAt != nil {
		gameInfo["対局日時"] = t.StartedAt.Format("2006-01-02 03:04")
	}
	return gameInfo
}

func (t *Kifu) buildGameInfo(options []*KifuOption) map[string]string {
	gameInfo := map[string]string{}
	if t.BlackPlayer != nil {
		gameInfo["先手"] = *t.BlackPlayer
	}
	if t.WhitePlayer != nil {
		gameInfo["後手"] = *t.WhitePlayer
	}
	if t.StartedAt != nil {
		gameInfo["対局開始日時"] = t.StartedAt.Format("2006-01-02 03:04")
	}
	if t.EndedAt != nil {
		gameInfo["対局終了日時"] = t.StartedAt.Format("2006-01-02 03:04")
	}
	if t.TimeRule != nil {
		timeRules := strings.Split(*t.TimeRule, "+")
		if len(timeRules) > 0 && timeRules[0] != "0" {
			gameInfo["持ち時間"] = fmt.Sprintf("%s秒", timeRules[0])
		}
		if len(timeRules) > 1 && timeRules[1] != "0" {
			gameInfo["秒読み"] = fmt.Sprintf("%s秒", timeRules[1])
		}
		if len(timeRules) > 2 && timeRules[2] != "0" {
			gameInfo["秒加算"] = fmt.Sprintf("%s秒", timeRules[2])
		}
	}
	for _, option := range options {
		gameInfo[option.Name] = option.Value
	}
	return gameInfo
}

func (t *Kifu) buildTags(kifuTags []*KifuTag) []string {
	tags := []string{}
	for _, kifuTag := range kifuTags {
		tags = append(tags, kifuTag.Name)
	}
	return tags
}

func (t *Kifu) buildInitialPosition(pieces []*KifuInitialPiece) InitialPositionResponse {
	resp := InitialPositionResponse{}

	if len(pieces) == 0 {
		pieces = InitialPositionHirate
	}

	// 盤上の駒と持ち駒
	usedPieces := map[PieceType]int{}
	for _, p := range pieces {
		isBlack := p.IsBlack
		isInHand := p.File == nil && p.Rank == nil

		resp = append(resp, &PiecePositionResponse{
			Piece:    p.Piece,
			IsBlack:  &isBlack,
			IsInHand: &isInHand,
			File:     p.File,
			Rank:     p.Rank,
		})

		basePiece := p.Piece & ^PIECE_PROMOTE // 元の駒種
		usedPieces[basePiece]++
	}

	// 駒箱の駒を追加
	for pieceType, standardCount := range AllPieceCount {
		remainingCount := standardCount - usedPieces[pieceType]
		for i := 0; i < remainingCount; i++ {
			resp = append(resp, &PiecePositionResponse{Piece: pieceType})
		}
	}

	return resp
}

func (t *Kifu) buildMoves(branches []*KifuBranchWithMoves) KifuMoveLineResponse {
	// メインラインを特定（RootBranchIDがnullのもの）
	var mainBranch *KifuBranchWithMoves
	branchMap := map[string]*KifuBranchWithMoves{}
	for _, branch := range branches {
		branchMap[branch.ID] = branch
		if branch.RootBranchID == nil {
			mainBranch = branch
		}
	}
	if mainBranch == nil {
		return KifuMoveLineResponse{}
	}

	// メインラインから再帰的にKifuMoveResponseを作成
	moveResponses := make([]*KifuMoveResponse, len(mainBranch.Moves))
	for i, branchMove := range mainBranch.Moves {
		moveResponses[i] = branchMove.ToResponse(mainBranch.ID, branches)
	}
	return moveResponses
}

func (t *KifuMove) ToResponse(branchID string, allBranches []*KifuBranchWithMoves) *KifuMoveResponse {
	resp := &KifuMoveResponse{
		Number:      t.Number,
		Piece:       t.Piece,
		ToFile:      t.ToFile,
		ToRank:      t.ToRank,
		FromFile:    t.FromFile,
		FromRank:    t.FromRank,
		Comment:     t.Comment,
		TimeSpentMs: t.TimeSpentMs,
	}

	// 分岐の追加
	variations := []KifuMoveLineResponse{}
	for _, branch := range allBranches {
		if branch.RootBranchID != nil && *branch.RootBranchID == branchID && *branch.RootNumber == t.Number {
			// 分岐ラインから再帰的にKifuMoveResponseを作成
			moveResponses := make([]*KifuMoveResponse, len(branch.Moves))
			for i, branchMove := range branch.Moves {
				moveResponses[i] = branchMove.ToResponse(branch.ID, allBranches)
			}
			variations = append(variations, moveResponses)
		}
	}
	if len(variations) > 0 {
		resp.Variations = &variations
	}

	// 以下の情報は局面の遷移を追跡する必要があるため未実装:
	// - Promote（成り）：直前の局面と現在の駒種を比較して判定
	// - CatchPiece（取られた駒）：直前の局面から消失した駒を特定
	// - DirectionSign（上、寄、引などの符号）：直前の局面での他の同種の駒の位置から判定
	//
	// 実装には以下の処理が必要:
	// 1. 初期局面から指し手を順に適用して各局面を生成
	// 2. 連続する2局面の差分から上記の情報を算出
	// これらは局面生成ロジックの実装後に追加予定

	return resp
}
