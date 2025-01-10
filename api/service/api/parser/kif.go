// service/api/parser/kif.go

package parser

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/service/model"
)

func ParseFromKIF(aid string, content string) (*model.ParsedKifu, error) {
	kifuID := auxi.NewULID()       // dummy id
	mainBranchID := auxi.NewULID() // dummy id
	result := &model.ParsedKifu{
		Kifu: &model.Kifu{
			AccountID:       aid,
			Title:           "新規棋譜",
			IsPublic:        false,
			InitialPosition: model.SfenHirate.PSFEN(),
		},
		Options: []*model.KifuOption{},
		Branches: []*model.KifuBranchWithMoves{
			{
				KifuBranch: &model.KifuBranch{
					ID:           mainBranchID,
					KifuID:       kifuID,
					RootBranchID: nil,
					RootNumber:   nil,
				},
			},
		},
	}
	tmpTimeRule := &model.GameInfo{}
	nextNumber := 1 // 次の指し手番号
	lastPlace := model.PIECE_PLACE_IN_HAND

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "：") {
			// 棋譜情報の行
			if err := parseGameInfoLineForKIF(line, result, tmpTimeRule); err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, fmt.Sprintf("%d ", nextNumber)) {
			// 指し手の行
			if err := parseMoveLineForKIF(line, result.Branches[0], &lastPlace); err != nil {
				return nil, err
			}
			nextNumber++
		}
		// その他は無視
		// 空行、コメント行、テーブルヘッダー行、勝敗宣言の行
	}
	result.Kifu.TimeRule = tmpTimeRule.GetTimeRule()

	return result, nil
}

func parseGameInfoLineForKIF(line string, result *model.ParsedKifu, tmpTimeRule *model.GameInfo) error {
	parts := strings.Split(line, "：")
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	switch key {
	case "開始日時", "対局日":
		t, err := time.Parse("2006/01/02 15:04:05", value)
		if err != nil {
			t, err = time.Parse("2006/01/02", value)
			if err != nil {
				return err
			}
		}
		result.Kifu.StartedAt = &t
	case "終了日時":
		t, err := time.Parse("2006/01/02 15:04:05", value)
		if err != nil {
			t, err = time.Parse("2006/01/02", value)
			if err != nil {
				return err
			}
		}
		opt := &model.KifuOption{
			KifuID: "",
			Name:   key,
			Value:  t.Format("2006-01-02T15:04"),
		}
		result.Options = append(result.Options, opt)
	case "手合割":
		switch value {
		case "平手":
			result.Kifu.InitialPosition = model.SfenHirate.PSFEN()
		case "香落ち":
			result.Kifu.InitialPosition = model.SfenKyoOchi.PSFEN()
		case "右香落ち":
			result.Kifu.InitialPosition = model.SfenMigiKyoOchi.PSFEN()
		case "角落ち":
			result.Kifu.InitialPosition = model.SfenKakuOchi.PSFEN()
		case "飛車落ち":
			result.Kifu.InitialPosition = model.SfenHishaOchi.PSFEN()
		case "飛香落ち":
			result.Kifu.InitialPosition = model.SfenHiKyoOchi.PSFEN()
		case "二枚落ち":
			result.Kifu.InitialPosition = model.SfenNimaiOchi.PSFEN()
		case "三枚落ち":
			result.Kifu.InitialPosition = model.SfenSanmaiOchi.PSFEN()
		case "四枚落ち":
			result.Kifu.InitialPosition = model.SfenYonmaiOchi.PSFEN()
		case "五枚落ち":
			result.Kifu.InitialPosition = model.SfenGomaiOchi.PSFEN()
		case "左五枚落ち":
			result.Kifu.InitialPosition = model.SfenHidariGomaiOchi.PSFEN()
		case "六枚落ち":
			result.Kifu.InitialPosition = model.SfenRokumaiOchi.PSFEN()
		case "左七枚落ち":
			result.Kifu.InitialPosition = model.SfenHidariNanamaiOchi.PSFEN()
		case "右七枚落ち":
			result.Kifu.InitialPosition = model.SfenMigiNanamaiOchi.PSFEN()
		case "八枚落ち":
			result.Kifu.InitialPosition = model.SfenHachimaiOchi.PSFEN()
		case "十枚落ち":
			result.Kifu.InitialPosition = model.SfenJumaiOchi.PSFEN()
		default: // 不明
			result.Kifu.InitialPosition = model.SfenHirate.PSFEN()
		}
	case "先手", "下手":
		result.Kifu.BlackPlayer = &value
	case "後手", "上手":
		result.Kifu.WhitePlayer = &value
	case "持ち時間":
		s, err := getSecondsStringFromString(value)
		if err != nil {
			return err
		}
		(*tmpTimeRule)["持ち時間"] = s + "秒"
	case "秒読み":
		s, err := getSecondsStringFromString(value)
		if err != nil {
			return err
		}
		(*tmpTimeRule)["秒読み"] = s + "秒"
	case "表題":
		result.Kifu.Title = value
	default:
		// 他の情報は全てKifuOptionとする
		opt := &model.KifuOption{
			KifuID: "",
			Name:   key,
			Value:  value,
		}
		result.Options = append(result.Options, opt)
	}
	return nil
}

func getSecondsStringFromString(s string) (string, error) {
	var hours int64 = 0
	var minutes int64 = 0
	var seconds float64 = 0
	var err error
	if strings.Contains(s, "秒") {
		if seconds, err = strconv.ParseFloat(strings.Split(s, "秒")[0], 64); err != nil {
			return "", err
		}
	}
	if strings.Contains(s, "分") {
		if minutes, err = strconv.ParseInt(strings.Split(s, "分")[0], 10, 64); err != nil {
			return "", err
		}
	}
	if strings.Contains(s, "時間") {
		if hours, err = strconv.ParseInt(strings.Split(s, "時間")[0], 10, 64); err != nil {
			return "", err
		}
	}
	return strconv.FormatFloat(float64(hours*3600)+float64(minutes*60)+seconds, 'f', -1, 64), nil
}

func parseMoveLineForKIF(line string, resultBranch *model.KifuBranchWithMoves, lastPlace *model.PiecePlace) error {
	// まず行を3つのパートに分ける
	//   "1 ７六歩(77) ( 0:16/00:00:16)"           -> ["1", "７六歩(77)", "(0:16/00:00:16)"]
	//   "1 ７六歩(77)   (0:16/00:00:16)"          -> ["1", "７六歩(77)", "(0:16/00:00:16)"]
	//   "   9 同　歩(87)        ( 0:02/00:00:13)" -> ["9", "同　歩(87)", "(0:02/00:00:13)"]
	//   "3 中断 ( 0:03/ 0:00:19)"                 -> ["3", "中断",       "(0:03/0:00:19)"]

	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' '
	})
	if len(parts) < 2 {
		return fmt.Errorf("invalid move format: %s", line)
	}
	numString := parts[0]
	moveString := parts[1]
	timeString := strings.Join(parts[2:], "")
	slog.Debug("parseMoveLine", "numString", numString, "moveString", moveString, "timeString", timeString)

	// 番号の取得
	num, err := strconv.ParseInt(numString, 10, 64)
	if err != nil {
		return err
	}

	// 移動元
	moveFrom := model.PIECE_PLACE_IN_HAND
	if strings.Contains(moveString, "(") {
		splited := strings.Split(moveString, "(")
		fromPlaceString := strings.Trim(splited[1], " ()")
		if len(fromPlaceString) != 2 {
			return fmt.Errorf("invalid from place of move: %s", fromPlaceString)
		}
		fileString := string([]rune(fromPlaceString)[0:1])
		file, err := strconv.ParseInt(fileString, 10, 64)
		if err != nil {
			return err
		}
		rankString := string([]rune(fromPlaceString)[1:2])
		rank, err := strconv.ParseInt(rankString, 10, 64)
		if err != nil {
			return err
		}
		moveFrom = model.NewPiecePlaceFromFileRank(int(file), int(rank))
		moveString = splited[0]
	}

	// エンディングの場合
	if ending, ok := model.EndingNameToEndingTypeKIF[moveString]; ok {
		resultBranch.EndingNumber = auxi.PInt64(int64(num)) // ToDo: 不要では？
		resultBranch.EndingType = &ending
		return nil // ToDo: エンディングの消費時間への対応
	}

	// 移動先
	if strings.HasPrefix(moveString, "同") {
		if num <= 1 {
			return fmt.Errorf("cannot set last place first")
		}
		if strings.HasPrefix(moveString, "同　") {
			moveString, _ = strings.CutPrefix(moveString, "同　")
		} else {
			moveString, _ = strings.CutPrefix(moveString, "同")
		}
	} else {
		moveRunes := []rune(moveString)
		fileChar := string(moveRunes[0])
		rankChar := string(moveRunes[1])
		if !strings.Contains(model.FullWidthFileString, fileChar) || !strings.Contains(model.FullWidthRankString, rankChar) {
			return fmt.Errorf("invalid to place: %s", moveString)
		}

		fileRunes := []rune(model.FullWidthFileString)
		rankRunes := []rune(model.FullWidthRankString)
		var file, rank int
		for i := 0; i < 9; i++ {
			if string(fileRunes[i]) == fileChar {
				file = i + 1
			}
			if string(rankRunes[i]) == rankChar {
				rank = i + 1
			}
		}
		*lastPlace = model.NewPiecePlaceFromFileRank(file, rank)
		moveString, _ = strings.CutPrefix(moveString, string(moveRunes[0:2]))
	}

	// 駒の種類
	piece := model.PIECE_VACANCY
	moveRunes := []rune(moveString)
	if p, ok := model.PieceTypeFromStringKIF[string(moveRunes[0:1])]; ok {
		piece = p
		moveString, _ = strings.CutPrefix(moveString, string(moveRunes[0:1]))
	} else if p, ok := model.PieceTypeFromStringKIF[string(moveRunes[0:2])]; ok {
		piece = p
		moveString, _ = strings.CutPrefix(moveString, string(moveRunes[0:2]))
	} else {
		return fmt.Errorf("invalid piece type string: %s", moveString)
	}
	if strings.HasPrefix(moveString, "打") {
		// moveFromはデフォルト値でOK
		moveString, _ = strings.CutPrefix(moveString, "打")
	}
	if strings.HasPrefix(moveString, "成") {
		piece = piece | model.PIECE_PROMOTE
	}

	// 消費時間の解析
	var timeSpent *int64 = nil
	if timeString != "" {
		spents := strings.Split(strings.Split(strings.Trim(timeString, "()"), "/")[0], ":")
		var seconds int64 = 0
		var minutes int64 = 0
		var hours int64 = 0
		if len(spents) >= 1 {
			seconds, err = strconv.ParseInt(spents[len(spents)-1], 10, 64)
			if err != nil {
				return err
			}
		}
		if len(spents) >= 2 {
			minutes, err = strconv.ParseInt(spents[len(spents)-2], 10, 64)
			if err != nil {
				return err
			}
		}
		if len(spents) >= 3 {
			hours, err = strconv.ParseInt(spents[len(spents)-3], 10, 64)
			if err != nil {
				return err
			}
		}
		seconds += hours*3600 + minutes*60
		timeSpent = &seconds
	}

	// 指し手の追加
	move := &model.KifuMove{
		BranchID:    resultBranch.ID,
		Number:      num,
		Piece:       piece,
		FromPlace:   moveFrom,
		ToPlace:     *lastPlace,
		TimeSpentMs: timeSpent,
	}
	resultBranch.Moves = append(resultBranch.Moves, move)
	// slog.Debug("move appended", "ToPlace", move.ToPlace)

	return nil
}
