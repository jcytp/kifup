<!-- src/routes/+page.svelte -->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { register } from '$lib/apis/account';
  import { login } from '$lib/apis/session';
  import { account, sessionToken } from '$lib/stores/session';

  // ログイン済みの場合はホームへ遷移する
  $: if ($account) {
    goto('/home');
  }

  // フォーム表示の状態管理
  let activeSection: 'login' | 'register' | 'reset' = 'login';
  let errorMessage = '';

  // 各フォームの入力値
  let loginForm = {
    email: '',
    password: '',
  };
  let registerForm = {
    name: '',
    email: '',
    password: '',
    passwordConfirm: '',
  };
  let resetForm = {
    email: '',
  };

  // フォームの切り替え
  const switchToLogin = () => {
    activeSection = 'login';
    errorMessage = '';
  };
  const switchToRegister = () => {
    activeSection = 'register';
    errorMessage = '';
  };
  const switchToReset = () => {
    activeSection = 'reset';
    errorMessage = '';
  };

  const handleLogin = async () => {
    const result = await login(loginForm.email, loginForm.password);
    if (result.ok) {
      console.log('login success');
      sessionToken.set(result.data);
    } else {
      errorMessage = 'ログイン処理中にエラーが発生しました';
      return;
    }
    goto('/home');
  };

  const handleRegister = async () => {
    console.debug('handleRegister');
    if (registerForm.password !== registerForm.passwordConfirm) {
      errorMessage = 'パスワードが一致していません';
      return;
    }
    const result = await register(registerForm.name, registerForm.email, registerForm.password);
    if (!result.ok) {
      errorMessage = '新規登録処理中にエラーが発生しました';
      return;
    }
    const result_login = await login(registerForm.email, registerForm.password);
    if (result_login.ok) {
      console.log('login success');
      sessionToken.set(result_login.data);
    } else {
      errorMessage = 'ログイン処理中にエラーが発生しました';
      return;
    }
    goto('/home');
  };

  function handleReset() {
    // ToDo: パスワードリセット処理
    console.log('Reset:', resetForm);
  }
</script>

<div class="page">
  {#if activeSection === 'login'}
    <section class="basic form-section">
      <form on:submit|preventDefault={handleLogin} class="basic">
        <h2>ログイン</h2>
        <div class="form-group">
          <label for="login-email">メールアドレス</label>
          <input
            type="email"
            id="login-email"
            bind:value={loginForm.email}
            autocomplete="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="login-password">パスワード</label>
          <input
            type="password"
            id="login-password"
            bind:value={loginForm.password}
            autocomplete="current-password"
            required
          />
        </div>
        <button type="submit" class="submit">ログイン</button>
        {#if errorMessage}
          <div class="error-message">
            {errorMessage}
          </div>
        {/if}
      </form>
      <ul class="view-control">
        <li><button on:click={switchToRegister}>新規登録</button></li>
        <li><button on:click={switchToReset}>パスワードを忘れた</button></li>
      </ul>
    </section>
  {/if}

  {#if activeSection === 'register'}
    <section class="basic form-section">
      <form on:submit|preventDefault={handleRegister} class="basic">
        <h2>新規登録</h2>
        <div class="form-group">
          <label for="register-name">名前</label>
          <input
            type="text"
            id="register-name"
            bind:value={registerForm.name}
            autocomplete="nickname"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-email">メールアドレス</label>
          <input
            type="email"
            id="register-email"
            bind:value={registerForm.email}
            autocomplete="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-password">パスワード</label>
          <input
            type="password"
            id="register-password"
            bind:value={registerForm.password}
            autocomplete="new-password"
            required
          />
        </div>
        <div class="form-group">
          <label for="register-password-confirm">パスワード（確認）</label>
          <input
            type="password"
            id="register-password-confirm"
            bind:value={registerForm.passwordConfirm}
            autocomplete="off"
            required
          />
        </div>
        <button type="submit" class="submit">登録</button>
        {#if errorMessage}
          <div class="error-message">
            {errorMessage}
          </div>
        {/if}
      </form>
      <ul class="view-control">
        <li><button on:click={switchToLogin}>ログイン</button></li>
      </ul>
    </section>
  {/if}

  {#if activeSection === 'reset'}
    <section class="basic form-section">
      <form on:submit|preventDefault={handleReset} class="basic">
        <h2>パスワードの再設定</h2>
        <div class="form-group">
          <label for="reset-email">メールアドレス</label>
          <input
            type="email"
            id="reset-email"
            bind:value={resetForm.email}
            autocomplete="email"
            required
          />
        </div>
        <button type="submit" class="submit">送信</button>
        {#if errorMessage}
          <div class="error-message">
            {errorMessage}
          </div>
        {/if}
      </form>
      <ul class="view-control">
        <li><button on:click={switchToLogin}>ログイン</button></li>
        <li><button on:click={switchToRegister}>新規登録</button></li>
      </ul>
    </section>
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

<style>
  .page {
    .form-section {
      display: flex;
      padding: 1rem;

      form {
        margin-top: 0;
        width: 70%;

        h2 {
          color: var(--secondary-color);
          font-size: 1.4rem;
          line-height: 2rem;
          margin-bottom: 1rem;
        }

        .error-message {
          color: #dc2626;
          padding: 0.5rem;
          font-size: small;
          text-align: center;
        }
      }

      ul.view-control {
        display: flex;
        flex-direction: column;
        justify-content: end;
        gap: 0.5rem;
        width: 30%;
        padding: 0.5rem 1rem;

        button {
          background: var(--background-color);
          padding: 0.3rem 1rem;
          border: 1px solid var(--secondary-color);
          border-radius: 0.5rem;
          width: 100%;

          &:hover {
            background-color: var(--secondary-color);
            color: var(--background-color);
          }
        }
      }
    }

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
