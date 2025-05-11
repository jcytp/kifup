<!-- src/lib/components/AccountForm/ResetForm.svelte -->

<script lang="ts">
  import { resetPassword, verifyCode, verifyEmail } from '$lib/apis/account';
  import { ArrowRightFromLine } from 'lucide-svelte';

  export let switchToLogin: () => void;
  export let switchToRegister: () => void;

  let step: 1 | 2 | 3 | 4 = 1;
  let resetForm = {
    email: '',
    code: '',
    password: '',
    confirmPassword: '',
  };
  let errorMessage = '';
  let loading = false;

  const handleSendEmail = async () => {
    if (!resetForm.email) {
      errorMessage = 'メールアドレスを入力してください';
      return;
    }
    loading = true;
    errorMessage = '';
    const result = await verifyEmail(resetForm.email);
    if (result.ok) {
      step = 2;
    } else {
      errorMessage = '処理に失敗しました';
    }
    loading = false;
  };

  const handleVerifyCode = async () => {
    if (!resetForm.code) {
      errorMessage = '認証コードを入力してください';
      return;
    }
    loading = true;
    errorMessage = '';
    const result = await verifyCode(resetForm.email, resetForm.code);
    if (result.ok) {
      step = 3;
    } else {
      errorMessage = '処理に失敗しました';
    }
    loading = false;
  };

  const handleResetPassword = async () => {
    if (!resetForm.password || !resetForm.confirmPassword) {
      errorMessage = 'パスワードを入力してください';
      return;
    }
    if (resetForm.password != resetForm.confirmPassword) {
      errorMessage = 'パスワード確認が一致しません';
      return;
    }
    loading = true;
    errorMessage = '';
    const result = await resetPassword(resetForm.email, resetForm.password, resetForm.code);
    if (result.ok) {
      step = 4;
    } else {
      errorMessage = '処理に失敗しました';
    }
    loading = false;
  };

  const handleSubmit = () => {
    console.debug('Reset:', resetForm);
    switch (step) {
      case 1:
        handleSendEmail();
        break;
      case 2:
        handleVerifyCode();
        break;
      case 3:
        handleResetPassword();
        break;
    }
  };

  const goBack = () => {
    errorMessage = '';
    if (step > 1) {
      step--;
    }
  };
</script>

<section class="basic form-section">
  <form onsubmit={handleSubmit} class="basic register-form">
    <h2>パスワードの再設定</h2>

    {#if step === 1}
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
      <button type="submit" class="submit">認証コードを送信</button>
    {:else if step === 2}
      <div class="form-group">
        <label for="reset-code">認証コード</label>
        <input
          type="text"
          id="reset-code"
          bind:value={resetForm.code}
          placeholder="6桁のコードを入力"
          maxlength="6"
          required
        />
      </div>
      <button type="submit" class="submit">コードを検証</button>
      <button type="button" class="cancel" onclick={goBack}>戻る</button>
    {:else if step === 3}
      <div class="form-group">
        <label for="reset-password">新しいパスワード</label>
        <input
          type="password"
          id="reset-password"
          bind:value={resetForm.password}
          autocomplete="new-password"
          required
        />
      </div>
      <div class="form-group">
        <label for="reset-password-confirm">パスワード確認</label>
        <input
          type="password"
          id="reset-password-confirm"
          bind:value={resetForm.confirmPassword}
          autocomplete="new-password"
          required
        />
      </div>
      <button type="submit" class="submit">パスワードをリセット</button>
      <button type="button" class="cancel" onclick={goBack}>戻る</button>
    {:else if step === 4}
      <p>パスワードをリセットしました。</p>
      <p>ログインしてください。</p>
    {/if}
    {#if errorMessage}
      <div class="error-message">
        {errorMessage}
      </div>
    {/if}
  </form>

  <ul class="view-control">
    <li>
      <button onclick={switchToLogin}>
        <ArrowRightFromLine size={20} color="var(--primary-color)" />
        <span>ログイン</span>
      </button>
    </li>
    <li>
      <button onclick={switchToRegister}>
        <ArrowRightFromLine size={20} color="var(--primary-color)" />
        <span>新規登録</span>
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
