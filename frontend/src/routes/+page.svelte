<!-- src/routes/+page.svelte -->

<script lang="ts">
  let loginForm = {
    email: '',
    password: ''
  };

  let registerForm = {
    name: '',
    email: '',
    password: '',
    passwordConfirm: ''
  };

  function handleLogin() {
    // ToDo: ログイン処理
    console.log('Login:', loginForm);
  }

  function handleRegister() {
    // ToDo: 新規登録処理
    console.log('Register:', registerForm);
  }
</script>

<div class="top">
  <section class="login">
    <h2>ログイン</h2>
    <form on:submit|preventDefault={handleLogin}>
      <div class="form-group">
        <label for="login-email">メールアドレス</label>
        <input
          type="email"
          id="login-email"
          bind:value={loginForm.email}
          required
        />
      </div>
      <div class="form-group">
        <label for="login-password">パスワード</label>
        <input
          type="password"
          id="login-password"
          bind:value={loginForm.password}
          required
        />
      </div>
      <button type="submit">ログイン</button>
    </form>
  </section>

  <section class="register">
    <h2>新規登録</h2>
    <form on:submit|preventDefault={handleRegister}>
      <div class="form-group">
        <label for="register-name">名前</label>
        <input
          type="text"
          id="register-name"
          bind:value={registerForm.name}
          required
        />
      </div>
      <div class="form-group">
        <label for="register-email">メールアドレス</label>
        <input
          type="email"
          id="register-email"
          bind:value={registerForm.email}
          required
        />
      </div>
      <div class="form-group">
        <label for="register-password">パスワード</label>
        <input
          type="password"
          id="register-password"
          bind:value={registerForm.password}
          required
        />
      </div>
      <div class="form-group">
        <label for="register-password-confirm">パスワード（確認）</label>
        <input
          type="password"
          id="register-password-confirm"
          bind:value={registerForm.passwordConfirm}
          required
        />
      </div>
      <button type="submit">登録</button>
    </form>
  </section>

  <section class="about">
    <h2>kifupとは？</h2>
    <p>
      kifupは、将棋の棋譜を管理・共有するためのWebアプリケーションです。
      対局の記録を保存し、振り返りや研究に活用できます。
    </p>
    <div class="features">
      <div class="feature">
        <h3>棋譜の管理</h3>
        <p>棋譜を簡単にアップロードし、整理できます。</p>
      </div>
      <div class="feature">
        <h3>対局の再生</h3>
        <p>保存した棋譜を再生して、局面を確認できます。</p>
      </div>
      <div class="feature">
        <h3>コミュニティ</h3>
        <p>他のユーザーと棋譜を共有し、コメントを交換できます。</p>
      </div>
    </div>
  </section>
</div>

<style>
  .top {
    display: grid;
    gap: 2rem;
  }

  section {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .form-group {
    margin-bottom: 1rem;
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    color: var(--text-color);
  }

  input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    box-sizing: border-box;
  }

  button {
    background-color: var(--primary-color);
    color: white;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    width: 100%;
  }

  button:hover {
    background-color: var(--secondary-color);
  }

  .features {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
  }

  .feature {
    padding: 1rem;
    background-color: var(--background-color);
    border-radius: 4px;
  }

  @media (min-width: 768px) {
    .top {
      grid-template-columns: repeat(2, 1fr);
    }

    .about {
      grid-column: 1 / -1;
    }
  }
</style>
