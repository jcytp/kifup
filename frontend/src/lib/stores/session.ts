// src/lib/stores/session.ts

import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { Account } from '$lib/types/Account';

const createTokenStore = () => {
    const initialValue = browser ? localStorage.getItem('token') : null;
    const { subscribe, set } = writable<string | null>(initialValue);
    const token = {
        subscribe,
        set: (token: string | null) => {
            // if (browser) {
                if (token) {
                    localStorage.setItem('token', token);
                } else {
                    localStorage.removeItem('token');
                }
            // }
            set(token);
        },
        clear: () => {
            // if (browser) {
                localStorage.removeItem('token');
            // }
            set(null);
        }
    };
    return token;
}

export const token = createTokenStore();

export const account = writable<Account | null>(null);
