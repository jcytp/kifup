<!-- src/routes/kifu/new/+page.svelte -->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { createKifu } from '$lib/apis/kifu';
  import PositionEditor from '$lib/components/PositionEditor.svelte';
  import { readTextFile } from '$lib/utils/textEncoding';

  // ----------------------------------------
  // データから作成

  const ALLOWED_FORMATS = ['KIF形式', 'CSA形式'];

  let kifuContent = '';

  async function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];

    if (!file) return;

    try {
      kifuContent = await readTextFile(file);
    } catch (error) {
      console.error('ファイル読み込みエラー:', error);
      alert(error instanceof Error ? error.message : 'ファイルの読み込みに失敗しました');
    }
  }

  async function createFromKifu() {
    const data = kifuContent.trim();
    if (!data) {
      alert('棋譜データを入力してください');
      return;
    }

    console.log('Creating kifu from data:', data);
    const result = await createKifu('file', data, undefined);
    if (result.ok && result.data) {
      const kifuID = result.data;
      goto(`/kifu/edit/?id=${kifuID}`); // 編集画面へ遷移
    } else {
      console.error('Failed to create kifu from file: ', result);
    }
  }

  // ----------------------------------------
  // 初期局面から作成

  let sfen: string | undefined = undefined;

  function handlePositionChange(newSfen?: string) {
    sfen = newSfen;
  }

  async function createFromPosition() {
    const result = await createKifu('position', undefined, sfen);
    if (result.ok && result.data) {
      const kifuID = result.data;
      goto(`/kifu/edit/?id=${kifuID}`); // 編集画面へ遷移
    } else {
      console.error('Failed to create kifu from position: ', result);
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
            accept={ALLOWED_FORMATS.join(',')}
            on:change={handleFileSelect}
          />
          <p>対応形式： {ALLOWED_FORMATS.join(', ')}</p>
        </div>
      </div>
      <button type="submit" class="submit">棋譜データから作成</button>
      <p>※現在KIF形式のみ、分岐には未対応です</p>
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
