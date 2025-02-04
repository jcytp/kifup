// src/lib/apis/social.ts

import { API, type ApiResult } from '$lib/types/API';

export const addKifuLike = async (kifuId: string): Promise<ApiResult> => {
  const result = await API.post(`/api/kifu/${kifuId}/like`, null, true);
  if (!result.ok) {
    console.error('add kifu like error');
    result.data = 'いいねの追加に失敗しました。';
  }
  return result;
};

export const removeKifuLike = async (kifuId: string): Promise<ApiResult> => {
  const result = await API.delete(`/api/kifu/${kifuId}/like`, null, true);
  if (!result.ok) {
    console.error('remove kifu like error');
    result.data = 'いいねの削除に失敗しました。';
  }
  return result;
};

export const getKifuComments = async (kifuId: string): Promise<ApiResult> => {
  const result = await API.get(`/api/kifu/${kifuId}/comments`, null, false);
  if (!result.data) {
    console.error('get kifu comments error: no data');
    result.ok = false;
    result.data = 'コメントの取得に失敗しました。';
  }
  return result;
};

export const addKifuComment = async (kifuId: string, content: string): Promise<ApiResult> => {
  const params = { content };
  const result = await API.post(`/api/kifu/${kifuId}/comment`, params, true);
  if (!result.ok) {
    console.error('add kifu comment error');
    result.data = 'コメントの投稿に失敗しました。';
  }
  return result;
};
