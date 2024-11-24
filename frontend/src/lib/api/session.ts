import { API, type ApiResult } from '$lib/classes/API';
import { token } from '$lib/stores/session';

export const Login = async (email: string, password: string): Promise<ApiResult> => {
	const params = {
		email: email,
		password: password
	};
	const result = await API.post('/api/session/login', params, false);
	if (!result.data) {
		console.error('login error: no data');
		result.ok = false;
		result.data = 'ログインに失敗しました。';
	}
	if (result.ok) {
		token.set(result.data)
		// ToDO: refreshAccountInfo
	}
	return result;
};
