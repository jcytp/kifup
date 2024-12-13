<!-- src/routes/home/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import type { KifuSummary } from '$lib/types/Kifu';
  import KifuList from '$lib/components/KifuList.svelte';
  import { deleteKifu, searchKifus, updateKifuInfo } from '$lib/apis/kifu';
  import type { PaginationResponse } from '$lib/types/API';
  import { account } from '$lib/stores/session';

  // ----------------------------------------
  // 通知リスト
  interface Notification {
    id: string;
    type: 'like' | 'comment';
    kifuId: string;
    kifuTitle: string;
    userId: string;
    userName: string;
    createdAt: string;
    read: boolean;
  }

  let notifications: Notification[] = [];
  let unreadCount = 0;

  async function fetchNotificationList() {
    notifications = [
      {
        id: '1',
        type: 'like',
        kifuId: 'kifu-1',
        kifuTitle: '第1局：四間飛車',
        userId: 'user-2',
        userName: '許褚クリスティーナ',
        createdAt: '2024-01-10 15:30',
        read: false,
      },
      {
        id: '2',
        type: 'comment',
        kifuId: 'kifu-2',
        kifuTitle: '第2局：矢倉',
        userId: 'user-3',
        userName: '程昱ティノラベッラ',
        createdAt: '2024-01-09 18:45',
        read: true,
      },
    ];

    unreadCount = notifications.filter((n) => !n.read).length;
  }

  function markAsRead(notificationId: string) {
    notifications = notifications.map((n) => (n.id === notificationId ? { ...n, read: true } : n));
    unreadCount = notifications.filter((n) => !n.read).length;
  }

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

  async function fetchKifuList() {
    isLoadingKifuList = true;

    const result = await searchKifus('me', pagination.page, pagination.page_size, true);
    if (result.ok && result.data && result.pagination) {
      kifuList = result.data as KifuSummary[];
      pagination = result.pagination;
    } else {
      kifuList = []; // エラー時はリストをクリア
      console.error('Failed to fetch kifu list: ', result);
    }

    isLoadingKifuList = false;
  }

  async function handleChangePage(page: number) {
    pagination.page = page;
    await fetchKifuList();
  }

  async function handleTogglePublic(kifu: KifuSummary) {
    const result = await updateKifuInfo(
      kifu.id,
      kifu.title,
      !kifu.is_public,
      kifu.game_info,
      kifu.tags
    );
    if (result.ok) {
      await fetchKifuList();
    } else {
      console.error('Failed to toggle public: ', result);
    }
  }

  async function handleDeleteKifu(kifuId: string) {
    const result = await deleteKifu(kifuId);
    if (result.ok) {
      await fetchKifuList();
    } else {
      console.error('Failed to delete kifu: ', result);
    }
  }

  // ----------------------------------------
  // 初回データロード

  $: if (account) {
    fetchNotificationList();
    fetchKifuList();
  }
</script>

<div class="page">
  <!-- 通知セクション -->
  {#if notifications.length > 0}
    <section class="basic notification">
      <h2>
        通知 {#if unreadCount > 0}<span class="unread-count">{unreadCount}</span>{/if}
      </h2>
      <div class="notification-list">
        {#each notifications as notification}
          <div class="card notification-item" class:unread={!notification.read}>
            <div class="notification-content">
              <span class="user-name">{notification.userName}</span>
              さんが
              <a href={`/kifu/view?id=${notification.kifuId}`} class="kifu-link">
                {notification.kifuTitle}
              </a>
              を
              {#if notification.type === 'like'}
                いいね！しました
              {:else}
                コメントしました
              {/if}
              <span class="notification-date">{notification.createdAt}</span>
            </div>
            {#if !notification.read}
              <button class="mark-read-button" on:click={() => markAsRead(notification.id)}>
                既読にする
              </button>
            {/if}
          </div>
        {/each}
      </div>
    </section>

    <hr />
  {/if}

  <!-- 棋譜リストセクション -->
  <section class="basic kifu-list">
    <h2>自分の棋譜</h2>
    <KifuList
      {kifuList}
      {pagination}
      loading={isLoadingKifuList}
      mode="with-actions"
      changePage={handleChangePage}
      togglePublic={handleTogglePublic}
      deleteKifu={handleDeleteKifu}
    />
  </section>
</div>

<style lang="scss">
  section.notification {
    h2 {
      display: flex;
      align-items: center;
      gap: 0.5rem;

      .unread-count {
        display: inline-block;
        background-color: var(--primary-color);
        color: white;
        border-radius: 0.8rem;
        padding: 0.4rem;
        min-width: 1.6rem;
        font-size: 0.8rem;
        line-height: 0.8rem;
        text-align: center;
      }
    }

    .notification-list {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      .notification-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-left: 1.5rem;
        gap: 1rem;

        &.unread {
          background-color: var(--pickup-color);
          border-left: 0.5rem solid var(--primary-color);
          padding-left: 1rem;
        }

        .notification-content {
          flex: 1;

          .user-name {
            font-weight: bold;
            color: var(--primary-color);
          }

          .kifu-link {
            color: var(--secondary-color);
            text-decoration: underline;

            &:hover {
              opacity: 0.8;
            }
          }

          .notification-date {
            margin-left: 1rem;
            color: #666;
            font-size: 0.9rem;
          }
        }

        .mark-read-button {
          padding: 0.25rem 0.5rem;
          background-color: transparent;
          border: 1px solid var(--primary-color);
          color: var(--primary-color);
          border-radius: 4px;
          font-size: 0.9rem;

          &:hover {
            background-color: var(--primary-color);
            color: white;
          }
        }
      }
    }
  }
</style>
