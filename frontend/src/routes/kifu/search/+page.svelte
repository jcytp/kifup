<!-- src/routes/kifu/search/+page.svelte -->

<script lang="ts">
  import type { Kifu } from '$lib/types/Kifu';
  import { onMount } from 'svelte';

  // 検索条件の状態管理
  let keyword = '';
  let accountId = '';
  let tags: string[] = [];
  let startDate = '';
  let endDate = '';

  // ページネーションの状態管理
  let currentPage = 1;
  const itemsPerPage = 10;
  let totalPages = 1;

  // 棋譜リストの状態管理
  let kifuList: Kifu[] = [];
  let isLoading = false;

  // タグの入力管理
  let tagInput = '';
  
  // タグの追加
  function addTag() {
    if (tagInput && !tags.includes(tagInput)) {
      tags = [...tags, tagInput];
      tagInput = '';
    }
  }

  // タグの削除
  function removeTag(index: number) {
    tags = tags.filter((_, i) => i !== index);
  }

  // 検索実行
  async function handleSearch() {
    isLoading = true;
    // TODO: APIリクエストの実装
    // 仮のデータを表示
    kifuList = Array(itemsPerPage).fill(null).map((_, i) => ({
      id: `kifu-${i}`,
      ownerId: 'user-1',
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
    isLoading = false;
  }

  // ページ変更
  function changePage(page: number) {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
      handleSearch();
    }
  }

  onMount(() => {
    handleSearch();
  });
</script>

<div class="container">
  <section class="search-section">
    <h2>棋譜検索</h2>
    <form on:submit|preventDefault={handleSearch} class="search-form">
      <div class="form-group">
        <h3 class="search-label">キーワード</h3>
        <input
          type="text"
          id="keyword"
          bind:value={keyword}
          placeholder="タイトル、対局者名で検索"
        />
      </div>

      <div class="form-group">
        <h3 class="search-label">投稿アカウント</h3>
        <input
          type="text"
          id="accountId"
          bind:value={accountId}
          placeholder="アカウント名で検索"
        />
      </div>

      <div class="form-group">
        <h3 class="search-label">タグ</h3>
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
          {#each tags as tag, i}
            <span class="tag">
              {tag}
              <button type="button" on:click={() => removeTag(i)} class="remove-tag">×</button>
            </span>
          {/each}
        </div>
      </div>

      <div class="form-group date-range">
        <h3 class="search-label">対局日</h3>
        <div class="date-inputs">
          <input type="date" bind:value={startDate} />
          <span>～</span>
          <input type="date" bind:value={endDate} />
        </div>
      </div>

      <button type="submit" class="search-button">検索</button>
    </form>
  </section>

  <section class="results-section">
    {#if isLoading}
      <div class="loading">検索中...</div>
    {:else}
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
    {/if}
  </section>
</div>

<style lang="scss">
  .container {
    max-width: 1000px;
    margin: 0 auto;
    padding: 2rem;
  }

  .search-section {
    margin-bottom: 2rem;

    h2 {
      margin-bottom: 1rem;
      color: var(--primary-color);
    }
  }

  .search-form {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .form-group {
    margin-bottom: 1rem;

    .search-label {
      display: block;
      margin-bottom: 0.5rem;
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

  .date-range {
    .date-inputs {
      display: flex;
      gap: 1rem;
      align-items: center;

      input {
        width: calc(50% - 1rem);
      }
    }
  }

  .tag-input-container {
    display: flex;
    gap: 0.5rem;

    input {
      flex: 1;
    }
  }

  .add-tag-btn {
    padding: 0.5rem 1rem;
    background-color: var(--secondary-color);
    color: white;
    border-radius: 4px;
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 0.5rem;
  }

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
    }
  }

  .search-button {
    width: 100%;
    padding: 0.75rem;
    background-color: var(--primary-color);
    color: white;
    border-radius: 4px;
    font-weight: bold;
    margin-top: 1rem;

    &:hover {
      opacity: 0.9;
    }
  }

  .results-section {
    .loading {
      text-align: center;
      padding: 2rem;
      color: var(--text-color);
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
