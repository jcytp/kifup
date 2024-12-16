<!-- src/routes/kifu/new/+page.svelte -->

<script lang="ts">
  import PositionEditor from '$lib/components/PositionEditor.svelte';
  import { readTextFile } from '$lib/utils/textEncoding';

  // 許可する拡張子のリスト
  const ALLOWED_EXTENSIONS = ['.txt', '.kif', '.csa'];

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

  let currentPosition: string | undefined;

  function handlePositionChange(position?: string) {
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

<div class="page">
  <section class="basic">
    <h2>棋譜を作成 - データから作成</h2>
    <form on:submit|preventDefault={createFromKifu} class="basic">
      <div class="form-group kifu-input-layout">
        <textarea
          id="kifu-textarea"
          bind:value={kifuContent}
          placeholder={'棋譜データを入力してください。\rファイルから読み込むか、直接テキストを貼り付けることができます。'}
        ></textarea>
        <div class="kifu-file-container">
          <input
            class="kifu-file-input"
            type="file"
            accept={ALLOWED_EXTENSIONS.join(',')}
            on:change={handleFileSelect}
          />
          <p>対応形式： {ALLOWED_EXTENSIONS.join(', ')}</p>
        </div>
      </div>
      <button type="submit" class="submit">棋譜データから作成</button>
    </form>
  </section>

  <hr />

  <section class="basic">
    <h2>棋譜を作成 - 局面から作成</h2>
    <PositionEditor onChange={handlePositionChange} />
    <button on:click={createFromPosition} class="submit">この局面から作成</button>
  </section>
</div>

<style lang="scss">
  .kifu-input-layout {
    display: flex;
    gap: 1.5rem;

    #kifu-textarea {
      flex: 1;
      height: calc(1.2rem * 8 + 0.8rem + 2px);
    }

    .kifu-file-container {
      width: 30%;
      display: flex;
      flex-direction: column;
      justify-content: end;

      input.kifu-file-input {
        padding: 0;
        border: none;
        border-radius: 0;
        background: none;
        font-size: 0.9rem;
      }

      p {
        font-size: 0.9rem;
        color: #666;
      }
    }
  }
</style>
