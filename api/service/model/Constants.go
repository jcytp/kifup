// service/model/Constants.go

package model

import "github.com/jcytp/kifup-api/common/auxi"

// --------------------------------------------------------------------------------
type PieceType int64

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

var PieceTypeName = map[PieceType]string{
	PIECE_FU:      "FU",
	PIECE_KY:      "KY",
	PIECE_KE:      "KE",
	PIECE_GI:      "GI",
	PIECE_KI:      "KI",
	PIECE_KA:      "KA",
	PIECE_HI:      "HI",
	PIECE_OU:      "OU",
	PIECE_TO:      "TO",
	PIECE_NY:      "NY",
	PIECE_NK:      "NK",
	PIECE_NG:      "NG",
	PIECE_UM:      "UM",
	PIECE_RY:      "RY",
	PIECE_VACANCY: "--",
}

// 駒種ごとの数
var AllPieceCount = map[PieceType]int{
	PIECE_FU: 18, // 歩：9x2
	PIECE_KY: 4,  // 香：2x2
	PIECE_KE: 4,  // 桂：2x2
	PIECE_GI: 4,  // 銀：2x2
	PIECE_KI: 4,  // 金：2x2
	PIECE_KA: 2,  // 角：1x2
	PIECE_HI: 2,  // 飛：1x2
	PIECE_OU: 2,  // 王：1x2（玉将と王将）
}

// 平手の初期配置（rank順 -> file順）
var InitialPositionHirate = []*KifuInitialPiece{
	{"", PIECE_KY, false, auxi.PInt64(1), auxi.PInt64(1)},
	{"", PIECE_KE, false, auxi.PInt64(2), auxi.PInt64(1)},
	{"", PIECE_GI, false, auxi.PInt64(3), auxi.PInt64(1)},
	{"", PIECE_KI, false, auxi.PInt64(4), auxi.PInt64(1)},
	{"", PIECE_OU, false, auxi.PInt64(5), auxi.PInt64(1)},
	{"", PIECE_KI, false, auxi.PInt64(6), auxi.PInt64(1)},
	{"", PIECE_GI, false, auxi.PInt64(7), auxi.PInt64(1)},
	{"", PIECE_KE, false, auxi.PInt64(8), auxi.PInt64(1)},
	{"", PIECE_KY, false, auxi.PInt64(9), auxi.PInt64(1)},
	{"", PIECE_KA, false, auxi.PInt64(2), auxi.PInt64(2)},
	{"", PIECE_HI, false, auxi.PInt64(8), auxi.PInt64(2)},
	{"", PIECE_FU, false, auxi.PInt64(1), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(2), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(3), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(4), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(5), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(6), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(7), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(8), auxi.PInt64(3)},
	{"", PIECE_FU, false, auxi.PInt64(9), auxi.PInt64(3)},
	{"", PIECE_FU, true, auxi.PInt64(1), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(2), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(3), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(4), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(5), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(6), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(7), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(8), auxi.PInt64(7)},
	{"", PIECE_FU, true, auxi.PInt64(9), auxi.PInt64(7)},
	{"", PIECE_HI, true, auxi.PInt64(2), auxi.PInt64(8)},
	{"", PIECE_KA, true, auxi.PInt64(8), auxi.PInt64(8)},
	{"", PIECE_KY, true, auxi.PInt64(1), auxi.PInt64(9)},
	{"", PIECE_KE, true, auxi.PInt64(2), auxi.PInt64(9)},
	{"", PIECE_GI, true, auxi.PInt64(3), auxi.PInt64(9)},
	{"", PIECE_KI, true, auxi.PInt64(4), auxi.PInt64(9)},
	{"", PIECE_OU, true, auxi.PInt64(5), auxi.PInt64(9)},
	{"", PIECE_KI, true, auxi.PInt64(6), auxi.PInt64(9)},
	{"", PIECE_GI, true, auxi.PInt64(7), auxi.PInt64(9)},
	{"", PIECE_KE, true, auxi.PInt64(8), auxi.PInt64(9)},
	{"", PIECE_KY, true, auxi.PInt64(9), auxi.PInt64(9)},
}

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

// --------------------------------------------------------------------------------

// // 盤上ビットボード： 81bit -> 11byte
// // 駒種ごとのビットボード：8 * 11byte -> 88byte
// // 成駒のビットボード：11byte
// // 後手のビットボード：11byte
// // 盤上の表現：110byte
// // 190 256
// // 15   16
// // 6     8

// // 63byte

// type BitBoard [12]byte

// // 盤上：     9x9　→81byte
// // 先手駒台： FU:2byte + KY,KE,GI,KI,KA,HI:1byte -> 8byte
// // 後手駒台： FU:2byte + KY,KE,GI,KI,KA,HI:1byte -> 8byte
// // 局面：     81 + 8 + 8 -> 97byte

// type BoardPosition [97]byte

// func (t *BoardPosition) PieceOnBoard(file int, rank int) PieceType {
// 	return PieceType((*t)[(rank-1)*9+(file-1)])
// }

// func (t *BoardPosition) PieceCountInHands(isBlack bool, piece PieceType) int {
// 	offset := 81
// 	if !isBlack {
// 		offset += 8
// 	}
// 	if piece == PIECE_FU {
// 		return int((*t)[offset])*0xF + int((*t)[offset+1])
// 	} else {
// 		return int((*t)[offset+1+int(piece)])
// 	}
// }

// func (t *BoardPosition) SetPieceOnBoard(file int, rank int, piece PieceType) {
// 	(*t)[(rank-1)*9+(file-1)] = byte(piece)
// }

// var INITIAL_POSITION BoardPosition

// func init() {
// 	b := &BoardPosition{}
// 	for i := 0; i < 81; i++ {
// 		b[i] = byte(PIECE_VACANCY)
// 	}
// 	for i := 81; i < 97; i++ {
// 		b[i] = byte(0)
// 	}
// 	b.SetPieceOnBoard(1, 1, PIECE_KY)
// 	b.SetPieceOnBoard(2, 1, PIECE_KE)
// 	b.SetPieceOnBoard(3, 1, PIECE_GI)
// 	b.SetPieceOnBoard(4, 1, PIECE_KI)
// 	b.SetPieceOnBoard(5, 1, PIECE_OU)
// }
