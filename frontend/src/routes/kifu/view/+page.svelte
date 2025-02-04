<!-- src/routes/kifu/view/+page.svelte -->

<script lang="ts">
  import type { KifuComment, KifuDetail, KifuMove } from '$lib/types/Kifu';
  import { page } from '$app/stores';
  import { account } from '$lib/stores/session';
  import { getKifu } from '$lib/apis/kifu';
  import { addKifuLike, removeKifuLike, getKifuComments, addKifuComment } from '$lib/apis/social';
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
  let newComment = '';
  let commentsError = false;

  const fetchKifuComments = async () => {
    commentsError = false;
    if (!kifuId) return;

    const result = await getKifuComments(kifuId);
    if (result.ok && result.data) {
      comments = result.data.map((comment: KifuComment) => ({
        id: comment.id,
        userId: comment.account.id,
        userName: comment.account.name,
        content: comment.content,
        createdAt: formatDateTime(comment.created_at),
      }));
    } else {
      console.error('Failed to fetch comments: ', result);
      commentsError = true;
    }
  };

  const handleCommentSubmit = async (e: Event) => {
    e.preventDefault();
    if (!$account || !kifuId || !newComment.trim()) return;

    const result = await addKifuComment(kifuId, newComment.trim());
    if (result.ok) {
      // コメント投稿成功後、コメント一覧を再取得
      await fetchKifuComments();
      newComment = '';
    } else {
      console.error('Failed to post comment: ', result);
    }
  };

  const toggleLike = async () => {
    if (!$account || !kifuId) return;

    const result = await (kifu.has_like ? removeKifuLike(kifuId) : addKifuLike(kifuId));
    if (result.ok) {
      kifu.has_like = !kifu.has_like;
      kifu.like_count += kifu.has_like ? 1 : -1;
    } else {
      console.error('Failed to toggle like: ', result);
    }
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
      <button
        class="like-button"
        class:liked={kifu.has_like}
        onclick={toggleLike}
        disabled={!$account}
      >
        {kifu.has_like ? '★' : '☆'} いいね！ ({kifu.like_count})
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

      {#if $account}
        <form class="basic" onsubmit={handleCommentSubmit}>
          <div class="form-group">
            <label for="comment-input">コメント</label>
            <textarea id="comment-input" bind:value={newComment} placeholder="コメントを入力..."
            ></textarea>
          </div>
          <button type="submit" class="submit">コメントを送信</button>
        </form>
      {/if}
    </section>
  {/if}
</div>

<style lang="scss">
  @import '../../../lib/styles/mixins.scss';

  .kifu-info {
    display: flex;
    flex-wrap: wrap;

    .kifu-info-block {
      @include sp {
        min-width: 100%;
      }
      @include pc {
        min-width: 35%;
      }
      flex: 1;
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
