<!-- src/lib/components/AccountForm/RegisterForm.svelte -->

<script lang="ts">
  import { goto } from '$app/navigation';
  import { register, verifyEmail, verifyCode } from '$lib/apis/account';
  import { login } from '$lib/apis/session';
  import { sessionToken } from '$lib/stores/session';
  import { ArrowRightFromLine } from 'lucide-svelte';

  export let switchToLogin: () => void;

  // 新規登録の各ステップ
  type Step = 'email' | 'verify' | 'info';
  let currentStep: Step = 'email';
  let errorMessage = '';

  // フォームの入力値
  let email = '';
  let verificationCode = '';
  let name = '';
  let password = '';
  let passwordConfirm = '';

  // メールアドレス送信
  const handleEmailSubmit = async () => {
    const result = await verifyEmail(email);
    if (result.ok) {
      currentStep = 'verify';
      errorMessage = '';
    } else {
      errorMessage = 'メールアドレスの送信に失敗しました';
    }
  };

  // 認証コード検証
  const handleVerifySubmit = async () => {
    const result = await verifyCode(email, verificationCode);
    if (result.ok) {
      currentStep = 'info';
      errorMessage = '';
    } else {
      errorMessage = '認証コードの検証に失敗しました';
    }
  };

  // アカウント情報の登録
  const handleRegister = async () => {
    if (password !== passwordConfirm) {
      errorMessage = 'パスワードが一致していません';
      return;
    }

    const result = await register(name, email, password, verificationCode);
    if (!result.ok) {
      errorMessage = 'アカウントの作成に失敗しました';
      return;
    }

    // 登録成功したら自動ログイン
    const loginResult = await login(email, password);
    if (loginResult.ok) {
      sessionToken.set(loginResult.data);
      goto('/home');
    } else {
      errorMessage = 'ログインに失敗しました';
    }
  };
</script>

<section class="basic form-section">
  <form
    onsubmit={(e) => {
      e.preventDefault();
    }}
    class="basic register-form"
  >
    <h2>新規登録</h2>

    {#if currentStep === 'email'}
      <!-- メールアドレス入力 -->
      <div class="form-group">
        <label for="register-email">メールアドレス</label>
        <input type="email" id="register-email" bind:value={email} autocomplete="email" required />
      </div>
      <button type="submit" class="submit" onclick={handleEmailSubmit}>認証コードを送信</button>
    {:else if currentStep === 'verify'}
      <!-- 認証コード入力 -->
      <div class="info-text">
        {email}宛に認証コードを送信しました。 メールに記載された6桁のコードを入力してください。
      </div>
      <div class="form-group">
        <label for="verification-code">認証コード</label>
        <input
          type="text"
          id="verification-code"
          bind:value={verificationCode}
          maxlength="6"
          pattern="[0-9]{6}"
          required
        />
      </div>
      <button type="submit" class="submit" onclick={handleVerifySubmit}>認証コードを確認</button>
    {:else}
      <!-- アカウント情報入力 -->
      <div class="form-group">
        <label for="register-name">名前</label>
        <input type="text" id="register-name" bind:value={name} autocomplete="nickname" required />
      </div>
      <div class="form-group">
        <label for="register-password">パスワード</label>
        <input
          type="password"
          id="register-password"
          bind:value={password}
          autocomplete="new-password"
          required
        />
      </div>
      <div class="form-group">
        <label for="register-password-confirm">パスワード（確認）</label>
        <input
          type="password"
          id="register-password-confirm"
          bind:value={passwordConfirm}
          autocomplete="off"
          required
        />
      </div>
      <button type="submit" class="submit" onclick={handleRegister}>登録</button>
    {/if}

    {#if errorMessage}
      <div class="error-message">
        {errorMessage}
      </div>
    {/if}
  </form>

  <ul class="view-control">
    <li>
      <button onclick={switchToLogin}
        ><ArrowRightFromLine size={20} color="var(--primary-color)" /><span>ログイン</span></button
      >
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
