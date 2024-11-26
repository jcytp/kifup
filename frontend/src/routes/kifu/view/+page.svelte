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
      await new Promise((resolve) => setTimeout(resolve, 500)); // ローディング表示確認用

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
            piece: '歩',
            from: { x: 7, y: 7 },
            to: { x: 7, y: 6 },
            comment: '歩を突く',
            variations: [
              {
                moveNumber: 1,
                piece: '歩',
                from: { x: 5, y: 7 },
                to: { x: 5, y: 6 },
                comment: '５六歩と指す作戦もある',
              },
            ],
          },
          {
            moveNumber: 2,
            piece: '歩',
            from: { x: 3, y: 3 },
            to: { x: 3, y: 4 },
          },
          {
            moveNumber: 3,
            piece: '歩',
            from: { x: 6, y: 7 },
            to: { x: 6, y: 6 },
          },
        ],
        createdAt: '2024-01-01 12:00',
        updatedAt: '2024-01-01 15:30',
      };

      comments = [
        {
          id: '1',
          userId: 'user-2',
          userName: '観戦者A',
          content: '素晴らしい一手でした！',
          createdAt: '2024-01-01 12:00',
        },
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
        createdAt: formatDateTime(new Date().toISOString()),
      },
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

<div class="page">
  {#if kifu}
    <section class="basic">
      <h2>{kifu.title}</h2>
      <form class="basic kifu-info">
        <p>先手： {kifu.matchInfo.black}</p>
        <p>後手： {kifu.matchInfo.white}</p>
        <p>対局日時： {kifu.matchInfo.date}</p>
        <p>対局場所： {kifu.matchInfo.place}</p>
        <p>
          持ち時間： {kifu.matchInfo.timeLimit?.initial || 0}分, 秒読み {kifu.matchInfo.timeLimit
            ?.byoyomi || 0}秒, 加算 {kifu.matchInfo.timeLimit?.increment || 0}秒
        </p>
        <p class="tags">
          {#each kifu.tags as tag}
            <span class="tag">{tag}</span>
          {/each}
        </p>
        <p>
          アップロード： <a href={`/account?id=${kifu.ownerId}`} class="owner-link">ユーザー名</a>
        </p>
        {#if isOwner}
          <p class="edit-link"><a href={`/kifu/edit?id=${kifu.id}`}>編集する</a></p>
        {/if}
      </form>

      <KifuPlayer {kifu} />
    </section>

    <section class="basic">
      <button class="like-button" class:liked={isLiked} on:click={toggleLike}>
        {isLiked ? '★' : '☆'} いいね！ ({likeCount})
      </button>

      <div class="comments-list">
        {#if comments.length === 0}
          <p class="no-comments">コメントはまだありません。</p>
        {:else}
          {#each comments as comment}
            <div class="card comment">
              <div class="comment-header">
                <span class="comment-author">{comment.userName}</span>
                <span class="comment-date">{comment.createdAt}</span>
              </div>
              <p class="comment-content">{comment.content}</p>
            </div>
          {/each}
        {/if}
      </div>

      <form class="basic">
        <div class="form-group">
          <label for="comment-input">コメント</label>
          <textarea id="comment-input" bind:value={newComment} placeholder="コメントを入力..."
          ></textarea>
        </div>
        <button type="submit" class="submit">コメントを送信</button>
      </form>
    </section>
  {/if}
</div>

<style lang="scss">
  .owner-link {
    color: var(--secondary-color);
    text-decoration: underline;

    &:hover {
      opacity: 0.8;
    }
  }

  .edit-link {
    margin-top: 1.5rem;

    a {
      display: block;
      background-color: var(--primary-color);
      color: var(--background-color);
      padding: 0.5rem 1rem;
      border: none;
      border-radius: 0.5rem;
      width: 100%;
      text-align: center;

      &:hover {
        background-color: var(--secondary-color);
      }
    }
  }

  .like-button {
    padding: 0.4rem 1rem;
    background-color: white;
    border: 1px solid var(--secondary-color);
    color: var(--secondary-color);
    border-radius: 0.4rem;
    line-height: 1.2rem;
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

  .comments-list {
    margin: 1rem 0;

    .no-comments {
      text-align: center;
      color: var(--text-color);
      font-style: italic;
    }

    .comment {
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
</style>
