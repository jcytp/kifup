<!-- src/routes/kifu/view/+page.svelte -->

<script lang="ts">
  import type { KifuDetail, KifuMove } from '$lib/types/Kifu';
  import { page } from '$app/stores';
  import { account } from '$lib/stores/session';
  import { getKifu } from '$lib/apis/kifu';
  import { formatDateTime, formatTimeRule } from '$lib/utils/textFormat';
  import KifuPlayer from '$lib/components/KifuPlayer.svelte';

  // ----------------------------------------
  // 棋譜データの状態管理

  const kifuId = $page.url.searchParams.get('id'); // URLからkifuIDを取得

  let kifu: KifuDetail;
  let kifuError = false;
  $: isOwner = kifu?.owner.id === $account?.id; // 所有者本人かどうか
  let timeRuleString = '--';
  let moves: KifuMove[] = []; // 表示する指し手のライン

  const fetchKifuData = async () => {
    kifuError = false;
    if (!kifuId) {
      console.error('No kifuId');
      kifuError = true;
      return;
    }

    const result = await getKifu(kifuId, $account ? true : false);
    if (result.ok && result.data) {
      kifu = result.data as KifuDetail;
    } else {
      console.error('Failed to fetch kifu detail: ', result);
      kifuError = true;
      return;
    }

    timeRuleString =
      formatTimeRule(kifu.game_info.持ち時間, kifu.game_info.秒読み, kifu.game_info.秒加算) || '--';
    moves = kifu.moves;
  };

  // ----------------------------------------
  // コメント・いいねの状態管理

  let comments: Array<{
    id: string;
    userId: string;
    userName: string;
    content: string;
    createdAt: string;
  }> = [];
  let likeCount = 0;
  let isLiked = false;
  let newComment = '';
  let commentsError = false;

  const fetchKifuComments = async () => {
    commentsError = false;

    // ToDo: API実装後にデータ取得を実装
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
    isLiked = false;
  };

  const handleCommentSubmit = async () => {
    if (!newComment.trim()) return;

    // ToDo: API実装後にコメント投稿処理を実装
    comments = [
      ...comments,
      {
        id: String(Date.now()),
        userId: $account ? $account.id : '',
        userName: $account ? $account.name : 'anonymous',
        content: newComment,
        createdAt: formatDateTime(new Date()),
      },
    ];

    newComment = '';
  };

  const toggleLike = async () => {
    // ToDo: API実装後にいいね処理を実装
    isLiked = !isLiked;
    likeCount += isLiked ? 1 : -1;
  };

  // ----------------------------------------
  let preinit = true;
  $: if (preinit) {
    preinit = false;
    fetchKifuData();
    fetchKifuComments();
  }
  $: if ($account) {
    preinit = false;
    fetchKifuData();
    fetchKifuComments();
  }
</script>

<div class="page">
  {#if kifu}
    <section class="basic">
      <h2>{kifu.title}</h2>
      <form class="basic kifu-info">
        <div class="kifu-info-block">
          <p>先手： {kifu.game_info.先手}</p>
          <p>後手： {kifu.game_info.後手}</p>
        </div>
        <div class="kifu-info-block">
          <p>対局日時： {formatDateTime(kifu.game_info.対局日時)}</p>
          <p>対局場所： {kifu.game_info.対局場所}</p>
          <p>持ち時間： {timeRuleString}</p>
        </div>
        <div class="kifu-info-block">
          <p>
            作成： <a href={`/account?id=${kifu.owner.id}`} class="owner-link">{kifu.owner.name}</a>
          </p>
        </div>
        <div class="kifu-info-block">
          <p class="tags">
            {#each kifu.tags as tag}
              <span class="tag">{tag}</span>
            {/each}
          </p>
        </div>
        {#if isOwner}
          <div class="kifu-info-block">
            <p class="edit-link"><a href={`/kifu/edit?id=${kifu.id}`}>編集する</a></p>
          </div>
        {/if}
      </form>

      <KifuPlayer initialSfen={kifu.initial_position} moveList={moves} />
    </section>

    <section class="basic">
      <button class="like-button" class:liked={isLiked} onclick={toggleLike}>
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

      <form class="basic" onsubmit={handleCommentSubmit}>
        <div class="form-group">
          <label for="comment-input">コメント</label>
          <textarea id="comment-input" bind:value={newComment} placeholder="コメントを入力..."
          ></textarea>
        </div>
        <button type="submit" class="submit">{$account ? '' : '匿名で'}コメントを送信</button>
      </form>
    </section>
  {/if}
</div>

<style lang="scss">
  .kifu-info {
    display: flex;
    flex-wrap: wrap;

    .kifu-info-block {
      flex: 1;
      min-width: 35%;
    }

    .tags {
      display: flex;
      flex-wrap: wrap;
      gap: 0.4rem;
      margin: 0.5rem 0 0 0;

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
      }
    }

    .owner-link {
      color: var(--secondary-color);
      text-decoration: underline;

      &:hover {
        opacity: 0.8;
      }
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
