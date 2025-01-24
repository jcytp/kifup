// service/model/Social.go

package model

import "time"

// table: `kifu_likes`
type KifuLike struct {
	KifuID    string    `db:"kifu_id"`
	AccountID string    `db:"account_id"`
	CreatedAt time.Time `db:"created_at"`
}

// table: `kifu_comments`
type KifuComment struct {
	ID        string    `db:"id"`
	KifuID    string    `db:"kifu_id"`
	AccountID string    `db:"account_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

type KifuCommentResponse struct {
	ID        string           `json:"id"`
	KifuID    string           `json:"kifu_id"`
	Account   *AccountResponse `json:"account"`
	Content   string           `json:"content"`
	CreatedAt time.Time        `json:"created_at"`
}

func (t *KifuComment) ToResponse(account *Account) *KifuCommentResponse {
	resp := &KifuCommentResponse{
		ID:        t.ID,
		KifuID:    t.KifuID,
		Account:   account.ToResponse(),
		Content:   t.Content,
		CreatedAt: t.CreatedAt,
	}
	return resp
}
