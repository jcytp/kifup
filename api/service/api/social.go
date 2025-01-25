// service/api/social.go

package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/service/dao"
	"github.com/jcytp/kifup-api/service/model"
)

func LikeKifu(c *gin.Context) (string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// 公開棋譜の存在確認
	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return "Failed to get kifu", err
	}
	if !kifu.IsPublic { // 公開棋譜のみ
		return "Access denied", fmt.Errorf("unauthorized access to private kifu")
	}

	// いいねを追加してカウンターを更新
	if err := dao.InsertKifuLike(kifuID, accountID); err != nil {
		return "Failed to insert like", err
	}
	if err := dao.IncrementKifuLikeCount(kifuID); err != nil {
		return "Failed to update like count", err
	}

	return "", nil
}

func UnlikeKifu(c *gin.Context) (string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// いいねを削除してカウンターを更新
	if err := dao.DeleteKifuLike(kifuID, accountID); err != nil {
		return "Failed to delete kifu like", err
	}
	if err := dao.DecrementKifuLikeCount(kifuID); err != nil {
		return "Failed to update kifu like count", err
	}

	return "", nil
}

type requestPostKifuComment struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

func PostKifuComment(c *gin.Context, req requestPostKifuComment) (string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// 棋譜の存在確認
	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return "Failed to get kifu", err
	}
	if !kifu.IsPublic { // 公開棋譜のみ
		return "Access denied", fmt.Errorf("unauthorized access to private kifu")
	}

	// コメントを追加してカウンターを更新
	comment := &model.KifuComment{
		KifuID:    kifuID,
		AccountID: accountID,
		Content:   req.Content,
	}
	commentID, err := dao.InsertKifuComment(comment)
	if err != nil {
		return "Failed to insert kifu comment", err
	}
	if err := dao.IncrementKifuCommentCount(kifuID); err != nil {
		return "Failed to update kifu comment count", err
	}

	return commentID, nil
}

func ListKifuComments(c *gin.Context) (*[]*model.KifuCommentResponse, string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// 棋譜の存在確認
	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return nil, "Failed to get kifu", err
	}
	if !kifu.IsPublic && kifu.AccountID != accountID { // 非公開の場合は所有者のみ
		return nil, "Access denied", fmt.Errorf("unauthorized access to private kifu")
	}

	// コメント一覧を取得
	comments, err := dao.ListKifuComments(kifuID)
	if err != nil {
		return nil, "Failed to get comments", err
	}

	// レスポンスを構築（アカウント情報も含める）
	responses := make([]*model.KifuCommentResponse, 0, len(comments))
	for _, comment := range comments {
		account, err := dao.GetAccountByID(comment.AccountID)
		if err != nil {
			return nil, "Failed to get account info", err
		}
		responses = append(responses, comment.ToResponse(account))
	}

	return &responses, "", nil
}
