// src/lib/types/BoardPosition.ts

import type { KifuMove } from './Kifu';
import {
  PieceCharOfSFEN,
  PieceOrderForSFEN,
  PiecePlace,
  PieceType,
  PieceTypeOfSFEN,
} from './Piece';

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

  copy(): BoardPosition {
    return new BoardPosition(this.toSfen(1));
  }

  toSfen(moveCount: number): string {
    if (moveCount < 1) {
      console.error('Invalid move count: must be positive');
      return SFEN_HIRATE;
    }

    const parts: string[] = [];

    // 盤面作成
    const rows: string[] = [];
    for (let i = 0; i < 9; i++) {
      let row = '';
      let emptyCount = 0;

      for (let j = 0; j < 9; j++) {
        const blackPiece = this.blackBoard[i][j];
        const whitePiece = this.whiteBoard[i][j];
        if (blackPiece === PieceType.VACANCY && whitePiece === PieceType.VACANCY) {
          emptyCount++;
          continue;
        }
        if (emptyCount > 0) {
          // 空マス
          row += emptyCount.toString();
          emptyCount = 0;
        }
        if (whitePiece == PieceType.VACANCY) {
          // 先手の駒
          const char = PieceCharOfSFEN.get(blackPiece & ~PieceType.PROMOTE);
          if (char) {
            if (blackPiece & PieceType.PROMOTE) {
              row += '+';
            }
            row += char.toUpperCase();
          }
        } else if (blackPiece === PieceType.VACANCY) {
          // 後手の駒
          const char = PieceCharOfSFEN.get(whitePiece & ~PieceType.PROMOTE);
          if (char) {
            if (whitePiece & PieceType.PROMOTE) {
              row += '+';
            }
            row += char.toLowerCase();
          }
        }
      }

      if (emptyCount > 0) {
        // 行末の空マス
        row += emptyCount.toString();
      }

      rows.push(row);
    }
    parts.push(rows.join('/'));

    // 手番
    parts.push(this.isBlackTurn ? 'b' : 'w');

    // 持ち駒
    const hands: string[] = [];
    for (const piece of PieceOrderForSFEN) {
      const count = this.blackHands.get(piece);
      if (count) {
        if (count > 1) {
          hands.push(count.toString());
        }
        const char = PieceCharOfSFEN.get(piece);
        if (char) {
          hands.push(char.toUpperCase());
        }
      }
    }
    for (const piece of PieceOrderForSFEN) {
      const count = this.whiteHands.get(piece);
      if (count) {
        if (count > 1) {
          hands.push(count.toString());
        }
        const char = PieceCharOfSFEN.get(piece);
        if (char) {
          hands.push(char.toLowerCase());
        }
      }
    }
    parts.push(hands.length > 0 ? hands.join('') : '-');

    // 手数
    parts.push(moveCount.toString());

    return parts.join(' ');
  }

  next(move: KifuMove): void {
    const toPlace = new PiecePlace(move.to_place);
    const toRow = toPlace.row();
    const toCol = toPlace.col();
    const fromPlace = new PiecePlace(move.from_place);

    // 移動先に駒がある場合、持ち駒へ追加
    const target = this.isBlackTurn ? this.whiteBoard[toRow][toCol] : this.blackBoard[toRow][toCol];
    if (target !== PieceType.VACANCY) {
      if (target === PieceType.OU) {
        const num = this.pieceBox.get(target) ?? 0;
        this.pieceBox.set(target, num + 1);
      } else {
        const original = target & ~PieceType.PROMOTE;
        if (this.isBlackTurn) {
          const num = this.blackHands.get(original) ?? 0;
          this.blackHands.set(original, num + 1);
        } else {
          const num = this.whiteHands.get(original) ?? 0;
          this.whiteHands.set(original, num + 1);
        }
      }
    }

    // 移動先に駒を配置
    const piece = move.piece | (move.promote ? PieceType.PROMOTE : 0);
    if (this.isBlackTurn) {
      this.blackBoard[toRow][toCol] = piece;
      this.whiteBoard[toRow][toCol] = PieceType.VACANCY;
    } else {
      this.whiteBoard[toRow][toCol] = piece;
      this.blackBoard[toRow][toCol] = PieceType.VACANCY;
    }

    // 移動元の駒を削除
    if (fromPlace.isInHand()) {
      if (this.isBlackTurn) {
        const num = this.blackHands.get(piece) || 0;
        this.blackHands.set(piece, num - 1);
      } else {
        const num = this.whiteHands.get(piece) || 0;
        this.whiteHands.set(piece, num - 1);
      }
    } else {
      const fromRow = fromPlace.row();
      const fromCol = fromPlace.col();
      if (this.isBlackTurn) {
        this.blackBoard[fromRow][fromCol] = PieceType.VACANCY;
      } else {
        this.whiteBoard[fromRow][fromCol] = PieceType.VACANCY;
      }
    }

    // 手番を更新
    this.isBlackTurn = !this.isBlackTurn;
  }
}
