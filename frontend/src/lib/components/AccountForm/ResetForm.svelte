<!-- src/lib/components/AccountForm/ResetForm.svelte -->

<script lang="ts">
  import { ArrowRightFromLine } from 'lucide-svelte';

  export let switchToLogin: () => void;
  export let switchToRegister: () => void;

  let resetForm = {
    email: '',
  };
  let errorMessage = '';

  function handleReset() {
    // ToDo: パスワードリセット処理
    console.log('Reset:', resetForm);
  }
</script>

<section class="basic form-section">
  <form onsubmit={handleReset} class="basic register-form">
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
