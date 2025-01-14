// service/api/parser/csa.go

package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jcytp/kifup-api/common/auxi"
	"github.com/jcytp/kifup-api/service/model"
)

func ParseFromCSA(content string) (*model.ParsedKifu, error) {
	lines := strings.Split(content, "\n")
	if !checkKifuFormatCSA(lines) {
		return nil, nil
	}

	kifuID := auxi.NewULID()       // dummy id
	mainBranchID := auxi.NewULID() // dummy id
	result := &model.ParsedKifu{
		Kifu: &model.Kifu{
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

	position, err := model.NewBoardPosition(model.SfenAllInBox.PSFEN())
	if err != nil {
		return nil, err
	}

	flgAL := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "'") && !strings.HasPrefix(line, "'*") {
			continue // 空行とコメントをスキップ
		}

		statements := strings.Split(line, ",")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "'") && !strings.HasPrefix(stmt, "'*") {
				continue // 空行とコメントをスキップ
			}

			switch {
			case strings.HasPrefix(stmt, "V"): // CSAバージョン情報（処理なし）
			case strings.HasPrefix(stmt, "N+"): // 先手
				result.Kifu.BlackPlayer = auxi.PString(strings.TrimPrefix(stmt, "N+"))
			case strings.HasPrefix(stmt, "N-"): // 後手
				result.Kifu.WhitePlayer = auxi.PString(strings.TrimPrefix(stmt, "N-"))
			case strings.HasPrefix(stmt, "$"): // 棋譜情報
				if err := parseGameInfoLineForCSA(stmt, result); err != nil {
					return nil, err
				}
			case strings.HasPrefix(stmt, "P"): // 局面情報
				if err := parsePositionLineForCSA(stmt, position, &flgAL); err != nil {
					return nil, err
				}
			case stmt == "+" || stmt == "-": // 手番情報
				position.IsBlackTurn = (stmt == "+")
				sfen, err := position.ToSFEN(1)
				if err != nil {
					return nil, err
				}
				result.Kifu.InitialPosition = &sfen
			case strings.HasPrefix(stmt, "+") || strings.HasPrefix(stmt, "-"): // 指し手
				if err := parseMoveLineForCSA(line, result.Branches[0], position.IsBlackTurn); err != nil {
					return nil, err
				}
			case strings.HasPrefix(stmt, "%"): // エンディング
				moveString, _ := strings.CutPrefix(stmt, "%")
				if ending, ok := model.EndingNameToEndingTypeCSA[moveString]; ok {
					nextNumber := len(result.Branches[0].Moves) + 1
					result.Branches[0].EndingNumber = auxi.PInt64(int64(nextNumber)) // ToDo: 不要では？
					result.Branches[0].EndingType = &ending
				}
			case strings.HasPrefix(stmt, "T"): // 消費時間
				num := len(result.Branches[0].Moves) - 1
				if num < 0 {
					return nil, fmt.Errorf("spent time before any move")
				}
				timeString, _ := strings.CutPrefix(stmt, "T")
				seconds, err := strconv.ParseFloat(timeString, 64)
				if err != nil {
					return nil, err
				}
				result.Branches[0].Moves[num].TimeSpentMs = auxi.PInt64(int64(seconds * 1000))
			case strings.HasPrefix(stmt, "'*"): // プログラムが読むコメント -> 局面コメント
				num := len(result.Branches[0].Moves) - 1
				if num < 0 {
					break
				}
				comment := ""
				if result.Branches[0].Moves[num].Comment != nil {
					comment = *result.Branches[0].Moves[num].Comment + "\n"
				}
				additional, _ := strings.CutPrefix(stmt, "'")
				comment += additional
				result.Branches[0].Moves[num].Comment = &comment
			}
		}
	}

	return result, nil
}

func checkKifuFormatCSA(lines []string) bool {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		statements := strings.Split(line, ",")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			// 手番の指定行があればフォーマット適合とする
			if stmt == "+" || stmt == "-" {
				return true
			}
		}
	}
	return false
}

func parseGameInfoLineForCSA(line string, result *model.ParsedKifu) error {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid game info format of csa: %s", line)
	}

	key, _ := strings.CutPrefix(strings.TrimSpace(parts[0]), "$")
	value := strings.TrimSpace(parts[1])

	mapKeyToName := map[string]string{
		"EVENT":     "棋戦",
		"SITE":      "対局場所",
		"OPENING":   "戦型",
		"MAX_MOVES": "最大手数",
	}

	switch key {
	case "START_TIME":
		if t, err := time.Parse("2006/01/02 15:04:05", value); err == nil {
			result.Kifu.StartedAt = &t
		} else if t, err := time.Parse("2006/01/02", value); err == nil {
			result.Kifu.StartedAt = &t
		}
	case "END_TIME":
		var endTime time.Time
		if t, err := time.Parse("2006/01/02 15:04:05", value); err == nil {
			endTime = t
		} else if t, err := time.Parse("2006/01/02", value); err == nil {
			endTime = t
		} else {
			return err
		}
		opt := &model.KifuOption{
			Name:  key,
			Value: endTime.Format("2006-01-02T15:04"),
		}
		result.Options = append(result.Options, opt)
	case "TIME_LIMIT":
		if len(parts) < 3 {
			return fmt.Errorf("invalid time_limit in csa")
		}
		hours, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		splited := strings.Split(parts[2], "+")
		minutes, err := strconv.ParseInt(splited[0], 10, 64)
		if err != nil {
			return err
		}
		byoyomi, err := strconv.ParseInt(splited[1], 10, 64)
		if err != nil {
			return err
		}
		timeGameInfo := model.GameInfo{
			"持ち時間": fmt.Sprintf("%d秒", hours*3600+minutes*60),
			"秒読み":  fmt.Sprintf("%d秒", byoyomi),
		}
		result.Kifu.TimeRule = timeGameInfo.GetTimeRule()
	case "TIME":
		times := strings.Split(value, "+")
		timeGameInfo := model.GameInfo{}
		if len(times) >= 1 && times[0] != "0" {
			timeGameInfo["持ち時間"] = fmt.Sprintf("%s秒", times[0])
		}
		if len(times) >= 2 && times[1] != "0" {
			timeGameInfo["秒読み"] = fmt.Sprintf("%s秒", times[1])
		}
		if len(times) >= 3 && times[2] != "0" {
			timeGameInfo["秒加算"] = fmt.Sprintf("%s秒", times[2])
		}
		result.Kifu.TimeRule = timeGameInfo.GetTimeRule()
	case "JISHOGI":
		if value == "24" {
			result.Options = append(result.Options, &model.KifuOption{
				Name:  "持将棋ルール",
				Value: "24点法",
			})
		} else if value == "27" {
			result.Options = append(result.Options, &model.KifuOption{
				Name:  "持将棋ルール",
				Value: "27点法",
			})
		}
	case "NOTE":
		note := strings.ReplaceAll(value, "\\n", "\n")
		note = strings.ReplaceAll(note, "\\\\", "\\")
		result.Options = append(result.Options, &model.KifuOption{
			Name:  "備考",
			Value: note,
		})
	default:
		name, ok := mapKeyToName[key]
		if ok {
			result.Options = append(result.Options, &model.KifuOption{
				Name:  name,
				Value: value,
			})
		} else {
			result.Options = append(result.Options, &model.KifuOption{
				Name:  key,
				Value: value,
			})
		}
	}
	return nil
}

func parsePositionLineForCSA(line string, position *model.BoardPosition, flgAL *bool) error {
	if strings.HasPrefix(line, "PI") {
		if tmpPosition, err := model.NewBoardPosition(model.SfenHirate.PSFEN()); err != nil {
			return err
		} else {
			*position = *tmpPosition
		}

		line = strings.TrimPrefix(line, "PI")
		for i := 0; i*4 < len(line); i++ {
			if i*4+4 > len(line) {
				return fmt.Errorf("invalid piece drop format: %s", line)
			}
			pieceString := line[i*4 : i*4+4]
			file := int(pieceString[0] - '0')
			rank := int(pieceString[1] - '0')
			place := model.NewPiecePlaceFromFileRank(file, rank)
			row, col := place.RowCol()
			piece, ok := model.PieceTypeFromStringCSA[pieceString[2:]]
			if !ok {
				return fmt.Errorf("invalid piece drop format (invalid piece type): %s", line)
			}
			if position.WhiteBoard[row][col] != piece && position.BlackBoard[row][col] != piece {
				return fmt.Errorf("invalid piece drop format (piece not unmatched): %s", line)
			}
			position.WhiteBoard[row][col] = model.PIECE_VACANCY
			position.BlackBoard[row][col] = model.PIECE_VACANCY
		}
	} else if len(line) >= 2 && line[1] >= '1' && line[1] <= '9' {
		rank := int(line[1] - '0')
		row, _ := model.NewPiecePlaceFromFileRank(1, rank).RowCol()
		line = line[2:]
		if strings.HasSuffix(line, "*") {
			line += " " // TrimSpaceで削除された空白を付与する
		}
		for col := 0; col < 9; col++ {
			if col*3+3 > len(line) {
				return fmt.Errorf("invalid position format (rank set): %s", line)
			}
			pieceString := line[col*3 : col*3+3]
			if pieceString == " * " {
				// 初期状態は全て駒箱なので何もしない
			} else {
				if pieceString[0] != '+' && pieceString[0] != '-' {
					return fmt.Errorf("invalid piece string format (mark): %s", pieceString)
				}
				isBlack := pieceString[0] == '+'
				piece, ok := model.PieceTypeFromStringCSA[pieceString[1:]]
				if !ok {
					return fmt.Errorf("invalid piece string format (piece): %s", pieceString)
				}
				if isBlack {
					position.BlackBoard[row][col] = piece
				} else {
					position.WhiteBoard[row][col] = piece
				}
			}
		}
	} else if strings.HasPrefix(line, "P+") || strings.HasPrefix(line, "P-") {
		isBlack := line[1] == '+'
		line = line[2:]
		for i := 0; i*4 < len(line); i++ {
			if i*4+4 > len(line) {
				return fmt.Errorf("invalid position format: %s", line)
			}
			pieceString := line[i*4 : i*4+4]
			if pieceString[:2] == "00" { // 持ち駒
				if pieceString[2:] == "AL" {
					if *flgAL {
						return fmt.Errorf("invalid position format: multiple AL")
					}
					*flgAL = true

					inbox := position.AllPiecesInBox()
					for piece, cnt := range inbox {
						if piece == model.PIECE_OU {
							continue
						}
						if isBlack {
							position.BlackHands[piece] += cnt
						} else {
							position.WhiteHands[piece] += cnt
						}
					}
				} else {
					piece, ok := model.PieceTypeFromStringCSA[pieceString[2:]]
					if !ok {
						return fmt.Errorf("invalid piece type: %s", pieceString[2:])
					}
					if isBlack {
						position.BlackHands[piece]++
					} else {
						position.WhiteHands[piece]++
					}
				}
			} else { // 盤上
				file := int(pieceString[0] - '0')
				rank := int(pieceString[1] - '0')
				row, col := model.NewPiecePlaceFromFileRank(file, rank).RowCol()
				piece, ok := model.PieceTypeFromStringCSA[pieceString[2:]]
				if !ok {
					return fmt.Errorf("invalid piece string format (piece): %s", pieceString)
				}
				if position.BlackBoard[row][col] != model.PIECE_VACANCY ||
					position.WhiteBoard[row][col] != model.PIECE_VACANCY {
					return fmt.Errorf("piece already exists at %d%d", file, rank)
				}
				if isBlack {
					position.BlackBoard[row][col] = piece
				} else {
					position.WhiteBoard[row][col] = piece
				}
			}
		}
	}

	return nil
}

func parseMoveLineForCSA(line string, branch *model.KifuBranchWithMoves, isBlackFirst bool) error {
	if len(line) != 7 { // ex: "+7776FU"
		return fmt.Errorf("invalid move format: unmatch length: %s", line)
	}

	isBlackMove := line[0] == '+'
	number := int64(len(branch.Moves) + 1)
	if (number%2 == 1 && isBlackMove != isBlackFirst) || (number%2 == 0 && isBlackMove == isBlackFirst) {
		return fmt.Errorf("invalid turn sequence at move %d", number)
	}

	fromFile := int(line[1] - '0')
	fromRank := int(line[2] - '0')
	fromPlace := model.PIECE_PLACE_IN_HAND
	if fromFile == 0 && fromRank == 0 {
		// OK
	} else if fromFile >= 1 && fromFile <= 9 && fromRank >= 1 && fromRank <= 9 {
		fromPlace = model.NewPiecePlaceFromFileRank(fromFile, fromRank)
	} else {
		return fmt.Errorf("invalid from position: %d%d", fromFile, fromRank)
	}

	toFile := int(line[3] - '0')
	toRank := int(line[4] - '0')
	if toFile < 1 || toFile > 9 || toRank < 1 || toRank > 9 {
		return fmt.Errorf("invalid to position: %d%d", toFile, toRank)
	}
	toPlace := model.NewPiecePlaceFromFileRank(toFile, toRank)

	piece, ok := model.PieceTypeFromStringCSA[line[5:]]
	if !ok {
		return fmt.Errorf("invalid piece type: %s", line[5:])
	}

	move := &model.KifuMove{
		BranchID:  branch.ID,
		Number:    number,
		Piece:     piece,
		FromPlace: fromPlace,
		ToPlace:   toPlace,
	}
	branch.Moves = append(branch.Moves, move)
	return nil
}
