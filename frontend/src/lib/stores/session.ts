// src/lib/stores/session.ts

import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { Account } from '$lib/types/Account';
import { refreshSession } from '$lib/apis/session';
import { getAccount } from '$lib/apis/account';

let callbackLogoutMove: () => void;

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
        console.log('sessionToken.set(null)');
        // ログアウト
        localStorage.removeItem('token');
        logout();
        callbackLogoutMove();
      }
    },
    initialize: async (logoutMove: () => void) => {
      callbackLogoutMove = logoutMove;
      const value = browser ? localStorage.getItem('token') : null;
      if (value) {
        set(value);
        localStorage.setItem('token', value);
        await update();
      } else {
        callbackLogoutMove();
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
        const result = await refreshSession();
        if (result.ok) {
          sessionToken.set(result.data);
        } else {
          // セッション取得失敗のため、ログアウト
          sessionToken.set(null);
        }
      },
      30 * 60 * 1000
    ); // 30分
  }

  // アカウント情報を更新
  const result = await getAccount();
  if (result.ok) {
    account.set(result.data);
  } else {
    // アカウント取得失敗のため、ログアウト
    sessionToken.set(null);
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
