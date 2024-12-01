// service/model/Piece.go

package model

// --------------------------------------------------------------------------------
type PieceType byte

const (
	PIECE_PROMOTE PieceType = 0x8
	PIECE_FU      PieceType = 0x0
	PIECE_KY      PieceType = 0x1
	PIECE_KE      PieceType = 0x2
	PIECE_GI      PieceType = 0x3
	PIECE_KI      PieceType = 0x4
	PIECE_KA      PieceType = 0x5
	PIECE_HI      PieceType = 0x6
	PIECE_OU      PieceType = 0x7
	PIECE_TO      PieceType = PIECE_FU | PIECE_PROMOTE
	PIECE_NY      PieceType = PIECE_KY | PIECE_PROMOTE
	PIECE_NK      PieceType = PIECE_KE | PIECE_PROMOTE
	PIECE_NG      PieceType = PIECE_GI | PIECE_PROMOTE
	PIECE_UM      PieceType = PIECE_KA | PIECE_PROMOTE
	PIECE_RY      PieceType = PIECE_HI | PIECE_PROMOTE
	PIECE_VACANCY PieceType = 0xF
)

var PieceTypeNameCSA = map[PieceType]string{
	PIECE_FU: "FU",
	PIECE_KY: "KY",
	PIECE_KE: "KE",
	PIECE_GI: "GI",
	PIECE_KI: "KI",
	PIECE_KA: "KA",
	PIECE_HI: "HI",
	PIECE_OU: "OU",
	PIECE_TO: "TO",
	PIECE_NY: "NY",
	PIECE_NK: "NK",
	PIECE_NG: "NG",
	PIECE_UM: "UM",
	PIECE_RY: "RY",
}

var PieceTypeNameKIF = map[PieceType]string{
	PIECE_FU: "歩",
	PIECE_KY: "香",
	PIECE_KE: "桂",
	PIECE_GI: "銀",
	PIECE_KI: "金",
	PIECE_KA: "角",
	PIECE_HI: "飛",
	PIECE_OU: "王",
	PIECE_TO: "と",
	PIECE_NY: "成香",
	PIECE_NK: "成桂",
	PIECE_NG: "成銀",
	PIECE_UM: "馬",
	PIECE_RY: "龍",
}

const (
	DIRECTION_SIGN_GO_BACK    = "引"
	DIRECTION_SIGN_GO_LATERAL = "寄"
	DIRECTION_SIGN_GO_FORWARD = "上"

	DIRECTION_SIGN_FROM_LEFT   = "左"
	DIRECTION_SIGN_FROM_CENTER = "直"
	DIRECTION_SIGN_FROM_RIGHT  = "右"
)

// --------------------------------------------------------------------------------
type PiecePlace byte

// １一 ＝ 1筋1段 ＝ 8列0行 ＝ 0x81
// 持ち駒 = 0x99
const PIECE_PLACE_IN_HAND = 0xFF

func (place PiecePlace) RowCol() (row int, col int) {
	return int(place & 0xF0 >> 4), int(place & 0x0F)
}

func (place PiecePlace) FileRank() (file int, rank int) {
	row, col := place.RowCol()
	return 9 - col, row + 1
}
