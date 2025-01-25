// src/lib/types/Kifu.ts

import type { Account } from './Account';
import type { PieceType } from './Piece';

export interface KifuComment {
  id: string;
  kifu_id: string;
  account: Account;
  content: string;
  created_at: string;
}

export interface KifuSummary {
  id: string;
  owner: Account;
  title: string;
  is_public: boolean;
  updated_at: Date;
  game_info: { [key: string]: string };
  tags: string[];
  like_count: number;
  comment_count: number;
}

export interface KifuDetail {
  id: string;
  owner: Account;
  title: string;
  is_public: boolean;
  initial_position?: string; // SFEN format
  created_at: string;
  updated_at: string;
  game_info: { [key: string]: string };
  tags: string[];
  moves: KifuMove[];
  like_count: number;
  has_like: boolean;
}

export interface KifuMove {
  number: number;
  piece: PieceType;
  from_place: number; // Use PIECE_PLACE_IN_HAND for moves from hand
  to_place: number;
  promote?: boolean;
  catch_piece?: PieceType;
  direction_sign?: string;
  variations?: KifuMove[][]; // Array of move lines
  comment?: string;
  time_spent_ms?: number;
}
