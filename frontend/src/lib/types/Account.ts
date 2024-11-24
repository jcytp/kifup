// src/lib/types/Account.ts

export interface Account {
	id: string;
	name: string;
	icon_id?: string;
	introduction?: string;
	email?: string;
	created_at: Date;
	last_login_at: Date;
}

