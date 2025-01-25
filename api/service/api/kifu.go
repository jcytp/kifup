// service/api/kifu.go

package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/common/handler"
	"github.com/jcytp/kifup-api/service/api/parser"
	"github.com/jcytp/kifup-api/service/dao"
	"github.com/jcytp/kifup-api/service/model"
)

// ------------------------------------------------------------
type requestCreateKifu struct {
	Type            string  `json:"type" binding:"required,oneof=file position"`
	Content         *string `json:"content,omitempty" binding:"required_if=Type file"`
	InitialPosition *string `json:"initial_position,omitempty"`
}

func CreateKifu(c *gin.Context, req requestCreateKifu) (*string, string, error) {
	aid := handler.GetActorID(c)

	switch req.Type {
	case "file":
		return createKifuFromFile(aid, *req.Content)
	case "position":
		return createKifuFromPosition(aid, (*model.SFEN)(req.InitialPosition))
	default:
		return nil, "Invalid creation type", fmt.Errorf("invalid type: %s", req.Type)
	}
}

func createKifuFromFile(aid string, content string) (*string, string, error) {
	// 1. 棋譜フォーマットごとに棋譜テキストをパースする
	// 　※どの棋譜フォーマットにも合致しなければエラー
	// 2. 生成されたKifu・KifuOption・KifuBranch・KifuMoveをDBに保存する
	var err error
	var parsedKifu *model.ParsedKifu

	parsedKifu, err = parser.ParseFromKIF(content) // フォーマット対象外の場合、nil, nilが返る
	if err != nil {
		return nil, "error in Parsing from KIF", err
	} else if parsedKifu != nil {
		parsedKifu.Kifu.AccountID = aid
		return createKifuFromParsedKifu(parsedKifu) // DBへ保存
	}

	parsedKifu, err = parser.ParseFromCSA(content) // フォーマット対象外の場合、nil, nilが返る
	if err != nil {
		return nil, "error in Parsing from KIF", err
	} else if parsedKifu != nil {
		parsedKifu.Kifu.AccountID = aid
		return createKifuFromParsedKifu(parsedKifu) // DBへ保存
	}

	return nil, "Formats unmatched", fmt.Errorf("content unmatched any kifu formats")
}

func createKifuFromParsedKifu(parsedKifu *model.ParsedKifu) (*string, string, error) {
	kifuID, err := dao.InsertKifu(parsedKifu.Kifu)
	if err != nil {
		return nil, "Failed to create kifu", err
	}

	for i := range parsedKifu.Options {
		parsedKifu.Options[i].KifuID = kifuID
	}
	if err := dao.InsertKifuOptions(parsedKifu.Options); err != nil {
		return nil, "Failed to create kifu options", err
	}

	branchIDMap := make(map[string]string) // 仮ID->実IDの対応表
	for _, branch := range parsedKifu.Branches {
		branch.KifuID = kifuID
		if branch.RootBranchID != nil {
			rootBranchID := branchIDMap[*branch.RootBranchID]
			branch.RootBranchID = &rootBranchID
		}
		branchID, err := dao.InsertKifuBranch(branch.KifuBranch)
		if err != nil {
			return nil, "Failed to create kifu branch", err
		}

		for _, move := range branch.Moves {
			move.BranchID = branchID
		}
		if err := dao.InsertKifuMoves(branch.Moves); err != nil {
			return nil, "Failed to create kifu moves", err
		}
	}

	return &kifuID, "", nil
}

func createKifuFromPosition(aid string, sfen *model.SFEN) (*string, string, error) {
	// SFENの妥当性チェック
	_, err := model.NewBoardPosition(sfen)
	if err != nil {
		return nil, "Invalid initial position", err
	}

	// 棋譜レコードを作成
	kifu := &model.Kifu{
		AccountID:       aid,
		Title:           "新規棋譜",
		IsPublic:        false,
		InitialPosition: sfen,
	}
	kifuID, err := dao.InsertKifu(kifu)
	if err != nil {
		return nil, "Failed to create kifu", err
	}

	// メインラインの空branchを作成
	branch := &model.KifuBranch{
		KifuID: kifuID,
	}
	_, err = dao.InsertKifuBranch(branch)
	if err != nil {
		return nil, "Failed to create branch", err
	}

	return &kifuID, "", nil
}

// ------------------------------------------------------------
type requestListKifus struct {
	Owner *string `form:"owner"`
}

func ListKifus(c *gin.Context, req requestListKifus, pgreq *handler.PaginationRequest) (*[]*model.KifuSummaryResponse, *handler.PaginatedResponse, string, error) {
	limit, offset := pgreq.LimitOffset()

	// ownerに応じてKifuリストを取得
	var totalCount int
	var kifus []*model.Kifu
	var err error
	if req.Owner == nil {
		totalCount, _ = dao.CountPublicKifus()
		kifus, err = dao.ListPublicKifus(limit, offset) // 公開棋譜リスト
	} else {
		if *req.Owner == "me" {
			accountID := handler.GetActorID(c)
			totalCount, _ = dao.CountKifusByAccountID(accountID)
			kifus, err = dao.ListKifusByAccountID(accountID, limit, offset) // 自身の棋譜リスト
		} else {
			totalCount, _ = dao.CountPublicKifusByAccountID(*req.Owner)
			kifus, err = dao.ListPublicKifusByAccountID(*req.Owner, limit, offset) // 特定アカウントの公開棋譜リスト
		}
	}
	if err != nil {
		return nil, nil, "Failed to get kifu list", err
	}

	// レスポンス構築
	responses := make([]*model.KifuSummaryResponse, 0, len(kifus))
	for _, kifu := range kifus {
		owner, err := dao.GetAccountByID(kifu.AccountID)
		if err != nil {
			return nil, nil, "Failed to get account info", err
		}
		tags, err := dao.ListKifuTagsByKifuID(kifu.ID)
		if err != nil {
			return nil, nil, "Failed to get tags", err
		}
		response := kifu.ToSummaryResponse(owner, tags)
		responses = append(responses, response)
	}
	return &responses, pgreq.NewPaginatedResponse(totalCount), "", nil
}

// ------------------------------------------------------------
func GetKifu(c *gin.Context) (*model.KifuDetailResponse, string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return nil, "Failed to get kifu", err
	}

	// 非公開の棋譜は所有者のみアクセス可能
	if !kifu.IsPublic && (accountID != kifu.AccountID) {
		return nil, "Access denied", fmt.Errorf("unauthorized acces to private kifu")
	}

	owner, err := dao.GetAccountByID(kifu.AccountID)
	if err != nil {
		return nil, "Failed to get account info", err
	}
	options, err := dao.ListKifuOptionsByKifuID(kifuID)
	if err != nil {
		return nil, "Failed to get kifu options", err
	}
	tags, err := dao.ListKifuTagsByKifuID(kifuID)
	if err != nil {
		return nil, "Failed to get kifu tags", err
	}
	branches, err := dao.ListKifuBranchesByKifuID(kifuID)
	if err != nil {
		return nil, "Failed to get branches", err
	}
	branchesWithMoves := make([]*model.KifuBranchWithMoves, 0, len(branches))
	for _, branch := range branches {
		moves, err := dao.ListKifuMovesByBranchID(branch.ID)
		if err != nil {
			return nil, "Failed to get moves", err
		}
		branchWithMoves := &model.KifuBranchWithMoves{
			KifuBranch: branch,
			Moves:      moves,
		}
		branchesWithMoves = append(branchesWithMoves, branchWithMoves)
	}
	hasLike := false
	if accountID != "" {
		if like, err := dao.HasKifuLike(kifuID, accountID); err == nil {
			hasLike = like
		}
	}

	response := kifu.ToDetailResponse(owner, options, tags, branchesWithMoves, hasLike)
	return response, "", nil
}

// ------------------------------------------------------------
func DeleteKifu(c *gin.Context) (string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// 存在確認と所有者チェック
	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return "Failed to get kifu", err
	}
	if kifu.AccountID != accountID {
		return "Access denied", fmt.Errorf("unauthorized access")
	}

	// 棋譜の削除
	err = dao.DeleteKifu(kifuID, accountID)
	if err != nil {
		return "Failed to delete kifu", err
	}

	return "", nil
}

// ------------------------------------------------------------
type requestUpdateKifuInfo struct {
	Title    string         `json:"title" binding:"required"`
	IsPublic bool           `json:"is_public"`
	GameInfo model.GameInfo `json:"game_info"` // 対局情報（black_player, white_player, started_at, time_rule, その他オプション）
	Tags     []string       `json:"tags"`      // タグリスト
}

func UpdateKifuInfo(c *gin.Context, req requestUpdateKifuInfo) (string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// 存在確認と所有者チェック
	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return "Failed to get kifu", err
	}
	if kifu.AccountID != accountID {
		return "Access denied", fmt.Errorf("unauthorized access")
	}

	// Kifu更新
	kifu.Title = req.Title
	kifu.IsPublic = req.IsPublic
	kifu.BlackPlayer = req.GameInfo.GetBlackPlayer()
	kifu.WhitePlayer = req.GameInfo.GetWhitePlayer()
	kifu.StartedAt = req.GameInfo.GetStartedAt()
	kifu.TimeRule = req.GameInfo.GetTimeRule()

	err = dao.UpdateKifu(kifu)
	if err != nil {
		return "Failed to update kifu", err
	}

	// KifuOption更新
	if err := dao.ClearKifuOptionsByKifuID(kifuID); err != nil {
		return "Failed to clear existing kifu options", err
	}

	reservedKeys := []string{"先手", "後手", "対局日時", "持ち時間", "秒読み", "秒加算"}
	options := make([]*model.KifuOption, 0, len(req.GameInfo))
	for k, v := range req.GameInfo {
		if auxi.IsInArray(k, reservedKeys) {
			continue
		}
		options = append(options, &model.KifuOption{
			KifuID: kifuID,
			Name:   k,
			Value:  v,
		})
	}
	if err := dao.InsertKifuOptions(options); err != nil {
		return "Failed to insert kifu options", err
	}

	// KifuTag更新
	if err := dao.ClearKifuTagsByKifuID(kifuID); err != nil {
		return "Failed to clear existing kifu tags", err
	}

	tags := make([]*model.KifuTag, 0, len(req.Tags))
	for _, name := range req.Tags {
		tags = append(tags, &model.KifuTag{
			KifuID: kifuID,
			Name:   name,
		})
	}
	if err := dao.InsertKifuTags(tags); err != nil {
		return "Failed to insert kifu tags", err
	}

	return "", nil
}

// ------------------------------------------------------------
type KifuMoveRequest struct {
	Number        int64                  `json:"number"`                   // 何手目か（分岐の場合も初手からカウント）
	Piece         model.PieceType        `json:"piece"`                    // 動いた元の駒種
	FromPlace     model.PiecePlace       `json:"from_place"`               // 移動元の場所
	ToPlace       model.PiecePlace       `json:"to_place"`                 // 移動先の場所
	Promote       *bool                  `json:"promote,omitempty"`        // 成ったか（成らなければNULL）
	CatchPiece    *model.PieceType       `json:"catch_piece,omitempty"`    // 取った駒種（取ってなければNULL）
	DirectionSign *string                `json:"direction_sign,omitempty"` // 方向の符号（無ければNULL）
	Variations    *[]KifuMoveLineRequest `json:"variations,omitempty"`     // この手に変わる分岐
	Comment       *string                `json:"comment"`                  // コメント
	TimeSpentMs   *int64                 `json:"time_spent_ms"`            // 消費時間（ミリ秒）
}

type KifuMoveLineRequest []*KifuMoveRequest

func (move *KifuMoveRequest) ToKifuMove(branchID string) *model.KifuMove {
	return &model.KifuMove{
		BranchID:    branchID,
		Number:      move.Number,
		Piece:       move.Piece,
		FromPlace:   move.FromPlace,
		ToPlace:     move.ToPlace,
		Comment:     move.Comment,
		TimeSpentMs: move.TimeSpentMs,
	}
}

type requestUpdateKifuMoves struct {
	Moves KifuMoveLineRequest `json:"moves"` // メインラインと分岐を含む指し手情報
}

func UpdateKifuMoves(c *gin.Context, req requestUpdateKifuMoves) (string, error) {
	accountID := handler.GetActorID(c)
	kifuID := c.GetString("kifuID")

	// 存在確認と所有者チェック
	kifu, err := dao.GetKifu(kifuID)
	if err != nil {
		return "Failed to get kifu", err
	}
	if kifu.AccountID != accountID {
		return "Access denied", fmt.Errorf("unauthorized access")
	}

	// 既存データの削除
	if err := dao.ClearKifuBranchesByKifuID(kifuID); err != nil {
		return "Failed to clear existing branches", err
	}

	// 初期局面の生成（整合性チェックに使用）
	position, err := model.NewBoardPosition(kifu.InitialPosition)
	if err != nil {
		return "Invalid initial position", err
	}

	// ブランチの保存（IDの取得）と指し手の生成のための全体リスト
	branchWithMovesList := []*model.KifuBranchWithMoves{}

	// メインブランチを保存して全体リストに追加
	mainBranch := &model.KifuBranch{KifuID: kifuID}
	mainBranchID, err := dao.InsertKifuBranch(mainBranch)
	if err != nil {
		return "Failed to insert kifu branch", err
	}
	mainBranch.ID = mainBranchID
	mainBranchWithMoves := &model.KifuBranchWithMoves{
		KifuBranch: mainBranch,
		Moves:      []*model.KifuMove{},
	}
	branchWithMovesList = append(branchWithMovesList, mainBranchWithMoves)

	// メインブランチから再帰的にブランチの保存と指し手の生成
	if msg, err := createBranchWithMovesRecursive(kifuID, &branchWithMovesList, req.Moves, mainBranchWithMoves, position); err != nil {
		return msg, err
	}

	// 指し手の保存
	for _, branchWithMoves := range branchWithMovesList {
		if err := dao.InsertKifuMoves(branchWithMoves.Moves); err != nil {
			return "Failed to insert branch", err
		}
	}

	return "", nil
}

// ブランチに対する指し手生成の再帰処理
func createBranchWithMovesRecursive(kifuID string, branchWithMovesList *[]*model.KifuBranchWithMoves, moves KifuMoveLineRequest, currentBranchWithMoves *model.KifuBranchWithMoves, position *model.BoardPosition) (string, error) {
	for _, move := range moves {
		// 指し手の整合性チェック
		kifuMove := move.ToKifuMove(currentBranchWithMoves.ID)
		if err := position.Move(kifuMove); err != nil {
			return "Invalid move", err
		}

		// 分岐の処理
		if move.Variations != nil {
			for _, variation := range *move.Variations {
				// variationに対してブランチを作成し、再帰処理を呼び出す
				newBranch := &model.KifuBranch{
					KifuID:       kifuID,
					RootBranchID: &currentBranchWithMoves.ID,
					RootNumber:   &move.Number,
				}
				newBranchID, err := dao.InsertKifuBranch(newBranch)
				if err != nil {
					return "Failed to insert kifu branch", err
				}
				newBranch.ID = newBranchID
				newBranchWithMoves := &model.KifuBranchWithMoves{
					KifuBranch: newBranch,
					Moves:      []*model.KifuMove{},
				}
				(*branchWithMovesList) = append((*branchWithMovesList), newBranchWithMoves) // 全体リストに新ブランチを追加

				// 再帰呼び出し
				if msg, err := createBranchWithMovesRecursive(kifuID, branchWithMovesList, variation, newBranchWithMoves, position.Copy()); err != nil {
					return msg, err
				}
			}
		}

		// 指し手の追加
		currentBranchWithMoves.Moves = append(currentBranchWithMoves.Moves, move.ToKifuMove(currentBranchWithMoves.ID))
	}
	return "", nil
}
