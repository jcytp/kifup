<!-- src/lib/components/AccountForm/LoginForm.svelte -->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { sessionToken } from '$lib/stores/session';
  import { login } from '$lib/apis/session';
  import { ArrowRightFromLine } from 'lucide-svelte';

  export let switchToRegister: () => void;
  export let switchToReset: () => void;

  let loginForm = {
    email: '',
    password: '',
  };
  let errorMessage = '';

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
</script>

<section class="basic form-section">
  <form onsubmit={handleLogin} class="basic register-form">
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
    <li>
      <button onclick={switchToRegister}>
        <ArrowRightFromLine size={20} color="var(--primary-color)" />
        <span>新規登録</span>
      </button>
    </li>
    <li>
      <button onclick={switchToReset}>
        <ArrowRightFromLine size={20} color="var(--primary-color)" />
        <span style="font-size: x-small;">パスワード再設定</span>
      </button>
    </li>
  </ul>
</section>

<style lang="scss">
  @import '../../styles/mixins.scss';

  section.basic.form-section {
    display: flex;
    padding: 1rem;
    @include sp {
      flex-direction: column;
      align-items: end;
    }
    @include pc {
      flex-direction: row;
    }

    form.basic.register-form {
      margin-top: 0;
      @include sp {
        width: 100%;
      }
      @include pc {
        width: 70%;
      }

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
      @include sp {
        width: 10rem;
        padding: 0.5rem 0;
      }
      @include pc {
        width: 30%;
        padding: 0.5rem 1rem;
      }

      button {
        display: flex;
        align-items: center;
        justify-content: space-between;
        background: var(--background-color);
        padding: 0.3rem 1rem;
        border-bottom: 1px solid var(--secondary-color);
        width: 100%;

        &:hover {
          border-radius: 0.5rem;
          background-color: var(--secondary-color);
          color: var(--background-color);
        }
      }
    }
  }
</style>
