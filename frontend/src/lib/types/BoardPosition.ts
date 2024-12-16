// src/lib/types/BoardPosition.ts

import { PieceType, PieceTypeOfSFEN } from './Piece';

const SFEN_HIRATE = 'lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1';

export class BoardPosition {
  blackBoard: PieceType[][];
  whiteBoard: PieceType[][];
  blackHands: Map<PieceType, number>;
  whiteHands: Map<PieceType, number>;
  pieceBox: Map<PieceType, number>;
  isBlackTurn: boolean;

  constructor(sfen?: string) {
    // 初期化
    this.blackBoard = Array(9)
      .fill(null)
      .map(() => Array(9).fill(PieceType.VACANCY));
    this.whiteBoard = Array(9)
      .fill(null)
      .map(() => Array(9).fill(PieceType.VACANCY));
    this.blackHands = new Map<PieceType, number>();
    this.whiteHands = new Map<PieceType, number>();
    this.pieceBox = new Map<PieceType, number>();
    this.isBlackTurn = true;

    sfen = sfen || SFEN_HIRATE;
    console.debug(sfen);
    const parts = sfen.split(' ');
    if (parts.length < 4) {
      console.error(`Invalid SFEN format: expected 4 parts, got ${parts.length}`);
      return;
    }

    // 盤面の解析
    const rows = parts[0].split('/');
    if (rows.length !== 9) {
      console.error(`Invalid number of rows: expected 9, got ${rows.length}`);
      return;
    }
    for (let i = 0; i < 9; i++) {
      let x = 0;
      let isPromote = false;
      for (const c of rows[i]) {
        if (x > 8) {
          console.error(`Row ${i + 1} exceeds board size`);
          return;
        }
        if (/[1-9]/.test(c)) {
          const emptyCount = parseInt(c);
          x += emptyCount;
          continue;
        }
        let pieceType = PieceTypeOfSFEN.get(c);
        if (pieceType === undefined) {
          console.error(`Invalid piece character '${c}' at row ${i + 1} col ${x + 1}`);
          return;
        }
        if (pieceType == PieceType.PROMOTE) {
          if (isPromote) {
            console.error(`Consecutive promotion marks at row ${i + 1} col ${x + 1}`);
            return;
          }
          isPromote = true;
          continue;
        }
        if (isPromote) {
          pieceType = pieceType | PieceType.PROMOTE;
          isPromote = false;
        }
        if (c >= 'A' && c <= 'Z') {
          this.blackBoard[i][x] = pieceType;
        } else {
          this.whiteBoard[i][x] = pieceType;
        }
        x++;
      }
      if (x != 9) {
        console.error(`Row ${i + 1} has incorrect length: expected 9 positions, got ${x}`);
        return;
      }
    }

    // 手番の解析
  }
}
