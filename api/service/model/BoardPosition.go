// service/model/BoardPosition.go
// 保存用のSFENデータとシミュレーション用のBoardPositionデータの相互変換

package model

import (
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/jcytp/kifup-api/common/auxi"
)

type SFEN string

const sfenHirate SFEN = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"

type BoardPosition struct {
	BlackBoard  [9][9]PieceType
	WhiteBoard  [9][9]PieceType
	BlackHands  map[PieceType]int32
	WhiteHands  map[PieceType]int32
	isBlackTurn bool
}

var pieceTypeOfSFEN = map[rune]PieceType{
	'P': PIECE_FU,
	'L': PIECE_KY,
	'N': PIECE_KE,
	'S': PIECE_GI,
	'G': PIECE_KI,
	'B': PIECE_KA,
	'R': PIECE_HI,
	'K': PIECE_OU,
	'p': PIECE_FU,
	'l': PIECE_KY,
	'n': PIECE_KE,
	's': PIECE_GI,
	'g': PIECE_KI,
	'b': PIECE_KA,
	'r': PIECE_HI,
	'k': PIECE_OU,
	'+': PIECE_PROMOTE,
}

var sfenOfPieceType = map[PieceType]rune{
	PIECE_FU: 'p',
	PIECE_KY: 'l',
	PIECE_KE: 'n',
	PIECE_GI: 's',
	PIECE_KI: 'g',
	PIECE_KA: 'b',
	PIECE_HI: 'r',
	PIECE_OU: 'k',
}

// SFEN -> BoardPosition
func NewBoardPosition(sfen *SFEN) (*BoardPosition, error) {
	if sfen == nil {
		s := sfenHirate
		sfen = &s
	}

	bp := &BoardPosition{
		BlackBoard:  [9][9]PieceType{},
		WhiteBoard:  [9][9]PieceType{},
		BlackHands:  map[PieceType]int32{},
		WhiteHands:  map[PieceType]int32{},
		isBlackTurn: true,
	}

	parts := strings.Split(string(*sfen), " ")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid sfen format: expected 4 parts (position, turn, hands, move count), got %d", len(parts))
	}

	// Parse board position
	boardRows := strings.Split(parts[0], "/")
	if len(boardRows) != 9 {
		return nil, fmt.Errorf("invalid number of rows: expected 9, got %d", len(boardRows))
	}
	for i, row := range boardRows {
		x := 0
		flgPromote := false
		for _, c := range row {
			if x > 8 {
				return nil, fmt.Errorf("row %d exceeds board size: too many positions", i+1)
			}
			if c >= '1' && c <= '9' { // number
				for k := rune(0); k < c-'0'; k++ {
					bp.BlackBoard[i][x] = PIECE_VACANCY
					bp.WhiteBoard[i][x] = PIECE_VACANCY
					x++
				}
				continue
			}
			p, ok := pieceTypeOfSFEN[c]
			if !ok {
				return nil, fmt.Errorf("invalid piece character '%c' at row %d col %d", c, i+1, x+1)
			}
			if p == PIECE_PROMOTE { // promote mark
				if flgPromote {
					return nil, fmt.Errorf("consecutive promotion marks at row %d col %d", i+1, x+1)
				}
				flgPromote = true
				continue
			}
			if flgPromote {
				p = p | PIECE_PROMOTE
				flgPromote = false
			}
			if c >= 'A' && c <= 'Z' { // black piece
				bp.BlackBoard[i][x] = p
				bp.WhiteBoard[i][x] = PIECE_VACANCY
			} else { // white piece
				bp.BlackBoard[i][x] = PIECE_VACANCY
				bp.WhiteBoard[i][x] = p
			}
			x++
		}
		if x != 9 {
			return nil, fmt.Errorf("row %d has incorrect length: expected 9 positions, got %d", i+1, x)
		}
	}

	// Parse turn
	if parts[1] != "b" && parts[1] != "w" {
		return nil, fmt.Errorf("invalid turn indicator: expected 'b' or 'w', got '%s'", parts[1])
	}
	bp.isBlackTurn = parts[1] == "b"

	// Parse hands
	if parts[2] != "-" {
		count := int32(1)
		for i, c := range parts[2] {
			if c >= '1' && c <= '9' { // number
				if i == len(parts[2])-1 {
					return nil, fmt.Errorf("number at end of hands section")
				}
				count = c - '0'
				continue
			}
			p, ok := pieceTypeOfSFEN[c]
			if !ok || p == PIECE_PROMOTE {
				return nil, fmt.Errorf("invalid piece '%c' in hands", c)
			}
			if c >= 'A' && c <= 'Z' { // black piece
				if _, exists := bp.BlackHands[p]; exists {
					return nil, fmt.Errorf("duplicate piece '%c' in black hands", c)
				}
				bp.BlackHands[p] = count
			} else {
				if _, exists := bp.WhiteHands[p]; exists {
					return nil, fmt.Errorf("duplicate piece '%c' in white hands", c)
				}
				bp.WhiteHands[p] = count
			}
			count = int32(1)
		}
	}

	// Validate move number
	if _, err := strconv.Atoi(parts[3]); err != nil {
		return nil, fmt.Errorf("invalid move number: %s", parts[3])
	}

	return bp, nil
}

// BoardPosition -> SFEN
func (bp *BoardPosition) ToSFEN(moveCount int) (SFEN, error) {
	if moveCount < 1 {
		return "", fmt.Errorf("invalid move count: must be positive, got %d", moveCount)
	}

	var sb strings.Builder

	// Write board position
	for i := 0; i < 9; i++ {
		emptyCount := 0
		for j := 0; j < 9; j++ {
			blackPiece := bp.BlackBoard[i][j]
			whitePiece := bp.WhiteBoard[i][j]
			if blackPiece == PIECE_VACANCY && whitePiece == PIECE_VACANCY {
				emptyCount++
				continue
			}
			if emptyCount > 0 { // write vacancy brefore piece
				sb.WriteRune(rune('0' + emptyCount))
				emptyCount = 0
			}
			if whitePiece == PIECE_VACANCY { // black piece exists
				piece := blackPiece & ^PIECE_PROMOTE
				c, ok := sfenOfPieceType[piece]
				if !ok {
					return "", fmt.Errorf("invalid black piece at row %d col %d", i, j)
				}
				if blackPiece&PIECE_PROMOTE != 0 {
					sb.WriteRune('+')
				}
				sb.WriteRune(unicode.ToUpper(c))
			} else if blackPiece == PIECE_VACANCY { // white piece exists
				piece := whitePiece & ^PIECE_PROMOTE
				c, ok := sfenOfPieceType[piece]
				if !ok {
					return "", fmt.Errorf("invalid white piece at row %d col %d", i, j)
				}
				if whitePiece&PIECE_PROMOTE != 0 {
					sb.WriteRune('+')
				}
				sb.WriteRune(c)
			} else {
				return "", fmt.Errorf("both black and white piece exists at row %d col %d", i, j)
			}
		}
		if emptyCount > 0 { // write vacancy at end of the row
			sb.WriteRune(rune('0' + emptyCount))
		}
		if i < 8 {
			sb.WriteRune('/') // row separator
		}
	}

	// Write turn
	sb.WriteRune(' ')
	if bp.isBlackTurn {
		sb.WriteRune('b')
	} else {
		sb.WriteRune('w')
	}

	// Write hands
	sb.WriteRune(' ')
	if len(bp.BlackHands) == 0 && len(bp.WhiteHands) == 0 {
		sb.WriteRune('-')
	} else {
		// Prepare sorted piece types: krbgsnlp
		pieceTypes := make([]PieceType, 0, len(sfenOfPieceType))
		for p := range sfenOfPieceType {
			pieceTypes = append(pieceTypes, p)
		}
		sort.Slice(pieceTypes, func(i, j int) bool {
			return pieceTypes[i] > pieceTypes[j]
		})
		// Black hands
		for _, p := range pieceTypes {
			if count, ok := bp.BlackHands[p]; ok && count > 0 {
				if count > 1 {
					sb.WriteString(strconv.Itoa(int(count)))
				}
				sb.WriteRune(unicode.ToUpper(sfenOfPieceType[p]))
			}
		}
		// White hands
		for _, p := range pieceTypes {
			if count, ok := bp.WhiteHands[p]; ok && count > 0 {
				if count > 1 {
					sb.WriteString(strconv.Itoa(int(count)))
				}
				sb.WriteRune(sfenOfPieceType[p])
			}
		}
	}

	// Write move count
	sb.WriteRune(' ')
	sb.WriteString(strconv.Itoa(moveCount))

	return SFEN(sb.String()), nil
}

func (bp *BoardPosition) Move(move *KifuMove) error {
	var currentBoard, opponentBoard *[9][9]PieceType
	var currentHand *map[PieceType]int32
	if bp.isBlackTurn {
		currentBoard = &bp.BlackBoard
		opponentBoard = &bp.WhiteBoard
		currentHand = &bp.BlackHands
	} else {
		currentBoard = &bp.WhiteBoard
		opponentBoard = &bp.BlackBoard
		currentHand = &bp.WhiteHands
	}

	// 指し手のバリデーションチェック
	fromRow, fromCol := move.FromPlace.RowCol()
	toRow, toCol := move.ToPlace.RowCol()
	if move.FromPlace == PIECE_PLACE_IN_HAND {
		cnt, ok := (*currentHand)[move.Piece]
		if !(ok && cnt > 0) {
			return fmt.Errorf("failed to simulate move. no piece in hand")
		}
	} else {
		if fromRow < 0 || fromRow >= 9 || fromCol < 0 || fromCol >= 9 {
			return fmt.Errorf("invalid from position: row=%d, col=%d", fromRow, fromCol)
		}
		fromPiece := currentBoard[fromRow][fromCol]
		if fromPiece == PIECE_VACANCY {
			return fmt.Errorf("no piece at from position: row=%d, col=%d", fromRow, fromCol)
		}
		if fromPiece&PIECE_PROMOTE == 0 && move.Piece&PIECE_PROMOTE > 0 {
			fromPiece = fromPiece | PIECE_PROMOTE
		}
		if fromPiece != move.Piece {
			return fmt.Errorf("piece type mismatch: expected=%d, actual=%d", move.Piece, fromPiece)
		}
	}
	if toRow < 0 || toRow >= 9 || toCol < 0 || toCol >= 9 {
		return fmt.Errorf("invalid to position: row=%d, col=%d", toRow, toCol)
	}
	if currentBoard[toRow][toCol] != PIECE_VACANCY {
		return fmt.Errorf("destination is occupied by own piece: row=%d, col=%d", toRow, toCol)
	}

	// 移動先の敵駒を取る
	target := opponentBoard[toRow][toCol]
	if target != PIECE_VACANCY {
		cnt, ok := (*currentHand)[target & ^PIECE_PROMOTE]
		if !ok {
			cnt = 0
		}
		(*currentHand)[target & ^PIECE_PROMOTE] = cnt + 1 // 持ち駒を追加
		opponentBoard[toRow][toCol] = PIECE_VACANCY       // 移動先の駒を削除
	}

	// 自駒を移動する
	currentBoard[toRow][toCol] = move.Piece // 移動先に配置
	if move.FromPlace == PIECE_PLACE_IN_HAND {
		cnt, ok := (*currentHand)[move.Piece]
		if ok && cnt > 0 {
			(*currentHand)[move.Piece] = cnt - 1 // 持ち駒を消費
		}
	} else {
		currentBoard[fromRow][fromCol] = PIECE_VACANCY // 移動元を削除
	}

	// 手番を変更
	bp.isBlackTurn = !bp.isBlackTurn

	return nil
}

func (bp *BoardPosition) Copy() *BoardPosition {
	newPosition := &BoardPosition{
		BlackBoard:  bp.BlackBoard,
		WhiteBoard:  bp.WhiteBoard,
		BlackHands:  make(map[PieceType]int32, len(bp.BlackHands)),
		WhiteHands:  make(map[PieceType]int32, len(bp.WhiteHands)),
		isBlackTurn: bp.isBlackTurn,
	}
	for k, v := range bp.BlackHands {
		newPosition.BlackHands[k] = v
	}
	for k, v := range bp.WhiteHands {
		newPosition.WhiteHands[k] = v
	}
	return newPosition
}

func (bp *BoardPosition) IsPromote(move *KifuMove) *bool {
	if move.Piece&PIECE_PROMOTE > 0 { // 動いた結果が成駒
		if move.FromPlace == PIECE_PLACE_IN_HAND { // 持ち駒から打った
			slog.Error("failed to simulate board", "number", move.Number)
			return nil
		} else {
			var pieceBefore PieceType
			row, col := move.FromPlace.RowCol()
			if bp.isBlackTurn { // 先手番
				pieceBefore = bp.BlackBoard[row][col]
			} else { // 後手番
				pieceBefore = bp.WhiteBoard[row][col]
			}
			if pieceBefore == PIECE_VACANCY {
				slog.Error("failed to simulate board", "number", move.Number)
				return nil
			}
			if pieceBefore&PIECE_PROMOTE == 0 { // 元の駒が成駒ではない
				return auxi.PBool(true) // 成判定true
			}
			// 成駒が動いた場合は、成判定null
		}
	} else { // 動いた結果が成駒ではない
		if move.FromPlace != PIECE_PLACE_IN_HAND { // 持ち駒から打ったのではない
			// 敵陣から動いたか
			fromRow, _ := move.FromPlace.RowCol()
			if (bp.isBlackTurn && fromRow < 3) || (!bp.isBlackTurn && fromRow > 5) {
				return auxi.PBool(false) // 成判定false（不成）
			}
			// 敵陣に入ったか
			toRow, _ := move.ToPlace.RowCol()
			if (bp.isBlackTurn && toRow < 3) || (!bp.isBlackTurn && toRow > 5) {
				return auxi.PBool(false) // 成判定false（不成）
			}
			// その他は、成判定null
		}
	}
	return nil
}

func (bp *BoardPosition) CatchPiece(move *KifuMove) *PieceType {
	var result PieceType = PIECE_VACANCY
	row, col := move.ToPlace.RowCol()
	if bp.isBlackTurn {
		result = bp.WhiteBoard[row][col]
	} else {
		result = bp.BlackBoard[row][col]
	}
	if result != PIECE_VACANCY {
		return &result
	}
	return nil // 駒を取ってない場合はnull
}

func (bp *BoardPosition) DirectionSign(move *KifuMove) *string {
	// ToDo: ki2形式で出力するにはDirecrtionSignが必要
	return nil // 必要ない場合はnull
}
