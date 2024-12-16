// src/lib/apis/session.ts

import { API, type ApiResult } from '$lib/types/API';

export const login = async (email: string, password: string): Promise<ApiResult> => {
  const params = {
    email: email,
    password: password,
  };
  const result = await API.post('/api/session/login', params, false);
  if (!result.data) {
    console.error('login error: no data');
    result.ok = false;
    result.data = 'ログインに失敗しました。';
  }
  return result;
};

export const refreshSession = async (): Promise<ApiResult> => {
  const result = await API.post('/api/session/refresh', null, true);
  if (!result.data) {
    console.error('refresh session error: no data');
    result.ok = false;
    result.data = 'トークン取得に失敗しました。';
  }
  return result;
};
