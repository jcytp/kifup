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
          date: '2024-01-01',
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

  onMount(() => {
    if (kifuId) {
      fetchKifuData();
    } else {
      error = '棋譜IDが指定されていません。';
      isLoading = false;
    }
  });
</script>

<div class="container">
  {#if isLoading}
    <div class="loading">
      <p>棋譜を読み込んでいます...</p>
      <a href="/kifu/search" class="back-button">棋譜検索に戻る</a>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
      <a href="/kifu/search" class="back-button">棋譜検索に戻る</a>
    </div>
  {:else if kifu}
    <div class="edit-form">
      <header class="form-header">
        <h1>棋譜の編集</h1>
        <div class="header-actions">
          <label class="public-toggle">
            <input type="checkbox" bind:checked={formData.isPublic} />
            公開する
          </label>
          <button class="save-button" on:click={handleSave}>保存</button>
        </div>
      </header>

      <section class="basic-info">
        <div class="form-group">
          <label for="title">タイトル</label>
          <input
            type="text"
            id="title"
            bind:value={formData.title}
            placeholder="棋譜のタイトルを入力"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="black">先手</label>
            <input
              type="text"
              id="black"
              bind:value={formData.matchInfo.black}
              placeholder="先手の対局者名"
            />
          </div>
          <div class="form-group">
            <label for="white">後手</label>
            <input
              type="text"
              id="white"
              bind:value={formData.matchInfo.white}
              placeholder="後手の対局者名"
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="match-date">対局日</label>
            <input type="date" id="match-date" bind:value={formData.matchInfo.date} />
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

        <div class="form-row">
          <div class="form-group">
            <label for="match-title">大会名</label>
            <input
              type="text"
              id="match-title"
              bind:value={formData.matchInfo.title}
              placeholder="大会名や対局タイトル"
            />
          </div>
        </div>

        <div class="form-row">
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
        </div>

        <div class="form-group">
          <h3 class="form-heading">タグ</h3>
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
      </section>

      <section class="kifu-editor">
        <h2>棋譜の編集</h2>
        <div class="editor-container">
          <KifuPlayer {kifu} />
        </div>
      </section>
    </div>
  {/if}
</div>

<style lang="scss">
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  .loading,
  .error {
    text-align: center;
    padding: 2rem;

    .back-button {
      display: inline-block;
      margin-top: 1rem;
      padding: 0.5rem 1rem;
      background-color: var(--primary-color);
      color: white;
      border-radius: 4px;
      cursor: pointer;

      &:hover {
        opacity: 0.9;
      }
    }
  }

  .edit-form {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 2rem;

    .form-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 2rem;

      h1 {
        margin: 0;
        color: var(--primary-color);
      }

      .header-actions {
        display: flex;
        align-items: center;
        gap: 1rem;

        .public-toggle {
          display: flex;
          align-items: center;
          gap: 0.5rem;
          cursor: pointer;

          input[type='checkbox'] {
            cursor: pointer;
          }
        }

        .save-button {
          padding: 0.5rem 1.5rem;
          background-color: var(--primary-color);
          color: white;
          border-radius: 4px;
          cursor: pointer;

          &:hover {
            opacity: 0.9;
          }
        }
      }
    }

    .basic-info {
      margin-bottom: 2rem;
    }

    .form-heading {
      margin-bottom: 0.5rem;
      color: var(--text-color);
      font-weight: bold;
      font-size: 1rem; // 他のラベルと同じサイズに合わせる
    }

    .form-group {
      margin-bottom: 1rem;

      label {
        display: block;
        margin-bottom: 0.5rem;
        color: var(--text-color);
        font-weight: bold;
      }

      input {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid var(--border-color);
        border-radius: 4px;
        cursor: text;

        &:focus {
          border-color: var(--secondary-color);
        }
      }
    }

    .form-row {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 1rem;
      margin-bottom: 1rem;
    }

    .tag-input-container {
      display: flex;
      gap: 0.5rem;
      margin-bottom: 0.5rem;

      input {
        flex: 1;
      }

      .add-tag-btn {
        padding: 0.5rem 1rem;
        background-color: var(--secondary-color);
        color: white;
        border-radius: 4px;
        cursor: pointer;

        &:hover {
          opacity: 0.9;
        }
      }
    }

    .tags-container {
      display: flex;
      flex-wrap: wrap;
      gap: 0.5rem;

      .tag {
        background-color: var(--secondary-color);
        color: white;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        display: flex;
        align-items: center;
        gap: 0.25rem;

        .remove-tag {
          background: none;
          border: none;
          color: white;
          cursor: pointer;
          padding: 0 0.25rem;

          &:hover {
            opacity: 0.8;
          }
        }
      }
    }

    .kifu-editor {
      h2 {
        margin-bottom: 1rem;
        color: var(--primary-color);
      }

      .editor-container {
        border: 1px solid var(--border-color);
        border-radius: 4px;
        padding: 1rem;
      }
    }
  }
</style>
