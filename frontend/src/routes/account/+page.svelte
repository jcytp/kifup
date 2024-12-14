<!-- src/routes/account/+page.svelte -->

<script lang="ts">
  import { page } from '$app/stores';
  import { getAccountById } from '$lib/apis/account';
  import { searchKifus } from '$lib/apis/kifu';
  import KifuList from '$lib/components/KifuList.svelte';
  import { account } from '$lib/stores/session';
  import type { Account } from '$lib/types/Account';
  import type { PaginationResponse } from '$lib/types/API';
  import type { KifuSummary } from '$lib/types/Kifu';
  import { onMount } from 'svelte';

  // ----------------------------------------
  // アカウント情報
  const accountId = $page.url.searchParams.get('id');
  const image_url_base = 'http://example.com/icon/';
  let isLoadingAccount = true;
  let accountInfo: Account;

  const fetchAccountData = async () => {
    isLoadingAccount = true;

    if (accountId) {
      const result = await getAccountById(accountId);
      if (result.ok && result.data) {
        accountInfo = result.data as Account;
      } else {
        console.error('Failed to fetch account data:', result);
      }
    }

    isLoadingAccount = false;
  };

  // ----------------------------------------
  // 棋譜リスト
  let isLoadingKifuList = true;
  let kifuList: KifuSummary[] = [];
  let pagination: PaginationResponse = {
    total_count: 0,
    page: 1,
    page_size: 10,
    max_page: 1,
  };

  const changePage = async (page: number) => {
    pagination.page = page;
    await fetchKifuList();
  };

  const fetchKifuList = async () => {
    isLoadingKifuList = true;

    const result = await searchKifus(accountId, pagination.page, pagination.page_size, false);
    if (result.ok && result.data && result.pagination) {
      kifuList = result.data as KifuSummary[];
      pagination = result.pagination;
    } else {
      kifuList = []; // エラー時はリストをクリア
      console.error('Failed to fetch kifu list: ', result);
    }

    isLoadingKifuList = false;
  };

  // ----------------------------------------
  // 初回データロード

  onMount(() => {
    fetchAccountData();
    fetchKifuList();
  });
</script>

<div class="page">
  <section class="basic">
    {#if isLoadingAccount}
      <div class="loading">
        <p>アカウント情報を読み込んでいます...</p>
      </div>
    {:else if accountInfo}
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
    {:else}
      <div class="error">
        <p>アカウント情報の取得に失敗しました。</p>
      </div>
    {/if}
  </section>

  <hr />

  <section class="basic">
    <h2>公開棋譜</h2>
    <KifuList {kifuList} {pagination} loading={isLoadingKifuList} mode="view-only" {changePage} />
  </section>
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
</style>
