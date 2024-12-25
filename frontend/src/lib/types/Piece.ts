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

export const PieceCharOfSFEN = new Map<PieceType, string>([
  [PieceType.FU, 'p'],
  [PieceType.KY, 'l'],
  [PieceType.KE, 'n'],
  [PieceType.GI, 's'],
  [PieceType.KI, 'g'],
  [PieceType.KA, 'b'],
  [PieceType.HI, 'r'],
  [PieceType.OU, 'k'],
]);

export const PieceOrderForSFEN = [
  PieceType.OU,
  PieceType.HI,
  PieceType.KA,
  PieceType.KI,
  PieceType.GI,
  PieceType.KE,
  PieceType.KY,
  PieceType.FU,
];

// ------------------------------------------------------------
export class PiecePlace {
  static IN_HAND = 0xff;
  static FILE_CHARS = ['１', '２', '３', '４', '５', '６', '７', '８', '９'];
  static RANK_CHARS = ['一', '二', '三', '四', '五', '六', '七', '八', '九'];
  val: number;

  constructor(val?: number) {
    this.val = val === undefined ? PiecePlace.IN_HAND : val;
  }
  setRowCol(row: number, col: number): void {
    this.val = (row << 4) | col;
  }

  row(): number {
    return (this.val & 0xf0) >> 4;
  }
  col(): number {
    return this.val & 0x0f;
  }
  fileNum(): number {
    return 9 - this.col();
  }
  rankNum(): number {
    return this.row() + 1;
  }
  fileChar(): string {
    const n = this.fileNum();
    if (n > 0 && n <= 9) {
      return PiecePlace.FILE_CHARS[n - 1];
    }
    return '';
  }
  rankChar(): string {
    const n = this.rankNum();
    if (n > 0 && n <= 9) {
      return PiecePlace.RANK_CHARS[n - 1];
    }
    return '';
  }
  isInHand(): boolean {
    return this.val === PiecePlace.IN_HAND;
  }
}

// ------------------------------------------------------------
export type PieceClickEvent = {
  pieceType: PieceType;
  source:
    | {
        type: 'board';
        row: number;
        col: number;
        isBlack?: boolean;
      }
    | {
        type: 'stand';
        isBlack: boolean;
      }
    | {
        type: 'box';
      };
};

// ------------------------------------------------------------
const straightY = Array.from({ length: 8 }, (_, i) => ({ x: 0, y: -(i + 1) }));
const straightYR = Array.from({ length: 8 }, (_, i) => ({ x: 0, y: i + 1 }));
const straightX = Array.from({ length: 8 }, (_, i) => ({ x: i + 1, y: 0 }));
const straightXR = Array.from({ length: 8 }, (_, i) => ({ x: -(i + 1), y: 0 }));
const straightXY = Array.from({ length: 8 }, (_, i) => ({ x: i + 1, y: -(i + 1) }));
const straightXYR = Array.from({ length: 8 }, (_, i) => ({ x: i + 1, y: i + 1 }));
const straightXRY = Array.from({ length: 8 }, (_, i) => ({ x: -(i + 1), y: -(i + 1) }));
const straightXRYR = Array.from({ length: 8 }, (_, i) => ({ x: -(i + 1), y: i + 1 }));
const movablesFU = [{ x: 0, y: -1 }];
const movablesKY = straightY;
const movablesKE = [
  { x: -1, y: -2 },
  { x: 1, y: -2 },
];
const movablesGI = [
  { x: 1, y: -1 },
  { x: 0, y: -1 },
  { x: -1, y: -1 },
  { x: 1, y: 1 },
  { x: -1, y: 1 },
];
const movablesKI = [
  { x: 1, y: -1 },
  { x: 0, y: -1 },
  { x: -1, y: -1 },
  { x: 1, y: 0 },
  { x: -1, y: 0 },
  { x: 0, y: 1 },
];
const movablesKA = [...straightXY, ...straightXYR, ...straightXRY, ...straightXRYR];
const movablesHI = [...straightX, ...straightXR, ...straightY, ...straightYR];
const movablesUM = [
  ...movablesKA,
  { x: 1, y: 0 },
  { x: -1, y: 0 },
  { x: 0, y: 1 },
  { x: 0, y: -1 },
];
const movablesRY = [
  ...movablesHI,
  { x: 1, y: -1 },
  { x: -1, y: -1 },
  { x: 1, y: 1 },
  { x: -1, y: 1 },
];
const movablesOU = [...movablesKI, { x: 1, y: 1 }, { x: -1, y: 1 }];

export const PieceMovables: Map<PieceType, { x: number; y: number; blockable?: boolean }[]> =
  new Map([
    [PieceType.FU, movablesFU],
    [PieceType.KY, movablesKY],
    [PieceType.KE, movablesKE],
    [PieceType.GI, movablesGI],
    [PieceType.KI, movablesKI],
    [PieceType.KA, movablesKA],
    [PieceType.HI, movablesHI],
    [PieceType.OU, movablesOU],
    [PieceType.TO, movablesKI],
    [PieceType.NY, movablesKI],
    [PieceType.NG, movablesKI],
    [PieceType.UM, movablesUM],
    [PieceType.RY, movablesRY],
  ]);
