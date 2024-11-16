<!-- src/routes/kifu/view/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import type { Kifu } from '$lib/types/Kifu';
  import { page } from '$app/stores';
  import KifuPlayer from '$lib/components/KifuPlayer.svelte';
	import { initialPosition } from '$lib/test/positions';

  // 現在のユーザーID（実際にはログイン情報から取得）
  // TODO: 認証機能実装後に実際のユーザーID取得に置き換え
  const currentUserId = 'user-1';

  // URLからkifu_idを取得
  const kifuId = $page.url.searchParams.get('id');

  // 棋譜データの状態管理
  let kifu: Kifu | null = null;
  let isLoading = true;
  let error: string | null = null;

  // コメントの状態管理
  let comments: Array<{
    id: string;
    userId: string;
    userName: string;
    content: string;
    createdAt: string;
  }> = [];
  let newComment = '';

  // いいねの状態管理
  let isLiked = false;
  let likeCount = 0;

  // 所有者かどうかのチェック
  $: isOwner = kifu?.ownerId === currentUserId;

  // 棋譜データの取得
  async function fetchKifuData() {
    isLoading = true;
    try {
      // TODO: API実装後に実際のデータ取得に置き換え
      await new Promise(resolve => setTimeout(resolve, 500)); // ローディング表示確認用
      
      kifu = {
        id: kifuId || '',
        ownerId: 'user-1',
        title: 'テスト棋譜',
        matchInfo: {
          black: '先手太郎',
          white: '後手次郎',
          date: '2024-01-01',
        },
        tags: ['実戦', 'テスト'],
        isPublic: true,
        initialPosition: initialPosition,
        moves: [
          {
            moveNumber: 1,
            piece: "歩",
            from: { x: 7, y: 7 },
            to: { x: 7, y: 6 },
            comment: "歩を突く",
            variations: [
              {
                moveNumber: 1,
                piece: "歩",
                from: { x: 5, y: 7 },
                to: { x: 5, y: 6 },
                comment: "５六歩と指す作戦もある"
              }
            ]
          },
          {
            moveNumber: 2,
            piece: "歩",
            from: { x: 3, y: 3 },
            to: { x: 3, y: 4 }
          },
          {
            moveNumber: 3,
            piece: "歩",
            from: { x: 6, y: 7 },
            to: { x: 6, y: 6 }
          }
        ],
        createdAt: "2024-01-01 12:00",
        updatedAt: "2024-01-01 15:30"
      };
      
      comments = [
        {
          id: '1',
          userId: 'user-2',
          userName: '観戦者A',
          content: '素晴らしい一手でした！',
          createdAt: '2024-01-01 12:00'
        }
      ];
      
      likeCount = 42;
      
    } catch (e) {
      error = '棋譜データの取得に失敗しました。';
    } finally {
      isLoading = false;
    }
  }

  // 日付フォーマット用の関数
  function formatDateTime(dateStr: string): string {
    const date = new Date(dateStr);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}`;
  }

  // コメント投稿
  async function handleCommentSubmit() {
    if (!newComment.trim()) return;
    
    // TODO: API実装後に実際のコメント投稿処理に置き換え
    comments = [
      ...comments,
      {
        id: String(Date.now()),
        userId: 'current-user',
        userName: '自分',
        content: newComment,
        createdAt: formatDateTime(new Date().toISOString())
      }
    ];
    
    newComment = '';
  }

  // いいね処理
  async function toggleLike() {
    // TODO: API実装後に実際のいいね処理に置き換え
    isLiked = !isLiked;
    likeCount += isLiked ? 1 : -1;
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
      <a href="/kifu" class="back-button">棋譜検索に戻る</a>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
      <a href="/kifu" class="back-button">棋譜検索に戻る</a>
    </div>
  {:else if kifu}
    <header class="kifu-header">
      <div class="header-top">
        <h1>{kifu.title}</h1>
        {#if isOwner}
          <a href={`/kifu/edit?id=${kifu.id}`} class="edit-button">
            編集する
          </a>
        {/if}
      </div>
      <div class="kifu-meta">
        <div class="match-info">
          <p>対局日: {kifu.matchInfo.date}</p>
          <p>対局者: {kifu.matchInfo.black} vs {kifu.matchInfo.white}</p>
        </div>
        <div class="tags">
          {#each kifu.tags as tag}
            <span class="tag">{tag}</span>
          {/each}
        </div>
      </div>
    </header>

    <div class="kifu-content">
      <div class="kifu-view-placeholder">
        {#if kifu}
          <KifuPlayer {kifu} />
        {/if}
      </div>

      <div class="action-buttons">
        <button class="like-button" class:liked={isLiked} on:click={toggleLike}>
          {isLiked ? '★' : '☆'} いいね！ ({likeCount})
        </button>
      </div>

      <section class="comments-section">
        <h2>コメント</h2>
        
        <form class="comment-form" on:submit|preventDefault={handleCommentSubmit}>
          <textarea
            bind:value={newComment}
            placeholder="コメントを入力..."
            rows="3"
          ></textarea>
          <button type="submit" disabled={!newComment.trim()}>
            コメントする
          </button>
        </form>

        <div class="comments-list">
          {#if comments.length === 0}
            <p class="no-comments">コメントはまだありません。</p>
          {:else}
            {#each comments as comment}
              <div class="comment">
                <div class="comment-header">
                  <span class="comment-author">{comment.userName}</span>
                  <span class="comment-date">{comment.createdAt}</span>
                </div>
                <p class="comment-content">{comment.content}</p>
              </div>
            {/each}
          {/if}
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

  .loading, .error {
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

  .kifu-header {
    margin-bottom: 2rem;

    .header-top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;

      h1 {
        margin-bottom: 0; // 既存のmarginを削除
      }
    }

    .edit-button {
      padding: 0.5rem 1rem;
      background-color: var(--primary-color);
      color: white;
      border-radius: 4px;
      font-size: 0.9rem;
      cursor: pointer;
      transition: opacity 0.2s;

      &:hover {
        opacity: 0.9;
      }
    }

    .kifu-meta {
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 1rem;

      .match-info {
        color: var(--text-color);
      }

      .tags {
        display: flex;
        gap: 0.5rem;
        flex-wrap: wrap;

        .tag {
          background-color: var(--secondary-color);
          color: white;
          padding: 0.25rem 0.5rem;
          border-radius: 4px;
          font-size: 0.9rem;
        }
      }
    }
  }

  .kifu-view-placeholder {
    background: white;
    border: 2px dashed var(--border-color);
    border-radius: 8px;
    padding: 4rem;
    margin-bottom: 2rem;
    text-align: center;
    color: var(--text-color);
  }

  .action-buttons {
    margin-bottom: 2rem;

    .like-button {
      padding: 0.5rem 1rem;
      background-color: white;
      border: 1px solid var(--secondary-color);
      color: var(--secondary-color);
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s;

      &.liked {
        background-color: var(--secondary-color);
        color: white;
      }

      &:hover {
        opacity: 0.8;
      }
    }
  }

  .comments-section {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    h2 {
      margin-bottom: 1.5rem;
      color: var(--primary-color);
    }

    .comment-form {
      margin-bottom: 2rem;

      textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid var(--border-color);
        border-radius: 4px;
        resize: vertical;
        min-height: 100px;
        margin-bottom: 1rem;
        cursor: text;

        &:focus {
          border-color: var(--secondary-color);
        }
      }

      button {
        padding: 0.5rem 1.5rem;
        background-color: var(--primary-color);
        color: white;
        border-radius: 4px;
        cursor: pointer;

        &:disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }

        &:not(:disabled):hover {
          opacity: 0.9;
        }
      }
    }

    .comments-list {
      .no-comments {
        text-align: center;
        color: var(--text-color);
        font-style: italic;
      }

      .comment {
        padding: 1rem;
        border-bottom: 1px solid var(--border-color);

        &:last-child {
          border-bottom: none;
        }

        .comment-header {
          display: flex;
          justify-content: space-between;
          margin-bottom: 0.5rem;
          font-size: 0.9rem;

          .comment-author {
            font-weight: bold;
            color: var(--primary-color);
          }

          .comment-date {
            color: var(--text-color);
            opacity: 0.8;
          }
        }

        .comment-content {
          color: var(--text-color);
          white-space: pre-wrap;
        }
      }
    }
  }
</style>
