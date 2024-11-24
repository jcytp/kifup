<!-- src/routes/+layout.svelte -->

<script lang="ts">
	import '$lib/styles/default.scss';
  import { onMount } from 'svelte';
	import { account, token } from '$lib/stores/session';
	import { goto } from '$app/navigation';
	import { GetAccount } from '$lib/api/account';
	import type { Account } from '$lib/types/Account';

  // ロード時にLocalStorageのトークンを取得
  onMount(async () => {
    const currentToken = localStorage.getItem('token');
    if (currentToken) {
      token.set(currentToken);
    }
  })

  // トークンが更新された場合
  $: if ($token) {
    // アカウント情報を取得
    (async () => {
      const result = await GetAccount()
      if (result.ok) {
        const new_account = result.data as Account;
        account.set(new_account);
      } else {
        // アカウント情報を取得できなかった場合、トークンを削除
        token.clear();
      }
    })();
  } else {
    // トークンが削除された場合（ログアウト時）は、トップページへ遷移
    account.set(null);
    goto('/');
  }

  $: isLoggedIn = !!$account;

  const logout = () => {
    token.clear();
  }
</script>

<div class="app">
  <header>
    <h1><a href="/" class="logo">kifup</a></h1>
  </header>
  <nav>
    <ul>
      {#if isLoggedIn}
        <li><a href="/home">ホーム</a></li>
        <li><a href="/kifu/new">新規作成</a></li>
        <li><a href="/kifu/search">棋譜検索</a></li>
        <li><a href="/settings">設定</a></li>
        <li class="newgroup"><button on:click={logout}>ログアウト</button></li>
      {:else}
        <li><a href="/">ログイン</a></li>
        <li><a href="/kifu/search">棋譜検索</a></li>
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
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;

    header {
      background-color: var(--background-color);
      border-bottom: 0.3rem var(--primary-color) solid;
      padding: 0.5rem 0 0rem 11rem;
      z-index: var(--z-index-header);

      .logo {
        color: var(--primary-color);
        text-decoration: none;
        font-size: 1.7rem;
        font-weight: bold;
        line-height: 2.7rem;
      }
    }

    nav {
      position: fixed;
      top: 0;
      left: 0;
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
          padding-left: 1rem;
          line-height: 2rem;
          color: var(--primary-color);
          text-decoration: none;
          text-align: left;

          &:hover {
            background-color: var(--secondary-color);
            border-radius: 0.5rem;
            color: var(--background-color);
          }
        }

        li.newgroup {
          margin-top: 0.5rem;
          border-top: var(--border-color) 1px solid;
          padding-top: 1rem;
        }
      }
    }

    main {
      flex: 1;
      padding: 2rem 0 2rem 10rem;
      max-width: 1200px;
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
