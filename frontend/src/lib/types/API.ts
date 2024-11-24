// src/lib/classes/API.ts

import { sessionToken } from "$lib/stores/session";
import { get } from 'svelte/store';

export interface ApiResult {
    ok: boolean;
    data?: any;
}

export class API {
    static server = 'http://192.168.11.12:8080'; // 開発サーバー

    static async get(path: string, params: any, withToken=true): Promise<ApiResult> {
        const url = this.server + path;
        const headers: any = {};
        if (withToken) {
            const currentToken = get(sessionToken);
            if (currentToken) {
                headers['Authorization'] = `Bearer ${currentToken}`;
            } else {
                return {ok: false, data: 'no token error'}
            }
        }
        const response = await fetch(url, {
            method: 'GET',
            // mode: 'same-origin', // 開発サーバーではoriginが異なる
            cache: 'no-cache',
            headers: headers,
        });
        if (!response.ok) {
            return {
                ok: false,
                data: `${response.status}: ${response.statusText}`,
            };
        }
        return await response.json();
    }

    static async post(path: string, params: any, withToken=true): Promise<ApiResult> {
        const url = this.server + path;
        const headers: any = {
            'Content-Type': 'application/json',
        };
        if (withToken) {
            const currentToken = get(sessionToken);
            if (currentToken) {
                headers['Authorization'] = `Bearer ${currentToken}`;
            } else {
                return {ok: false, data: 'no token error'}
            }
        }
        const response = await fetch(url, {
            method: 'POST',
            // mode: 'same-origin', // 開発サーバーではoriginが異なる
            cache: 'no-cache',
            headers: headers,
            body: JSON.stringify(params),
        });
        if (!response.ok) {
            return {
                ok: false,
                data: `${response.status}: ${response.statusText}`,
            };
        }
        return await response.json();
    }
}
