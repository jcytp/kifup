<!-- src/routes/account/+page.svelte -->

<script lang="ts">
  import { page } from '$app/stores';
  import type { Account } from '$lib/types/Account';
  import type { Kifu } from '$lib/types/Kifu';

  // URLからアカウントIDを取得
  const accountId = $page.url.searchParams.get('id');

  // アカウント情報の状態管理
  let account: Account | null = null;
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
      await new Promise(resolve => setTimeout(resolve, 500)); // ローディング表示確認用
      
      account = {
        id: accountId || '',
        name: 'サンプルユーザー',
        email: 'sample@example.com'
      };

      kifuList = Array(itemsPerPage).fill(null).map((_, i) => ({
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

<div class="container">
  {#if isLoading}
    <div class="loading">
      <p>アカウント情報を読み込んでいます...</p>
      <a href="/kifu/view" class="back-button">棋譜詳細に戻る</a>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
      <a href="/kifu/view" class="back-button">棋譜詳細に戻る</a>
    </div>
  {:else if account}
    <div class="account-profile">
      <div class="profile-header">
        <div class="profile-image">
          <!-- TODO: プロフィール画像の実装 -->
          <div class="placeholder-image">
            {account.name[0]}
          </div>
        </div>
        <div class="profile-info">
          <h1>{account.name}</h1>
          <div class="profile-meta">
            <p>ID: {account.id}</p>
          </div>
        </div>
      </div>

      <div class="profile-bio">
        <h2>自己紹介</h2>
        <p>
          <!-- TODO: 自己紹介文の実装 -->
          自己紹介文がまだ設定されていません。
        </p>
      </div>
    </div>

    <div class="kifu-section">
      <h2>公開棋譜</h2>
      
      <div class="kifu-list">
        {#each kifuList as kifu}
          <a href={`/kifu/view?id=${kifu.id}`} class="kifu-card">
            <h3>{kifu.title}</h3>
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
    </div>
  {/if}
</div>

<style lang="scss">
  .container {
    max-width: 1000px;
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

  .account-profile {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;

    .profile-header {
      display: flex;
      gap: 2rem;
      margin-bottom: 2rem;

      .profile-image {
        flex: 0 0 120px;
        height: 120px;

        .placeholder-image {
          width: 100%;
          height: 100%;
          background-color: var(--secondary-color);
          color: white;
          font-size: 3rem;
          display: flex;
          align-items: center;
          justify-content: center;
          border-radius: 50%;
        }
      }

      .profile-info {
        flex: 1;

        h1 {
          margin-bottom: 0.5rem;
          color: var(--primary-color);
        }

        .profile-meta {
          color: var(--text-color);
          opacity: 0.8;
        }
      }
    }

    .profile-bio {
      h2 {
        color: var(--primary-color);
        margin-bottom: 1rem;
        font-size: 1.2rem;
      }

      p {
        color: var(--text-color);
        opacity: 0.8;
        font-style: italic;
      }
    }
  }

  .kifu-section {
    h2 {
      color: var(--primary-color);
      margin-bottom: 1.5rem;
    }
  }

  .kifu-list {
    display: grid;
    gap: 1rem;
  }

  .kifu-card {
    background: white;
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s;
    cursor: pointer;

    &:hover {
      transform: translateY(-2px);
    }

    h3 {
      margin-bottom: 0.5rem;
      color: var(--primary-color);
    }

    .kifu-info {
      display: flex;
      gap: 1rem;
      margin-bottom: 0.5rem;
      color: var(--text-color);
      font-size: 0.9rem;
    }

    .kifu-tags {
      display: flex;
      gap: 0.5rem;

      .tag {
        background-color: var(--secondary-color);
        color: white;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.9rem;
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
