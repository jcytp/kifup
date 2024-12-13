<!-- src/lib/components/KifuList.svelte -->

<script lang="ts">
  import type { PaginationResponse } from '$lib/types/API';
  import type { KifuSummary } from '$lib/types/Kifu';

  // Callback functions
  export let changePage: (page: number) => Promise<void>;
  export let togglePublic: (kifu: KifuSummary) => Promise<void> = async (kifu: KifuSummary) => {};
  export let deleteKifu: (kifuId: string) => Promise<void> = async (kifuId: string) => {};

  // Props
  export let kifuList: KifuSummary[] = [];
  export let pagination: PaginationResponse | null = null;
  export let loading = false;
  export let mode: 'view-only' | 'with-actions' = 'view-only';

  // 選択中の棋譜ID（アクション表示用）
  let selectedKifuId: string | null = null;

  // ページ変更ハンドラ
  function handleChangePage(page: number) {
    if (pagination && page >= 1 && page <= pagination.max_page) {
      changePage(page);
      selectedKifuId = null; // ページ変更時にポップアップを閉じる
    }
  }

  // イベントハンドラ（MouseEventとKeyboardEventの両方に対応）
  function handleCardClick(event: MouseEvent | KeyboardEvent, kifuId: string) {
    if (mode !== 'with-actions') return;

    // ポップアップメニュー内のクリックはカード自体のクリックイベントを発火させない
    if ((event.target as HTMLElement).closest('.popup-menu')) {
      return;
    }
    selectedKifuId = selectedKifuId === kifuId ? null : kifuId;
  }

  // アクションハンドラ
  function handleTogglePublic(kifu: KifuSummary) {
    togglePublic(kifu);
    selectedKifuId = null;
  }

  function handleDelete(kifuId: string) {
    deleteKifu(kifuId);
    selectedKifuId = null;
  }
</script>

{#if loading}
  <div class="loading">読み込み中...</div>
{:else if kifuList.length === 0}
  <div class="loading">棋譜はありません。</div>
{:else}
  <div class="kifu-list-container">
    {#each kifuList as kifu}
      {#if mode === 'with-actions'}
        <!-- アクション付きカード（ホームページ用） -->
        <div class="kifu-list-item">
          <button
            type="button"
            class="card kifu-card"
            class:selected={selectedKifuId === kifu.id}
            on:click={(e) => handleCardClick(e, kifu.id)}
            on:keydown={(e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                handleCardClick(e, kifu.id);
              }
            }}
            aria-expanded={selectedKifuId === kifu.id}
            aria-haspopup="menu"
            aria-controls={`menu-${kifu.id}`}
          >
            <div class="kifu-header">
              <h3>{kifu.title}</h3>
              <span class={`kifu-status ${kifu.is_public ? 'public' : 'private'}`}>
                {kifu.is_public ? '公開' : '非公開'}
              </span>
            </div>
            <div class="kifu-info">
              <span>対局日: {kifu.game_info.date}</span>
              <span>先手: {kifu.game_info.black}</span>
              <span>後手: {kifu.game_info.white}</span>
            </div>
            <div class="kifu-tags">
              {#each kifu.tags as tag}
                <span class="tag">{tag}</span>
              {/each}
            </div>
          </button>

          {#if selectedKifuId === kifu.id}
            <div id={`menu-${kifu.id}`} class="popup-menu" role="menu">
              <a href={`/kifu/view?id=${kifu.id}`} class="menu-item" role="menuitem">
                詳細を見る
              </a>
              <a href={`/kifu/edit?id=${kifu.id}`} class="menu-item" role="menuitem"> 編集する </a>
              <button
                type="button"
                class="menu-item"
                role="menuitem"
                on:click={() => handleTogglePublic(kifu)}
              >
                {kifu.is_public ? '非公開にする' : '公開する'}
              </button>
              <button
                type="button"
                class="menu-item delete"
                role="menuitem"
                on:click={() => handleDelete(kifu.id)}
              >
                削除する
              </button>
            </div>
          {/if}
        </div>
      {:else}
        <!-- 閲覧専用カード（検索・アカウントページ用） -->
        <a href={`/kifu/view?id=${kifu.id}`} class="card kifu-card">
          <div class="kifu-header">
            <h3>{kifu.title}</h3>
          </div>
          <div class="kifu-info">
            <span>対局日: {kifu.game_info.date}</span>
            <span>対局者: {kifu.game_info.black} vs {kifu.game_info.white}</span>
          </div>
          <div class="kifu-tags">
            {#each kifu.tags as tag}
              <span class="tag">{tag}</span>
            {/each}
          </div>
        </a>
      {/if}
    {/each}
  </div>

  <!-- ページネーション -->
  {#if pagination}
    <div class="pagination">
      <button
        class="page-button"
        disabled={pagination.page === 1}
        on:click={() => handleChangePage(pagination.page - 1)}
      >
        前へ
      </button>

      {#each Array(pagination.max_page) as _, i}
        <button
          class="page-button"
          class:active={pagination.page === i + 1}
          on:click={() => handleChangePage(i + 1)}
        >
          {i + 1}
        </button>
      {/each}

      <button
        class="page-button"
        disabled={pagination.page === pagination.max_page}
        on:click={() => handleChangePage(pagination.page + 1)}
      >
        次へ
      </button>
    </div>
  {/if}
{/if}

<style lang="scss">
  .loading {
    text-align: center;
    padding: 2rem;
    color: #666;
  }

  .kifu-list-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    .kifu-list-item {
      position: relative;
    }

    .kifu-card {
      display: block;
      width: 100%;
      padding: 1rem 1.5rem;
      transition:
        transform 0.2s,
        box-shadow 0.2s;

      .kifu-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .kifu-status {
          font-size: 0.8rem;
          padding: 0.25rem 0.5rem;
          border-radius: 0.2rem;

          &.public {
            background-color: var(--public-icon-background-color);
            color: var(--public-icon-text-color);
          }

          &.private {
            background-color: var(--private-icon-background-color);
            color: var(--private-icon-text-color);
          }
        }
      }

      .kifu-info {
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-bottom: 0.3rem;
        font-size: 0.9rem;
        color: #666;
      }

      .kifu-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        justify-content: end;

        .tag {
          background-color: var(--secondary-color);
          color: white;
          padding: 0.2rem 0.4rem;
          border-radius: 0.2rem;
          font-size: 0.8rem;
        }
      }
    }

    .popup-menu {
      position: absolute;
      top: -1rem;
      right: 1rem;
      background: var(--pickup-color);
      border-radius: 0.4rem;
      box-shadow: 0 0 0.8rem #696;
      z-index: var(--z-index-popup);
      min-width: 12rem;

      .menu-item {
        display: block;
        padding: 0.4rem 0.8rem;
        width: 100%;
        text-align: left;
        font-size: 1rem;

        &:hover {
          background-color: var(--pickup-strong-color);
        }

        &.delete {
          color: #e53e3e;
          border-top: 1px solid var(--border-color);
        }
      }
    }
  }
</style>
