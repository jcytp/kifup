// src/lib/stores/session.ts

import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { Account } from '$lib/types/Account';
import { refreshSession } from '$lib/apis/session';
import { getAccount } from '$lib/apis/account';

const createSessionTokenStore = () => {
  const { subscribe, set } = writable<string | null>(null);
  const token = {
    subscribe,
    set: async (value: string | null) => {
      set(value);
      if (value) {
        // ログイン or トークン更新
        localStorage.setItem('token', value);
        await update();
      } else {
        // ログアウト
        localStorage.removeItem('token');
        logout();
      }
    },
    initialize: async () => {
      const value = browser ? localStorage.getItem('token') : null;
      if (value) {
        set(value);
        localStorage.setItem('token', value);
        await update();
      }
    },
  };
  return token;
};

let sessionRefreshTimer: ReturnType<typeof setInterval> | null = null;
export const sessionToken = createSessionTokenStore();
export const account = writable<Account | null>(null);

const update = async () => {
  // リフレッシュタイマーを開始
  if (!sessionRefreshTimer) {
    sessionRefreshTimer = setInterval(
      async () => {
        await refreshSession();
      },
      30 * 60 * 1000
    ); // 30分
  }

  // アカウント情報を更新
  const result = await getAccount();
  if (result.ok) {
    const new_account = result.data as Account;
    account.set(new_account);
  } else {
    sessionToken.set(null); // アカウント取得失敗のため、ログアウト
  }
};

const logout = () => {
  // リフレッシュタイマーを停止
  if (sessionRefreshTimer) {
    clearInterval(sessionRefreshTimer);
    sessionRefreshTimer = null;
  }

  // アカウント情報を削除
  account.set(null);
};
