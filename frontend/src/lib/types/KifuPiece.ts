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

export const PIECE_PLACE_IN_HAND = 0xff;
