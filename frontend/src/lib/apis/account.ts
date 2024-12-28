// src/lib/apis/account.ts

import { API, type ApiResult } from '$lib/types/API';
import { account } from '$lib/stores/session';
import type { Account } from '$lib/types/Account';

export const register = async (
  name: string,
  email: string,
  password: string
): Promise<ApiResult> => {
  const params = {
    name: name,
    email: email,
    password: password,
  };
  const result = await API.post('/api/account', params, false);
  return result;
};

export const getAccount = async (): Promise<ApiResult> => {
  const result = await API.get('/api/account', null, true);
  if (!result.data) {
    console.error('get account error: no data');
    result.ok = false;
    result.data = 'アカウント情報の取得に失敗しました。';
  }
  if (result.ok) {
    account.set(result.data as Account);
  }
  return result;
};

export const getAccountById = async (accountId: string): Promise<ApiResult> => {
  const result = await API.get(`/api/account/${accountId}`, null, false);
  if (!result.data) {
    console.error('get account error: no data');
    result.ok = false;
    result.data = 'アカウント情報の取得に失敗しました。';
  }
  return result;
};

export const updateAccountInfo = async (
  name: string,
  icon_id?: string,
  introduction?: string
): Promise<ApiResult> => {
  const params = {
    name: name,
    icon_id: icon_id || '',
    introduction: introduction || '',
  };
  const result = await API.put('/api/account/info', params, true);
  return result;
};
