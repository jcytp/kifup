<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import type { GameInfo, Kifu, Move } from '$lib/types/Kifu';
  import KifuPlayer from '$lib/components/KifuPlayer.svelte';
  import { initialPosition } from '$lib/test/positions';

  // フォーム用の型定義
  interface KifuFormData {
    title: string;
    matchInfo: {
      black: string;
      white: string;
      date: string;
      place: string;
      title: string;
      timeLimit: {
        initial: number;
        byoyomi: number;
        increment: number;
      };
    };
    tags: string[];
    isPublic: boolean;
  }

  // URLからkifu_idを取得
  const kifuId = $page.url.searchParams.get('id');

  // 棋譜データの状態管理
  let kifu: Kifu | null = null;
  let isLoading = true;
  let error: string | null = null;

  // フォームの状態管理
  let formData: KifuFormData = {
    title: '',
    matchInfo: {
      black: '',
      white: '',
      date: '',
      place: '',
      title: '',
      timeLimit: {
        initial: 0,
        byoyomi: 0,
        increment: 0,
      },
    },
    tags: [],
    isPublic: false,
  };

  // タグの入力管理
  let tagInput = '';

  // 棋譜の指し手管理
  let moves: Move[] = [];

  // GameInfo型からフォームデータへの変換
  function gameInfoToFormData(gameInfo: GameInfo): KifuFormData['matchInfo'] {
    return {
      black: gameInfo.black,
      white: gameInfo.white,
      date: gameInfo.date,
      place: gameInfo.place || '',
      title: gameInfo.title || '',
      timeLimit: {
        initial: gameInfo.timeLimit?.initial || 0,
        byoyomi: gameInfo.timeLimit?.byoyomi || 0,
        increment: gameInfo.timeLimit?.increment || 0,
      },
    };
  }

  // フォームデータからGameInfo型への変換
  function formDataToGameInfo(formMatchInfo: KifuFormData['matchInfo']): GameInfo {
    const gameInfo: GameInfo = {
      black: formMatchInfo.black,
      white: formMatchInfo.white,
      date: formMatchInfo.date,
    };

    if (formMatchInfo.place) gameInfo.place = formMatchInfo.place;
    if (formMatchInfo.title) gameInfo.title = formMatchInfo.title;
    if (formMatchInfo.timeLimit.initial || formMatchInfo.timeLimit.byoyomi) {
      gameInfo.timeLimit = {
        initial: formMatchInfo.timeLimit.initial,
        byoyomi: formMatchInfo.timeLimit.byoyomi,
        increment: formMatchInfo.timeLimit.increment,
      };
    }

    return gameInfo;
  }

  // タグの追加
  function addTag() {
    if (tagInput && !formData.tags.includes(tagInput)) {
      formData.tags = [...formData.tags, tagInput];
      tagInput = '';
    }
  }

  // タグの削除
  function removeTag(index: number) {
    formData.tags = formData.tags.filter((_, i) => i !== index);
  }

  // 指し手の更新コールバック関数
  function handleMovesUpdate(new_moves: Move[]) {
    moves = new_moves;
  }

  // 棋譜データの取得
  async function fetchKifuData() {
    isLoading = true;
    try {
      // TODO: API実装後に実際のデータ取得に置き換え
      await new Promise((resolve) => setTimeout(resolve, 500)); // ローディング表示確認用

      // テストデータ
      const testData: Kifu = {
        id: kifuId || '',
        ownerId: 'user-1',
        title: 'テスト棋譜',
        matchInfo: {
          black: '先手太郎',
          white: '後手次郎',
          date: '2024-01-01T00:00',
          place: '対局場所',
          title: '大会名',
          timeLimit: {
            initial: 60,
            byoyomi: 30,
          },
        },
        tags: ['実戦', 'テスト'],
        isPublic: false,
        initialPosition: initialPosition,
        moves: [
          {
            moveNumber: 1,
            piece: '歩',
            from: { x: 7, y: 7 },
            to: { x: 7, y: 6 },
            comment: '序盤の一手',
          },
        ],
      };

      kifu = testData;
      moves = [...testData.moves];
      formData = {
        title: testData.title,
        matchInfo: gameInfoToFormData(testData.matchInfo),
        tags: [...testData.tags],
        isPublic: testData.isPublic,
      };
    } catch (e) {
      error = '棋譜データの取得に失敗しました。';
    } finally {
      isLoading = false;
    }
  }

  // 更新の保存
  async function handleSave() {
    try {
      if (!kifu) return;

      const updatedKifu: Kifu = {
        ...kifu,
        title: formData.title,
        matchInfo: formDataToGameInfo(formData.matchInfo),
        tags: formData.tags,
        isPublic: formData.isPublic,
        moves: moves,
      };

      // TODO: API実装後に実際の保存処理に置き換え
      console.log('Save kifu:', updatedKifu);
      alert('保存しました');
    } catch (e) {
      alert('保存に失敗しました');
    }
  }

  // 棋譜情報の更新
  async function updateKifuInfo() {
    console.debug('update kifu info');
  }

  // 指し手の更新
  async function updateKifuMoves() {
    console.debug('update kifu moves');
  }

  onMount(() => {
    if (kifuId) {
      fetchKifuData();
    } else {
      error = '棋譜IDが指定されていません。';
      isLoading = false;
    }
  });
</script>

<div class="page">
  {#if kifu}
    <section class="basic">
      <h2>棋譜の編集</h2>
      <form on:submit|preventDefault={updateKifuInfo} class="basic">
        <div class="form-group">
          <label for="title">タイトル</label>
          <input
            type="text"
            id="title"
            bind:value={formData.title}
            placeholder="棋譜のタイトル"
            required
          />
        </div>
        <div class="flex-arrange">
          <div class="form-group">
            <label for="black-player">先手</label>
            <input
              type="text"
              id="black-player"
              bind:value={formData.matchInfo.black}
              placeholder="先手の対局者名"
            />
          </div>
          <div class="form-group">
            <label for="white-player">後手</label>
            <input
              type="text"
              id="white-player"
              bind:value={formData.matchInfo.white}
              placeholder="後手の対局者名"
            />
          </div>
        </div>
        <div class="flex-arrange">
          <div class="form-group">
            <label for="start-date">対局日時</label>
            <input type="datetime-local" id="start-date" bind:value={formData.matchInfo.date} />
          </div>
          <div class="form-group">
            <label for="place">対局場所</label>
            <input
              type="text"
              id="place"
              bind:value={formData.matchInfo.place}
              placeholder="対局場所"
            />
          </div>
        </div>
        <div class="flex-arrange">
          <div class="form-group">
            <label for="initial-time">持ち時間（分）</label>
            <input
              type="number"
              id="initial-time"
              bind:value={formData.matchInfo.timeLimit.initial}
              min="0"
            />
          </div>
          <div class="form-group">
            <label for="byoyomi">秒読み（秒）</label>
            <input
              type="number"
              id="byoyomi"
              bind:value={formData.matchInfo.timeLimit.byoyomi}
              min="0"
            />
          </div>
          <div class="form-group">
            <label for="increment">加算（秒）</label>
            <input
              type="number"
              id="increment"
              bind:value={formData.matchInfo.timeLimit.increment}
              min="0"
            />
          </div>
        </div>
        <div class="form-group">
          <h3 class="label">タグ</h3>
          <div class="tag-input-container">
            <input
              type="text"
              bind:value={tagInput}
              placeholder="タグを入力"
              on:keydown={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
            />
            <button type="button" on:click={addTag} class="add-tag-btn">追加</button>
          </div>
          <div class="tags-container">
            {#each formData.tags as tag, i}
              <span class="tag">
                {tag}
                <button type="button" on:click={() => removeTag(i)} class="remove-tag">×</button>
              </span>
            {/each}
          </div>
        </div>
        <div class="form-group">
          <h3 class="label">公開設定</h3>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={formData.isPublic} />
            公開する
          </label>
        </div>
        <button type="submit" class="submit">棋譜情報を更新</button>
      </form>

      <KifuPlayer {kifu} />
      <button on:click={updateKifuMoves} class="submit">棋譜の指し手を更新</button>
    </section>
  {/if}
</div>

<style lang="scss">
  .tag-input-container {
    display: flex;
    gap: 0.5rem;

    input {
      flex: 1;
    }

    button {
      padding: 0.4rem 1rem;
      border: 1px solid var(--primary-color);
      border-radius: 0.4rem;
      background-color: var(--primary-color);
      color: white;

      &:hover,
      &:focus {
        border-color: var(--secondary-color);
        background-color: var(--secondary-color);
      }
    }
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-top: 0.5rem;

    .tag {
      background-color: var(--secondary-color);
      color: white;
      padding: 0.3rem 0.4rem 0.3rem 0.6rem;
      border-radius: 0.3rem;
      display: flex;
      align-items: center;
      gap: 0.4rem;
      font-size: 0.9rem;
      line-height: 1.4rem;

      .remove-tag {
        background: none;
        border: none;
        color: white;
        cursor: pointer;
      }
    }
  }
</style>
