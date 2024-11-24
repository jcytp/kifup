import { API, type ApiResult } from '$lib/classes/API';
import { account } from '$lib/stores/session';
import type { Account } from '$lib/types/Account';

export const Register = async (
	name: string,
	email: string,
	password: string
): Promise<ApiResult> => {
	const params = {
		name: name,
		email: email,
		password: password
	};
	const result = await API.post('/api/account', params, false);
	return result;
};

export const GetAccount = async (): Promise<ApiResult> => {
	const result = await API.get('/api/account', null, true);
	if (!result.data) {
		console.error('get account error: no data');
		result.ok = false;
		result.data = 'アカウント情報の取得に失敗しました。';
	}
	if (result.ok) {
        account.set(result.data as Account)
	}
	return result;
};
