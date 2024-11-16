<!-- src/routes/kifu/new/+page.svelte -->

<script lang="ts">
	import InitialPositionEditor from '$lib/components/InitialPositionEditor.svelte';
	import type { BoardPosition } from '$lib/types/Kifu';
	import { readTextFile } from '$lib/utils/textEncoding';

	// 許可する拡張子のリスト
	const ALLOWED_EXTENSIONS = ['.txt', '.kif', '.kifu', '.ki2', '.csa'];

	// 棋譜データの状態管理
	let kifuContent = '';

	// ファイル選択時の処理
	async function handleFileSelect(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];

		if (!file) return;

		try {
			// ファイル名から拡張子を取得
			const extension = '.' + file.name.split('.').pop()?.toLowerCase();

			// 拡張子チェック
			if (!ALLOWED_EXTENSIONS.includes(extension)) {
				alert('対応していないファイル形式です。\n' + `対応形式: ${ALLOWED_EXTENSIONS.join(', ')}`);
				return;
			}

			// FileReaderでファイルを読み込む
			kifuContent = await readTextFile(file);
		} catch (error) {
			console.error('ファイル読み込みエラー:', error);
			alert(error instanceof Error ? error.message : 'ファイルの読み込みに失敗しました');
		}
	}

	// 棋譜データから作成
	async function createFromKifu() {
		if (!kifuContent.trim()) {
			alert('棋譜データを入力してください');
			return;
		}

		try {
			// TODO: APIリクエストの実装
			console.log('Creating kifu from data:', kifuContent);
			// 成功したら編集ページへリダイレクト
			// window.location.href = `/kifu/edit?id=${newKifuId}`;
		} catch (error) {
			console.error('棋譜作成エラー:', error);
			alert('棋譜の作成に失敗しました');
		}
	}

	let currentPosition: BoardPosition | undefined;

	function handlePositionChange(position: BoardPosition) {
		currentPosition = position;
	}

	// 初期局面から作成
	async function createFromPosition() {
		if (!currentPosition) {
			alert('局面が設定されていません');
			return;
		}

		try {
			// TODO: APIリクエストの実装
			console.log('Creating kifu from position:', currentPosition);
			// 成功したら編集ページへリダイレクト
			// window.location.href = `/kifu/edit?id=${newKifuId}`;
		} catch (error) {
			console.error('棋譜作成エラー:', error);
			alert('棋譜の作成に失敗しました');
		}
	}
</script>

<div class="container">
	<h1>新規棋譜作成</h1>

	<div class="create-methods">
		<!-- 棋譜データから作成 -->
		<section class="create-section">
			<h2>棋譜データから作成</h2>
			<div class="kifu-input-layout">
				<div class="file-input-container">
					<input
						type="file"
						accept={ALLOWED_EXTENSIONS.join(',')}
						on:change={handleFileSelect}
						class="file-input"
					/>
					<p class="hint">
						対応形式: {ALLOWED_EXTENSIONS.join(', ')}
					</p>
				</div>

				<div class="textarea-container">
					<textarea
						bind:value={kifuContent}
						placeholder="棋譜データを入力してください。ファイルから読み込むか、テキストを直接貼り付けることができます。"
						class="kifu-textarea"
					></textarea>
				</div>
			</div>

			<button on:click={createFromKifu} disabled={!kifuContent.trim()} class="create-button">
				棋譜データから作成
			</button>
		</section>

		<!-- 初期局面から作成 -->
		<section class="create-section">
			<h2>初期局面から作成</h2>
			<div class="position-editor-container">
				<InitialPositionEditor change={handlePositionChange} />
			</div>

			<button on:click={createFromPosition} class="create-button"> 初期局面から作成 </button>
		</section>
	</div>
</div>

<style lang="scss">
	.container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem;

		h1 {
			color: var(--primary-color);
			margin-bottom: 2rem;
		}
	}

	.create-methods {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.create-section {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

		h2 {
			color: var(--primary-color);
			margin-bottom: 1.5rem;
			font-size: 1.25rem;
		}
	}

	.kifu-input-layout {
		display: flex;
		gap: 1.5rem;
		margin-bottom: 1rem;
	}

	.file-input-container {
		flex: 0 0 200px;

		.hint {
			margin-top: 0.5rem;
			font-size: 0.9rem;
			color: #666;
		}
	}

	.file-input {
		width: 100%;
		padding: 0.5rem;
		border: 1px solid var(--border-color);
		border-radius: 4px;
		cursor: pointer;
	}

	.textarea-container {
		flex: 1;
	}

	.kifu-textarea {
		width: 100%;
		height: 120px;
		padding: 1rem;
		border: 1px solid var(--border-color);
		border-radius: 4px;
		resize: vertical;
		font-family: monospace;
		line-height: 1.4;

		&:focus {
			border-color: var(--secondary-color);
			outline: none;
		}
	}

	.position-editor-container {
		margin-bottom: 1rem;
	}

	.create-button {
		width: 100%;
		padding: 0.75rem;
		background-color: var(--primary-color);
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-weight: bold;
		transition: opacity 0.2s;

		&:hover:not(:disabled) {
			opacity: 0.9;
		}

		&:disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}
	}

	@media (max-width: 768px) {
		.kifu-input-layout {
			flex-direction: column;
		}

		.file-input-container {
			flex: none;
		}
	}
</style>
