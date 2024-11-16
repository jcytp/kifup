// src/lib/test/positions.ts

import type { BoardPosition } from '$lib/types/Kifu';

// 初期局面
export const initialPosition: BoardPosition = {
  pieces: [
    // 後手の駒（上側）
    { position: { x: 1, y: 1 }, piece: '香', isBlack: false },
    { position: { x: 2, y: 1 }, piece: '桂', isBlack: false },
    { position: { x: 3, y: 1 }, piece: '銀', isBlack: false },
    { position: { x: 4, y: 1 }, piece: '金', isBlack: false },
    { position: { x: 5, y: 1 }, piece: '玉', isBlack: false },
    { position: { x: 6, y: 1 }, piece: '金', isBlack: false },
    { position: { x: 7, y: 1 }, piece: '銀', isBlack: false },
    { position: { x: 8, y: 1 }, piece: '桂', isBlack: false },
    { position: { x: 9, y: 1 }, piece: '香', isBlack: false },
    { position: { x: 2, y: 2 }, piece: '飛', isBlack: false },
    { position: { x: 8, y: 2 }, piece: '角', isBlack: false },
    { position: { x: 1, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 2, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 3, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 4, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 5, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 6, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 7, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 8, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 9, y: 3 }, piece: '歩', isBlack: false },

    // 先手の駒（下側）
    { position: { x: 1, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 2, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 3, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 4, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 5, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 6, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 7, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 8, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 9, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 2, y: 8 }, piece: '角', isBlack: true },
    { position: { x: 8, y: 8 }, piece: '飛', isBlack: true },
    { position: { x: 1, y: 9 }, piece: '香', isBlack: true },
    { position: { x: 2, y: 9 }, piece: '桂', isBlack: true },
    { position: { x: 3, y: 9 }, piece: '銀', isBlack: true },
    { position: { x: 4, y: 9 }, piece: '金', isBlack: true },
    { position: { x: 5, y: 9 }, piece: '玉', isBlack: true },
    { position: { x: 6, y: 9 }, piece: '金', isBlack: true },
    { position: { x: 7, y: 9 }, piece: '銀', isBlack: true },
    { position: { x: 8, y: 9 }, piece: '桂', isBlack: true },
    { position: { x: 9, y: 9 }, piece: '香', isBlack: true }
  ]
};

// 途中の局面（四間飛車vs角換わりの例）
export const middlePosition: BoardPosition = {
  pieces: [
    // 盤面の駒
    { position: { x: 1, y: 1 }, piece: '香', isBlack: false },
    { position: { x: 2, y: 1 }, piece: '桂', isBlack: false },
    { position: { x: 3, y: 1 }, piece: '銀', isBlack: false },
    { position: { x: 4, y: 1 }, piece: '金', isBlack: false },
    { position: { x: 5, y: 1 }, piece: '玉', isBlack: false },
    { position: { x: 6, y: 1 }, piece: '金', isBlack: false },
    { position: { x: 7, y: 1 }, piece: '銀', isBlack: false },
    { position: { x: 8, y: 1 }, piece: '桂', isBlack: false },
    { position: { x: 9, y: 1 }, piece: '香', isBlack: false },
    { position: { x: 6, y: 2 }, piece: '飛', isBlack: false },
    { position: { x: 3, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 4, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 5, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 6, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 7, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 8, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 9, y: 3 }, piece: '歩', isBlack: false },
    { position: { x: 2, y: 4 }, piece: '角', isBlack: true },
    { position: { x: 1, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 2, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 3, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 4, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 5, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 6, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 7, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 8, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 9, y: 7 }, piece: '歩', isBlack: true },
    { position: { x: 6, y: 8 }, piece: '飛', isBlack: true },
    { position: { x: 1, y: 9 }, piece: '香', isBlack: true },
    { position: { x: 2, y: 9 }, piece: '桂', isBlack: true },
    { position: { x: 3, y: 9 }, piece: '銀', isBlack: true },
    { position: { x: 4, y: 9 }, piece: '金', isBlack: true },
    { position: { x: 5, y: 9 }, piece: '玉', isBlack: true },
    { position: { x: 6, y: 9 }, piece: '金', isBlack: true },
    { position: { x: 7, y: 9 }, piece: '銀', isBlack: true },
    { position: { x: 8, y: 9 }, piece: '桂', isBlack: true },
    { position: { x: 9, y: 9 }, piece: '香', isBlack: true }
  ],
  hands: {
    black: {
      '歩': 2
    },
    white: {
      '角': 1,
      '歩': 1
    }
  }
};
