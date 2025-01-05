<!-- src/routes/kifu/search/+page.svelte -->

<script lang="ts">
  import { searchKifus } from '$lib/apis/kifu';
  import KifuList from '$lib/components/KifuList.svelte';
  import type { PaginationResponse } from '$lib/types/API';
  import type { KifuSummary } from '$lib/types/Kifu';
  import { onMount } from 'svelte';

  // ----------------------------------------
  // 棋譜リスト
  let loading = true;
  let kifuList: KifuSummary[] = [];
  let pagination: PaginationResponse = {
    total_count: 0,
    page: 1,
    page_size: 10,
    max_page: 1,
  };

  const changePage = async (page: number) => {
    pagination.page = page;
    handleSearch();
  };

  // ----------------------------------------
  // 検索条件
  let keyword = '';
  let accountId = '';
  let tags: string[] = [];
  let startDate = '';
  let endDate = '';

  // タグ
  let tagInput = '';
  function addTag() {
    if (tagInput && !tags.includes(tagInput)) {
      tags = [...tags, tagInput];
      tagInput = '';
    }
  }
  function removeTag(index: number) {
    tags = tags.filter((_, i) => i !== index);
  }

  // ----------------------------------------
  // 検索実行
  const handleSearch = async () => {
    loading = true;

    const result = await searchKifus(null, pagination.page, pagination.page_size, false);
    if (result.ok && result.data && result.pagination) {
      kifuList = result.data as KifuSummary[];
      pagination = result.pagination;
    } else {
      kifuList = []; // エラー時はリストをクリア
      console.error('Failed to fetch kifu list: ', result);
    }

    loading = false;
  };

  // ----------------------------------------
  // 初回データロード

  onMount(() => {
    handleSearch();
  });
</script>

<div class="page">
  <section class="basic">
    <h2>棋譜検索</h2>
    <form on:submit|preventDefault={handleSearch} class="basic search-form">
      <!-- <div class="form-group">
        <h3 class="label">キーワード</h3>
        <input
          type="text"
          id="keyword"
          bind:value={keyword}
          placeholder="タイトル、対局者名で検索"
        />
      </div> -->

      <div class="form-group">
        <h3 class="label">投稿アカウント</h3>
        <input type="text" id="accountId" bind:value={accountId} placeholder="アカウント名で検索" />
      </div>

      <!-- <div class="form-group">
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
      </div> -->

      <!-- <div class="form-group">
        <h3 class="label">対局日</h3>
        <div class="match-date-inputs">
          <input type="date" bind:value={startDate} />
          <span>～</span>
          <input type="date" bind:value={endDate} />
        </div>
      </div> -->

      <button type="submit" class="submit search-button">検索</button>
    </form>
  </section>

  <section class="basic kifu-list">
    <KifuList {kifuList} {pagination} {loading} mode="view-only" {changePage} />
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
</style>
