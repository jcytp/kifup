// src/lib/utils/textEncoding.ts

// BOMの定義
const BOM_PATTERNS = {
	UTF8: [0xef, 0xbb, 0xbf],
	UTF16BE: [0xfe, 0xff],
	UTF16LE: [0xff, 0xfe],
	UTF32BE: [0x00, 0x00, 0xfe, 0xff],
	UTF32LE: [0xff, 0xfe, 0x00, 0x00]
};

// 文字コード判定用の関数
export async function detectEncoding(file: File): Promise<string> {
	// まずバイナリとしてファイルを読み込む
	const buffer = await file.slice(0, Math.min(file.size, 4096)).arrayBuffer();
	const bytes = new Uint8Array(buffer);

	// BOMのチェック
	if (matchesBOM(bytes, BOM_PATTERNS.UTF8)) {
		return 'UTF-8';
	}
	if (matchesBOM(bytes, BOM_PATTERNS.UTF16BE)) {
		return 'UTF-16BE';
	}
	if (matchesBOM(bytes, BOM_PATTERNS.UTF16LE)) {
		return 'UTF-16LE';
	}
	if (matchesBOM(bytes, BOM_PATTERNS.UTF32BE)) {
		return 'UTF-32BE';
	}
	if (matchesBOM(bytes, BOM_PATTERNS.UTF32LE)) {
		return 'UTF-32LE';
	}

	// BOMがない場合は文字コードの特徴を確認

	// UTF-8の特徴をチェック
	if (looksLikeUTF8(bytes)) {
		return 'UTF-8';
	}

	// Shift_JISの特徴をチェック
	if (looksLikeShiftJIS(bytes)) {
		return 'Shift_JIS';
	}

	// 判定できない場合はShift_JISと仮定
	return 'Shift_JIS';
}

// BOMパターンのマッチング確認
function matchesBOM(bytes: Uint8Array, pattern: number[]): boolean {
	return pattern.every((value, index) => bytes[index] === value);
}

// UTF-8らしさの判定
function looksLikeUTF8(bytes: Uint8Array): boolean {
	let i = 0;
	while (i < bytes.length) {
		if (bytes[i] >= 0x00 && bytes[i] <= 0x7f) {
			// ASCII範囲
			i++;
		} else if (bytes[i] >= 0xc2 && bytes[i] <= 0xdf) {
			// 2バイト文字
			if (i + 1 >= bytes.length) return false;
			if (bytes[i + 1] < 0x80 || bytes[i + 1] > 0xbf) return false;
			i += 2;
		} else if (bytes[i] >= 0xe0 && bytes[i] <= 0xef) {
			// 3バイト文字
			if (i + 2 >= bytes.length) return false;
			if (bytes[i + 1] < 0x80 || bytes[i + 1] > 0xbf || bytes[i + 2] < 0x80 || bytes[i + 2] > 0xbf)
				return false;
			i += 3;
		} else {
			return false;
		}
	}
	return true;
}

// Shift_JISらしさの判定
function looksLikeShiftJIS(bytes: Uint8Array): boolean {
	let i = 0;
	while (i < bytes.length) {
		if (bytes[i] >= 0x00 && bytes[i] <= 0x7f) {
			// ASCII範囲
			i++;
		} else if ((bytes[i] >= 0x81 && bytes[i] <= 0x9f) || (bytes[i] >= 0xe0 && bytes[i] <= 0xfc)) {
			// 全角文字の1バイト目
			if (i + 1 >= bytes.length) return false;
			if (bytes[i + 1] < 0x40 || bytes[i + 1] > 0xfc || bytes[i + 1] === 0x7f) return false;
			i += 2;
		} else {
			return false;
		}
	}
	return true;
}

// ファイル読み込みの例
export async function readTextFile(file: File): Promise<string> {
	try {
		const encoding = await detectEncoding(file);
		console.log(`Detected encoding: ${encoding}`);

		// エンコーディングを指定してファイルを読み込む
		const reader = new FileReader();
		return new Promise<string>((resolve, reject) => {
			reader.onload = (e) => resolve(e.target?.result as string);
			reader.onerror = () => reject(new Error('ファイルの読み込みに失敗しました'));
			reader.readAsText(file, encoding);
		});
	} catch (error) {
		console.error('エンコーディング判定エラー:', error);
		throw error;
	}
}
