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
    this.pieceBox = new Map<PieceType, number>([
      [PieceType.FU, 18],
      [PieceType.KY, 4],
      [PieceType.KE, 4],
      [PieceType.GI, 4],
      [PieceType.KI, 4],
      [PieceType.KA, 2],
      [PieceType.HI, 2],
      [PieceType.OU, 2],
    ]);
    this.isBlackTurn = true;

    sfen = sfen || SFEN_HIRATE;
    console.debug(sfen);
    const parts = sfen.split(' ');
    if (parts.length < 4) {
      console.error(`Invalid SFEN format: expected 4 parts, got ${parts.length}`);
      return;
    }

    // 駒箱から取り出す
    const removeFromPieceBox = (pieceType: PieceType, count: number): boolean => {
      const cnt = this.pieceBox.get(pieceType);
      if (cnt && cnt >= count) {
        this.pieceBox.set(pieceType, cnt - count);
        return true;
      }
      console.error(`Cannot remove from piece box: ${pieceType}`);
      return false;
    };

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
        if (!removeFromPieceBox(pieceType & ~PieceType.PROMOTE, 1)) {
          return;
        }
        x++;
      }
      if (x != 9) {
        console.error(`Row ${i + 1} has incorrect length: expected 9 positions, got ${x}`);
        return;
      }
    }

    // 手番の解析
    if (parts[1] != 'b' && parts[1] != 'w') {
      console.error(`Invalid turn indicator: expected 'b' or 'w', got '${parts[1]}'`);
      return;
    }
    this.isBlackTurn = parts[1] === 'b';

    // 持ち駒の解析
    if (parts[2] !== '-') {
      let count = 1;
      for (let i = 0; i < parts[2].length; i++) {
        const c = parts[2][i];
        // ToDo: 2桁の数字への対応
        if (/[1-9]/.test(c)) {
          if (i === parts[2].length - 1) {
            console.error('Number at end of hands section');
            return;
          }
          count = parseInt(c);
          continue;
        }
        const pieceType = PieceTypeOfSFEN.get(c);
        if (pieceType === undefined || pieceType === PieceType.PROMOTE) {
          console.error(`Invalid piece '${c}' in hands`);
          return;
        }
        if (c >= 'A' && c <= 'Z') {
          if (this.blackHands.has(pieceType)) {
            console.error(`Duplicate piece '${c}' in black hands`);
            return;
          }
          this.blackHands.set(pieceType, count);
        } else {
          if (this.whiteHands.has(pieceType)) {
            console.error(`Duplicate piece '${c}' in white hands`);
            return;
          }
          this.whiteHands.set(pieceType, count);
        }
        if (!removeFromPieceBox(pieceType, count)) {
          return;
        }
        count = 1;
      }
    }

    // 手数
    const moveCount = parseInt(parts[3]);
    if (isNaN(moveCount)) {
      console.error(`Invalid move number: ${parts[3]}`);
      return;
    }
  }
}
