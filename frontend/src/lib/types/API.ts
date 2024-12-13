// src/lib/classes/API.ts

import { sessionToken } from '$lib/stores/session';
import { get } from 'svelte/store';

export interface PaginationResponse {
  total_count: number;
  page: number;
  page_size: number;
  max_page: number;
}

export interface ApiResult {
  ok: boolean;
  data?: any;
  pagination?: PaginationResponse;
}

export class API {
  static server = 'http://192.168.11.12:8080'; // 開発サーバー

  private static async call(
    method: string,
    path: string,
    queries?: any,
    params?: any,
    withToken = true
  ): Promise<ApiResult> {
    let queryString = '';
    if (queries) {
      const search = new URLSearchParams();
      Object.entries(queries).forEach(([key, value]) => {
        if (value != null && value != undefined) {
          if (Array.isArray(value)) {
            // 配列の場合は複数のパラメータとして追加
            value.forEach((v) => search.append(key, v.toString()));
          } else {
            // それ以外は単一のパラメーターとして追加
            search.append(key, value.toString());
          }
        }
      });
      queryString = search.toString();
    }
    const url = queryString ? this.server + path + '?' + queryString : this.server + path;

    const headers: any = {};
    if (params) {
      headers['Content-Type'] = 'application/json';
    }
    if (withToken) {
      const currentToken = get(sessionToken);
      if (currentToken) {
        headers['Authorization'] = `Bearer ${currentToken}`;
      } else {
        return { ok: false, data: 'no token error' };
      }
    }

    const body = params ? JSON.stringify(params) : undefined;

    const response = await fetch(url, {
      method: method,
      // mode: 'same-origin', // 開発サーバーではoriginが異なる
      cache: 'no-cache',
      headers: headers,
      body: body,
    });
    if (!response.ok) {
      return {
        ok: false,
        data: `${response.status}: ${response.statusText}`,
      };
    }
    return await response.json();
  }

  static async get(path: string, queries?: any, withToken = true): Promise<ApiResult> {
    return await this.call('GET', path, queries, undefined, withToken);
  }

  static async post(path: string, params?: any, withToken = true): Promise<ApiResult> {
    return await this.call('POST', path, undefined, params, withToken);
  }

  static async put(path: string, params?: any, withToken = true): Promise<ApiResult> {
    return await this.call('PUT', path, undefined, params, withToken);
  }

  static async delete(path: string, params?: any, withToken = true): Promise<ApiResult> {
    return await this.call('DELETE', path, undefined, params, withToken);
  }
}
