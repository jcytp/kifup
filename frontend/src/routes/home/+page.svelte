<!-- src/routes/home/+page.svelte -->

<script lang="ts">
  import { onMount } from 'svelte';
  import type { Kifu } from '$lib/types/Kifu';

  // モック通知データ
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

  // モック棋譜データ
  let kifuList: Kifu[] = [];
  let isLoading = false;

  // 選択中の棋譜ID（ポップアップ用）
  let selectedKifuId: string | null = null;

  // ページネーションの状態管理
  let currentPage = 1;
  const itemsPerPage = 10;
  let totalPages = 1;

  // モックデータを取得
  async function fetchData() {
    isLoading = true;
    try {
      // 通知データ
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

      // 棋譜リスト
      kifuList = Array(15)
        .fill(null)
        .map((_, i) => ({
          id: `kifu-${i + 1}`,
          ownerId: 'current-user',
          title: `第${i + 1}局：${['四間飛車', '矢倉', '角換わり', '横歩取り', '相掛かり'][i % 5]}`,
          matchInfo: {
            black: '曹操オッキマラ',
            white: ['許褚クリスティーナ', '程昱ティノラベッラ', '夏侯惇', '夏侯淵', '張遼'][i % 5],
            date: '2024-01-01',
          },
          tags: ['実戦', ['四間飛車', '矢倉', '角換わり', '横歩取り', '相掛かり'][i % 5]],
          isPublic: i % 2 === 0,
          moves: [],
        }));

      totalPages = Math.ceil(kifuList.length / itemsPerPage);
    } finally {
      isLoading = false;
    }
  }

  // ページネーションで表示する棋譜リスト
  $: paginatedKifuList = kifuList.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // ページ変更
  function changePage(page: number) {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
      selectedKifuId = null; // ページ変更時にポップアップを閉じる
    }
  }

  // 通知を既読にする
  function markAsRead(notificationId: string) {
    notifications = notifications.map((n) => (n.id === notificationId ? { ...n, read: true } : n));
    unreadCount = notifications.filter((n) => !n.read).length;
  }

  // 棋譜の公開状態を切り替え
  function togglePublic(kifuId: string) {
    kifuList = kifuList.map((k) => (k.id === kifuId ? { ...k, isPublic: !k.isPublic } : k));
    selectedKifuId = null;
  }

  // 棋譜を削除
  function deleteKifu(kifuId: string) {
    kifuList = kifuList.filter((k) => k.id !== kifuId);
    selectedKifuId = null;
  }

  // イベントハンドラ（MouseEventとKeyboardEventの両方に対応）
  function handleCardClick(event: MouseEvent | KeyboardEvent, kifuId: string) {
    // ポップアップメニュー内のクリックはカード自体のクリックイベントを発火させない
    if ((event.target as HTMLElement).closest('.popup-menu')) {
      return;
    }
    selectedKifuId = selectedKifuId === kifuId ? null : kifuId;
  }

  onMount(fetchData);
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
    {#if isLoading}
      <div class="loading">読み込み中...</div>
    {:else if currentPage == 1 && paginatedKifuList.length === 0}
      <div class="loading">棋譜はまだありません。</div>
    {:else}
      <div class="kifu-list-container">
        {#each paginatedKifuList as kifu}
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
                <span class={`kifu-status ${kifu.isPublic ? 'public' : 'private'}`}>
                  {kifu.isPublic ? '公開' : '非公開'}
                </span>
              </div>

              <div class="kifu-info">
                <span>対局日: {kifu.matchInfo.date}</span>
                <span>先手: {kifu.matchInfo.black}</span>
                <span>後手: {kifu.matchInfo.white}</span>
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
                <a href={`/kifu/edit?id=${kifu.id}`} class="menu-item" role="menuitem">
                  編集する
                </a>
                <button
                  type="button"
                  class="menu-item"
                  role="menuitem"
                  on:click={() => togglePublic(kifu.id)}
                >
                  {kifu.isPublic ? '非公開にする' : '公開する'}
                </button>
                <button
                  type="button"
                  class="menu-item delete"
                  role="menuitem"
                  on:click={() => deleteKifu(kifu.id)}
                >
                  削除する
                </button>
              </div>
            {/if}
          </div>
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

  section.kifu-list {
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
    }
  }
</style>
