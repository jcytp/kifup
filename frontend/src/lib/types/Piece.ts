// src/lib/types/Piece.ts

export const enum PieceType {
  PROMOTE = 0x8,
  FU = 0x0,
  KY = 0x1,
  KE = 0x2,
  GI = 0x3,
  KI = 0x4,
  KA = 0x5,
  HI = 0x6,
  OU = 0x7,
  TO = 0x8, // FU | PROMOTE
  NY = 0x9, // KY | PROMOTE
  NK = 0xa, // KE | PROMOTE
  NG = 0xb, // GI | PROMOTE
  UM = 0xd, // KA | PROMOTE
  RY = 0xe, // HI | PROMOTE
  VACANCY = 0xf,
}

export const PieceChar = new Map<PieceType, string>([
  [PieceType.FU, '歩'],
  [PieceType.KY, '香'],
  [PieceType.KE, '桂'],
  [PieceType.GI, '銀'],
  [PieceType.KI, '金'],
  [PieceType.KA, '角'],
  [PieceType.HI, '飛'],
  [PieceType.OU, '玉'],
  [PieceType.TO, 'と'],
  [PieceType.NY, '杏'],
  [PieceType.NK, '圭'],
  [PieceType.NG, '全'],
  [PieceType.UM, '馬'],
  [PieceType.RY, '龍'],
]);

export const PieceTypeOfSFEN = new Map<string, PieceType>([
  ['P', PieceType.FU],
  ['L', PieceType.KY],
  ['N', PieceType.KE],
  ['S', PieceType.GI],
  ['G', PieceType.KI],
  ['B', PieceType.KA],
  ['R', PieceType.HI],
  ['K', PieceType.OU],
  ['p', PieceType.FU],
  ['l', PieceType.KY],
  ['n', PieceType.KE],
  ['s', PieceType.GI],
  ['g', PieceType.KI],
  ['b', PieceType.KA],
  ['r', PieceType.HI],
  ['k', PieceType.OU],
  ['+', PieceType.PROMOTE],
]);

// ------------------------------------------------------------
export const PIECE_PLACE_IN_HAND = 0xff;
