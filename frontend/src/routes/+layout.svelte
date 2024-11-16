<!-- src/routes/+layout.svelte -->

<script lang="ts">
	import '$lib/styles/default.scss';
  import { getContext } from 'svelte';
  import type { Account } from '$lib/types/Account';

  // Account context
  const account = getContext<Account | null>('account');

  // $: isLoggedIn = !!account;
  $: isLoggedIn = true;
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
        <li><button class="logout">ログアウト</button></li>
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
    <p>©曹操オッキマラ 2024 kifup</p>
  </footer>
</div>

<style lang="scss">
  header {
    background-color: var(--primary-color);
    padding: 1rem;

    .logo {
      color: white;
      text-decoration: none;
      font-size: 1.5rem;
      font-weight: bold;
    }
  }

  nav {
    background-color: var(--secondary-color);
    padding: 0.5rem;

    ul {
      list-style: none;
      padding: 0;
      margin: 0;
      display: flex;
      gap: 1rem;
    }

    a {
      color: white;
      text-decoration: none;
      padding: 0.5rem 1rem;

      &:hover {
        background-color: rgba(255, 255, 255, 0.1);
        border-radius: 4px;
      }
    }

    .logout {
      background: none;
      border: 1px solid white;
      color: white;
      padding: 0.5rem 1rem;
      cursor: pointer;
      border-radius: 4px;

      &:hover {
        background-color: rgba(255, 255, 255, 0.1);
      }
    }
  }

  footer {
    background-color: var(--primary-color);
    color: white;
    padding: 1rem;
    text-align: center;
  }
</style>
