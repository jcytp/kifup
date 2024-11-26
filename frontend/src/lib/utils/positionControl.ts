// src/lib/utils/positionControl.ts

import type { BoardPosition, Move, PieceType, CellPosition } from '$lib/types/Kifu';

// 局面を複製する
export function clonePosition(pos: BoardPosition): BoardPosition {
  return {
    pieces: [...pos.pieces.map((p) => ({ ...p, position: { ...p.position } }))],
    hands: pos.hands
      ? {
          black: { ...pos.hands.black },
          white: { ...pos.hands.white },
        }
      : undefined,
  };
}

// 指し手を適用して新しい局面を生成する
export function applyMove(currentPosition: BoardPosition, move: Move): BoardPosition {
  const newPosition = clonePosition(currentPosition);

  // 1. 移動元の駒を削除
  if (move.from) {
    const moveIndex = newPosition.pieces.findIndex(
      (p) => p.position.x === move.from?.x && p.position.y === move.from?.y
    );
    if (moveIndex >= 0) {
      newPosition.pieces.splice(moveIndex, 1);
    }
  }

  // 2. 移動先に駒があれば取る（持ち駒に加える）
  const captureIndex = newPosition.pieces.findIndex(
    (p) => p.position.x === move.to.x && p.position.y === move.to.y
  );
  if (captureIndex >= 0) {
    const capturedPiece = newPosition.pieces[captureIndex];
    // 持ち駒の初期化
    if (!newPosition.hands) {
      newPosition.hands = { black: {}, white: {} };
    }
    // 成り駒は元の駒に戻して持ち駒に加える
    const originalPiece = getOriginalPiece(capturedPiece.piece);
    const hands = capturedPiece.isBlack ? newPosition.hands.white : newPosition.hands.black;
    hands[originalPiece] = (hands[originalPiece] || 0) + 1;
    // 盤上から削除
    newPosition.pieces.splice(captureIndex, 1);
  }

  // 3. 指し手を適用
  if (move.from) {
    // 盤上の駒を動かす場合
    newPosition.pieces.push({
      position: move.to,
      piece: move.isPromoted ? getPromotedPiece(move.piece) : move.piece,
      isBlack: move.moveNumber % 2 === 1, // 奇数手は先手
    });
  } else {
    // 持ち駒を打つ場合
    const isBlack = move.moveNumber % 2 === 1;
    const hands = isBlack ? newPosition.hands?.black : newPosition.hands?.white;
    if (hands && hands[move.piece] && hands[move.piece]! > 0) {
      // 持ち駒から減らす
      hands[move.piece]!--;
      // 盤上に配置
      newPosition.pieces.push({
        position: move.to,
        piece: move.piece,
        isBlack,
      });
    }
  }

  return newPosition;
}

// 指定手数までの局面を生成
export function generatePosition(
  initialPosition: BoardPosition,
  moves: Move[],
  moveIndex: number
): BoardPosition {
  let currentPosition = clonePosition(initialPosition);

  for (let i = 0; i <= moveIndex && i < moves.length; i++) {
    currentPosition = applyMove(currentPosition, moves[i]);
  }

  return currentPosition;
}

// 成り駒の元の駒を取得
function getOriginalPiece(piece: PieceType): PieceType {
  const mapping: { [K in PieceType]?: PieceType } = {
    と: '歩',
    成香: '香',
    成桂: '桂',
    成銀: '銀',
    馬: '角',
    龍: '飛',
  };
  return mapping[piece] || piece;
}

// 駒の成りを取得
function getPromotedPiece(piece: PieceType): PieceType {
  const mapping: { [K in PieceType]?: PieceType } = {
    歩: 'と',
    香: '成香',
    桂: '成桂',
    銀: '成銀',
    角: '馬',
    飛: '龍',
  };
  return mapping[piece] || piece;
}
