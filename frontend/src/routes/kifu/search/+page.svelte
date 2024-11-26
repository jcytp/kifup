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
    kifuList = Array(itemsPerPage)
      .fill(null)
      .map((_, i) => ({
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

<div class="page">
  <section class="basic">
    <h2>棋譜検索</h2>
    <form on:submit|preventDefault={handleSearch} class="basic search-form">
      <div class="form-group">
        <h3 class="label">キーワード</h3>
        <input
          type="text"
          id="keyword"
          bind:value={keyword}
          placeholder="タイトル、対局者名で検索"
        />
      </div>

      <div class="form-group">
        <h3 class="label">投稿アカウント</h3>
        <input type="text" id="accountId" bind:value={accountId} placeholder="アカウント名で検索" />
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
          {#each tags as tag, i}
            <span class="tag">
              {tag}
              <button type="button" on:click={() => removeTag(i)} class="remove-tag">×</button>
            </span>
          {/each}
        </div>
      </div>

      <div class="form-group">
        <h3 class="label">対局日</h3>
        <div class="match-date-inputs">
          <input type="date" bind:value={startDate} />
          <span>～</span>
          <input type="date" bind:value={endDate} />
        </div>
      </div>

      <button type="submit" class="submit search-button">検索</button>
    </form>
  </section>

  <section class="basic kifu-list">
    {#if isLoading}
      <div class="loading">検索中...</div>
    {:else}
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

      <!-- ページネーション -->
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
  .form-group {
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

    .match-date-inputs {
      display: flex;
      gap: 1rem;
      align-items: center;
    }
  }

  section.kifu-list {
    border-top: 0.2rem solid var(--primary-color);

    .loading {
      text-align: center;
      padding: 2rem;
      color: var(--text-color);
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
  }
</style>
