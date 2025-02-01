<!-- src/routes/+page.svelte -->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { account } from '$lib/stores/session';
  import RegisterForm from '$lib/components/AccountForm/RegisterForm.svelte';
  import LoginForm from '$lib/components/AccountForm/LoginForm.svelte';
  import ResetForm from '$lib/components/AccountForm/ResetForm.svelte';

  // ログイン済みの場合はホームへ遷移する
  $: if ($account) {
    goto('/home');
  }

  // フォーム表示の状態管理
  let activeSection: 'login' | 'register' | 'reset' = 'login';

  // フォームの切り替え
  const switchToLogin = () => {
    activeSection = 'login';
  };
  const switchToRegister = () => {
    activeSection = 'register';
  };
  const switchToReset = () => {
    activeSection = 'reset';
  };
</script>

<div class="page">
  {#if activeSection === 'login'}
    <LoginForm {switchToRegister} {switchToReset} />
  {:else if activeSection === 'register'}
    <RegisterForm {switchToLogin} />
  {:else if activeSection === 'reset'}
    <ResetForm {switchToLogin} {switchToRegister} />
  {/if}

  <hr />

  <section class="basic feature-section">
    <h2>kifupとは？</h2>
    <p>
      kifupは、将棋の棋譜を管理・共有するためのWebアプリケーションです。
      対局の記録を保存し、振り返りや研究に活用できます。
    </p>
    <div class="features">
      <div class="card">
        <h3>棋譜の管理</h3>
        <p>棋譜を簡単にアップロードし、整理できます。</p>
      </div>
      <div class="card">
        <h3>対局の再生</h3>
        <p>保存した棋譜を再生して、局面を確認できます。</p>
      </div>
      <div class="card">
        <h3>棋譜の共有</h3>
        <p>アップロードした棋譜を、共有リンクで公開できます。</p>
      </div>
    </div>
  </section>
</div>

<style lang="scss">
  .page {
    .feature-section {
      .features {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(16rem, 1fr));
        gap: 1rem;
        margin-top: 1rem;
      }
    }
  }
</style>
