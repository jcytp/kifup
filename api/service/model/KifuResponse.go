// service/model/KifuResponse.go
// APIレスポンス用の棋譜データモデルを定義

package model

import (
	"time"
)

// ------------------------------------------------------------
// リスト表示用のレスポンス

type KifuSummaryResponse struct {
	ID           string            `json:"id"`
	Owner        *AccountResponse  `json:"owner"`
	Title        string            `json:"title"`
	IsPublic     bool              `json:"is_public"`
	UpdatedAt    time.Time         `json:"updated_at"`
	GameInfo     map[string]string `json:"game_info"` // 対局情報
	Tags         []string          `json:"tags"`      // タグリスト
	LikeCount    int64             `json:"like_count"`
	CommentCount int64             `json:"comment_count"`
}

func (t *Kifu) ToSummaryResponse(owner *Account, kifuTags []*KifuTag) *KifuSummaryResponse {
	resp := &KifuSummaryResponse{
		ID:           t.ID,
		Owner:        owner.ToResponse(),
		Title:        t.Title,
		IsPublic:     t.IsPublic,
		UpdatedAt:    t.UpdatedAt,
		GameInfo:     t.buildSummaryGameInfo(),
		Tags:         t.buildTags(kifuTags),
		LikeCount:    t.LikeCount,
		CommentCount: t.CommentCount,
	}
	return resp
}

// ------------------------------------------------------------
// 詳細表示用のレスポンス

type KifuDetailResponse struct {
	ID              string               `json:"id"`
	Owner           *AccountResponse     `json:"owner"`
	Title           string               `json:"title"`
	IsPublic        bool                 `json:"is_public"`
	InitialPosition *SFEN                `json:"initial_position"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	GameInfo        GameInfo             `json:"game_info"` // 対局情報
	Tags            []string             `json:"tags"`      // タグリスト
	Moves           KifuMoveLineResponse `json:"moves"`     // 指し手（分岐を含む）
	LikeCount       int64                `json:"like_count"`
	HasLike         bool                 `json:"has_like"`
}

func (t *Kifu) ToDetailResponse(owner *Account, options []*KifuOption, kifuTags []*KifuTag, branches []*KifuBranchWithMoves, hasLike bool) *KifuDetailResponse {
	resp := &KifuDetailResponse{
		ID:              t.ID,
		Owner:           owner.ToResponse(),
		Title:           t.Title,
		IsPublic:        t.IsPublic,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
		GameInfo:        t.buildGameInfo(options),
		Tags:            t.buildTags(kifuTags),
		InitialPosition: t.InitialPosition,
		Moves:           t.buildMoves(branches),
		LikeCount:       t.LikeCount,
		HasLike:         hasLike,
	}
	return resp
}
