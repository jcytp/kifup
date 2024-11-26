// src/lib/types/Kifu.ts

export interface Kifu {
  id: string;
  ownerId: string;
  title: string;
  matchInfo: GameInfo;
  tags: string[];
  isPublic: boolean;
  initialPosition?: BoardPosition;
  moves: Move[]; // メインラインの指し手リスト
  createdAt?: string; // 棋譜の作成日時
  updatedAt?: string; // 最終更新日時
}

export interface GameInfo {
  black: string; // 先手の対局者名
  white: string; // 後手の対局者名
  date: string; // 対局日
  title?: string; // 大会名や対局タイトル
  place?: string; // 対局場所
  timeLimit?: {
    // 持ち時間
    initial: number; // 初期時間（分）
    byoyomi?: number; // 秒読み（秒）
  };
}

export interface BoardPosition {
  pieces: {
    position: CellPosition;
    piece: PieceType;
    isBlack: boolean; // true: 先手の駒, false: 後手の駒
  }[];
  hands?: {
    // 初期状態での持ち駒
    black: { [K in PieceType]?: number };
    white: { [K in PieceType]?: number };
  };
}

export interface Move {
  moveNumber: number; // 手数
  piece: PieceType; // 駒の種類
  from?: CellPosition; // 移動元（持ち駒の場合はundefined）
  to: CellPosition; // 移動先
  isPromoted?: boolean; // 成りの有無
  isCapture?: boolean; // 駒を取ったかどうか
  captured?: PieceType; // 取った駒の種類（取った場合のみ）
  comment?: string; // この手に対するコメント
  variations?: Move[]; // この手の代わりとなる指し手（分岐）のリスト
}

export interface CellPosition {
  x: number; // 1-9の横位置
  y: number; // 1-9の縦位置
}

export type PieceType =
  | '歩'
  | '香'
  | '桂'
  | '銀'
  | '金'
  | '角'
  | '飛'
  | 'と'
  | '成香'
  | '成桂'
  | '成銀'
  | '馬'
  | '龍'
  | '玉';

// 使用例：
/*
const exampleKifu: Kifu = {
  id: "kifu-1",
  ownerId: "user-1",
  title: "第1回天竜戦",
  matchInfo: {
    black: "先手太郎",
    white: "後手次郎",
    date: "2024-01-01",
    title: "天竺戦",
    timeLimit: {
      initial: 120,
      byoyomi: 30
    }
  },
  tags: ["実戦", "居飛車"],
  isPublic: true,
  moves: [
    {
      moveNumber: 1,
      piece: "歩",
      from: { x: 7, y: 7 },
      to: { x: 7, y: 6 },
      variations: [
        {
          moveNumber: 1,
          piece: "歩",
          from: { x: 5, y: 7 },
          to: { x: 5, y: 6 },
          comment: "５六歩と指す作戦もある"
        }
      ]
    },
    {
      moveNumber: 2,
      piece: "歩",
      from: { x: 3, y: 3 },
      to: { x: 3, y: 4 }
    }
    // ... 以降の指し手
  ],
  createdAt: "2024-01-01 12:00",
  updatedAt: "2024-01-01 15:30"
};
*/
