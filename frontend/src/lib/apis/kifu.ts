// src/lib/apis/kifu.ts

import { API, type ApiResult } from '$lib/types/API';
import type { KifuMove } from '$lib/types/Kifu';

export const searchKifus = async (
  owner: string | null,
  page: number,
  page_size: number,
  isLoggedIn: boolean
): Promise<ApiResult> => {
  const params = {
    owner: owner,
    page: page,
    page_size: page_size,
  };
  const result = await API.get('/api/kifu', params, isLoggedIn);
  if (!result.data) {
    console.error('search kifus error: no data');
    result.ok = false;
    result.data = '棋譜リストの取得に失敗しました。';
  }
  return result;
};

export const createKifu = async (
  type: 'file' | 'position',
  content?: string,
  initialPosition?: string
): Promise<ApiResult> => {
  const params = {
    type: type,
    content: content,
    initial_position: initialPosition,
  };
  const result = await API.post('/api/kifu', params, true);
  if (!result.data) {
    console.error('create kifu error: no data');
    result.ok = false;
    result.data = '棋譜の作成に失敗しました。';
  }
  return result;
};

export const getKifu = async (kifuId: string, withToken: boolean): Promise<ApiResult> => {
  const result = await API.get(`/api/kifu/${kifuId}`, null, withToken);
  if (!result.data) {
    console.error('get kifu error: no data');
    result.ok = false;
    result.data = '棋譜の取得に失敗しました。';
  }
  return result;
};

export const deleteKifu = async (kifuId: string): Promise<ApiResult> => {
  const result = await API.delete(`/api/kifu/${kifuId}`, null, true);
  if (!result.ok) {
    console.error('delete kifu error');
    result.data = '棋譜の削除に失敗しました。';
  }
  return result;
};

export const updateKifuInfo = async (
  kifuId: string,
  title: string,
  isPublic: boolean,
  gameInfo: { [key: string]: string },
  tags: string[]
): Promise<ApiResult> => {
  const params = {
    title,
    is_public: isPublic,
    game_info: gameInfo,
    tags,
  };
  const result = await API.put(`/api/kifu/${kifuId}`, params, true);
  if (!result.ok) {
    console.error('update kifu info error');
    result.data = '棋譜情報の更新に失敗しました。';
  }
  return result;
};

export const updateKifuMoves = async (kifuId: string, moves: KifuMove[]): Promise<ApiResult> => {
  const params = { moves };
  const result = await API.put(`/api/kifu/${kifuId}/moves`, params, true);
  if (!result.ok) {
    console.error('update kifu moves error');
    result.data = '棋譜の指し手の更新に失敗しました。';
  }
  return result;
};
