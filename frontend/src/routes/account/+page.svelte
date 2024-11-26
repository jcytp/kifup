<!-- src/routes/account/+page.svelte -->

<script lang="ts">
  import { page } from '$app/stores';
  import type { Account } from '$lib/types/Account';
  import type { Kifu } from '$lib/types/Kifu';

  // URLからアカウントIDを取得
  const accountId = $page.url.searchParams.get('id');

  const image_url_base = 'http://example.com/icon/';

  // アカウント情報の状態管理
  let accountInfo: Account | null = null;
  let kifuList: Kifu[] = [];
  let isLoading = true;
  let error: string | null = null;

  // ページネーションの状態管理
  let currentPage = 1;
  const itemsPerPage = 10;
  let totalPages = 1;

  // アカウント情報とその棋譜リストを取得
  async function fetchAccountData() {
    isLoading = true;
    try {
      // TODO: API実装後に実際のデータ取得に置き換え
      await new Promise((resolve) => setTimeout(resolve, 500)); // ローディング表示確認用

      accountInfo = {
        id: accountId || '',
        name: 'サンプルユーザー',
        introduction: '自己紹介～～～～～～～～～～～～～～',
        created_at: new Date(),
        last_login_at: new Date(),
      };

      kifuList = Array(itemsPerPage)
        .fill(null)
        .map((_, i) => ({
          id: `kifu-${i}`,
          ownerId: accountId || '',
          title: `テスト棋譜 ${i + 1}`,
          matchInfo: {
            black: '先手太郎',
            white: '後手次郎',
            date: '2024-01-01',
          },
          tags: ['実戦', 'テスト'],
          isPublic: true,
          moves: [],
        }));

      totalPages = 5;
    } catch (e) {
      error = 'アカウント情報の取得に失敗しました。';
    } finally {
      isLoading = false;
    }
  }

  // ページ変更
  function changePage(page: number) {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
      fetchAccountData();
    }
  }

  // 初期データ取得
  $: if (accountId) {
    fetchAccountData();
  }
</script>

<div class="page">
  {#if isLoading}
    <div class="loading">
      <p>アカウント情報を読み込んでいます...</p>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
    </div>
  {:else if accountInfo}
    <section class="basic">
      <h2>ユーザー情報</h2>
      <div class="card account-profile">
        <div class="profile-header">
          <div class="profile-image">
            {#if accountInfo.icon_id}
              <img src={`${image_url_base}${accountInfo.icon_id}`} alt="プロフィール画像" />
            {:else}
              <div class="placeholder-image">
                {accountInfo.name[0]}
              </div>
            {/if}
          </div>
          <div class="profile-info">
            <p>{accountInfo.name}</p>
            <p>{accountInfo.introduction}</p>
          </div>
        </div>
      </div>
    </section>

    <hr />

    <section class="basic">
      <h2>公開棋譜</h2>

      <div class="kifu-list-container">
        {#each kifuList as kifu}
          <a href={`/kifu/view?id=${kifu.id}`} class="card kifu-card">
            <div class="kifu-header">
              <h3>{kifu.title}</h3>
            </div>
            <div class="kifu-info">
              <span>対局日: {kifu.matchInfo.date}</span>
              <span>対局者: {kifu.matchInfo.black} vs {kifu.matchInfo.white}</span>
            </div>
            <div class="kifu-tags">
              {#each kifu.tags as tag}
                <span class="tag">{tag}</span>
              {/each}
            </div>
          </a>
        {/each}
      </div>

      <div class="pagination">
        <button
          class="page-button"
          disabled={currentPage === 1}
          on:click={() => changePage(currentPage - 1)}
        >
          前へ
        </button>

        {#each Array(totalPages) as _, i}
          <button
            class="page-button"
            class:active={currentPage === i + 1}
            on:click={() => changePage(i + 1)}
          >
            {i + 1}
          </button>
        {/each}

        <button
          class="page-button"
          disabled={currentPage === totalPages}
          on:click={() => changePage(currentPage + 1)}
        >
          次へ
        </button>
      </div>
    </section>
  {/if}
</div>

<style lang="scss">
  .account-profile {
    .profile-header {
      display: flex;
      gap: 2rem;
      align-items: flex-start;

      .profile-image {
        flex: 0 0 6rem;
        height: 6rem;
        overflow: hidden;
        border-radius: 50%;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }

        .placeholder-image {
          width: 100%;
          height: 100%;
          background-color: var(--secondary-color);
          color: white;
          font-size: 4rem;
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }

      .profile-info {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
      }
    }
  }

  .kifu-list-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;

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
  }

  .pagination {
    display: flex;
    justify-content: center;
    gap: 0.5rem;
    margin-top: 2rem;

    .page-button {
      padding: 0.5rem 1rem;
      border: 1px solid var(--border-color);
      border-radius: 4px;
      background: white;

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }

      &.active {
        background-color: var(--primary-color);
        color: white;
        border-color: var(--primary-color);
      }

      &:not(:disabled):hover {
        background-color: var(--secondary-color);
        color: white;
        border-color: var(--secondary-color);
      }
    }
  }
</style>
