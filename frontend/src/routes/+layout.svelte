<!-- src/routes/+layout.svelte -->

<script lang="ts">
  import '$lib/styles/default.scss';
  import { account, sessionToken } from '$lib/stores/session';
  import { onMount } from 'svelte';
  import { navigating } from '$app/stores';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Menu } from 'lucide-svelte';

  const isCurrentPath = (path: string) => {
    const currentPath = $page.url.pathname;
    return path.endsWith('*') ? currentPath.startsWith(path.slice(0, -1)) : currentPath === path;
  };

  // ----------------------------------------
  // セッション初期化
  onMount(async () => {
    await sessionToken.initialize(logoutMove);
  });
  $: isLoggedIn = !!$account;

  // ログアウト時の画面遷移
  const logoutMove = () => {
    // pageのurlがvisitorUrlsのいずれかに一致しない場合、"/"に遷移する
    const visitorUrls = ['/', '/kifu/search/', '/kifu/view/', '/account/*'];
    const currentPath = $page.url.pathname;
    const isVisitorUrl = visitorUrls.some((pattern) => {
      if (pattern.endsWith('*')) {
        const basePattern = pattern.slice(0, -1);
        return currentPath.startsWith(basePattern);
      }
      return currentPath === pattern;
    });
    if (!isVisitorUrl) {
      console.debug('unauthorized access');
      goto('/');
    }
  };

  // ログアウト
  const logout = () => {
    sessionToken.set(null);
  };

  // ----------------------------------------
  // メニューの展開
  let isMenuOpen = false;
  const toggleMenu = () => {
    isMenuOpen = !isMenuOpen;
  };

  $: if ($navigating) {
    // ページ遷移が発生したらメニューを閉じる
    isMenuOpen = false;
  }
</script>

<div class="app">
  <header>
    <button class="menu-button" onclick={toggleMenu}>
      <Menu size={27} color="var(--primary-color)" />
    </button>
    <h1><a href="/" class="logo">棋譜UP</a></h1>
  </header>
  {#if isMenuOpen}
    <div
      class="overlay"
      role="button"
      tabindex="0"
      onclick={toggleMenu}
      onkeydown={(e) => e.key === 'Escape' && toggleMenu()}
    ></div>
  {/if}
  <nav class:open={isMenuOpen}>
    <ul>
      {#if isLoggedIn}
        <li><a href="/home/" class:active={isCurrentPath('/home/')}>ホーム</a></li>
        <li><a href="/kifu/new/" class:active={isCurrentPath('/kifu/new/')}>新規作成</a></li>
        <li><a href="/kifu/search/" class:active={isCurrentPath('/kifu/search/')}>棋譜検索</a></li>
        <li><a href="/settings/" class:active={isCurrentPath('/settings/')}>設定</a></li>
        <li class="newgroup"><button onclick={logout}>ログアウト</button></li>
      {:else}
        <li><a href="/" class:active={isCurrentPath('/')}>ログイン</a></li>
        <li><a href="/kifu/search/" class:active={isCurrentPath('/kifu/search/')}>棋譜検索</a></li>
      {/if}
    </ul>
  </nav>
  <main>
    <slot />
  </main>
  <footer>
    <p>©jcytp 2024 kifup</p>
  </footer>
</div>

<style lang="scss">
  @import '../lib/styles/mixins.scss';

  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;

    header {
      @include sp {
        padding: 0.5rem 0 0 3.5rem;
      }
      @include pc {
        padding: 0.5rem 0 0 11rem;
      }
      background-color: var(--background-color);
      border-bottom: 0.3rem var(--primary-color) solid;
      padding: 0.5rem 0 0 11rem;
      z-index: var(--z-index-header);

      .menu-button {
        @include sp {
          position: absolute;
          top: 1rem;
          left: 1rem;
        }
        @include pc {
          display: none;
        }
      }

      .logo {
        color: var(--primary-color);
        text-decoration: none;
        font-size: 1.7rem;
        font-weight: bold;
        line-height: 2.7rem;
      }
    }

    .overlay {
      @include sp {
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        width: 100%;
        height: 100%;
        z-index: var(--z-index-nav);
        background: rgba(0, 0, 0, 0.2);
      }
      @include pc {
        display: none;
      }
    }

    nav {
      @include sp {
        left: -100%; // 初期状態は画面外
        transition: left 0.3s ease;
        box-shadow: 2px 0 8px rgba(0, 0, 0, 0.2);

        &.open {
          left: 0; // メニュー表示
        }
      }
      @include pc {
        left: 0;
      }
      position: fixed;
      top: 0;
      bottom: 0;
      z-index: var(--z-index-nav);
      background-color: var(--background-color);
      padding: 6rem 0.5rem 4rem 0.5rem;
      overflow-y: auto;

      ul {
        display: flex;
        flex-direction: column;
        list-style: none;
        gap: 0.5rem;

        li a,
        li button {
          display: block;
          width: 9rem;
          border: 2px solid transparent;
          border-radius: 0.5rem;
          padding-left: 1rem;
          line-height: 2rem;
          color: var(--primary-color);
          text-decoration: none;
          text-align: left;

          &:hover:not(.active) {
            background-color: var(--secondary-color);
            color: var(--background-color);
          }
        }

        li a.active {
          border-color: var(--border-color);
          cursor: default;
        }

        li.newgroup {
          margin-top: 0.5rem;
          border-top: var(--border-color) 1px solid;
          padding-top: 1rem;
        }
      }
    }

    main {
      @include sp {
        padding: 1.5rem 0 2rem 0;
      }
      @include pc {
        padding: 1.5rem 0 2rem 10rem;
      }
      flex: 1;
      max-width: 72rem;
      width: 100%;
      box-sizing: border-box;
    }

    footer {
      border-top: var(--primary-color) 0.2rem solid;
      background-color: var(--background-color);
      color: var(--primary-color);
      padding: 0.2rem 0 0.8rem 11rem;
      z-index: var(--z-index-header);

      p {
        font-size: 1rem;
        line-height: 2rem;
      }
    }
  }
</style>
