<!-- src/routes/settings/+page.svelte -->

<script lang="ts">
  import { deleteAccount, updateAccountInfo } from '$lib/apis/account';
  import { account, sessionToken } from '$lib/stores/session';
  import type { Account } from '$lib/types/Account';
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';

  const account_url_base = '/account/';
  const image_url_base = 'http://example.com/icon/';
  let imagePreview: string | null = null;
  let viewMode: 'setting' | 'password' | 'delete' = 'setting';

  // プロフィール画像の処理
  function handleImageChange(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];

    if (file) {
      // プレビュー用URLを生成
      imagePreview = URL.createObjectURL(file);
    }
  }

  // プロフィール画像のプレビューのクリーンアップ
  const cleanupImagePreview = () => {
    if (imagePreview) {
      URL.revokeObjectURL(imagePreview);
    }
  };

  // 設定の保存
  const handleSubmit = async () => {
    if (!accountInfo) {
      return;
    }
    const result = await updateAccountInfo(
      accountInfo.name,
      accountInfo.icon_id,
      accountInfo.introduction
    );
    if (!result.ok) {
      console.error('Failed to update account info: ', result);
    }
  };

  // コンポーネントのアンマウント時にプレビューをクリーンアップ
  onMount(() => {
    return cleanupImagePreview;
  });

  let accountInfo: Account | null = null;
  $: if ($account) {
    accountInfo = get(account);
  }

  // ToDo: アカウント情報に追加が必要
  let likes_notification = true;
  let comments_notification = true;

  // ----------------------------------------
  // パスワードリセット
  const handleResetPassword = async () => {};

  // ----------------------------------------
  // アカウント削除

  const handleDeleteAccount = async () => {
    const result = await deleteAccount();
    if (result.ok) {
      sessionToken.set(null);
    } else {
      console.error('Failed to delete account: ', result);
    }
  };
</script>

<div class="page">
  {#if accountInfo}
    <section class="basic">
      <h2>設定</h2>

      {#if viewMode === 'setting'}
        <form onsubmit={handleSubmit} class="basic settings-form">
          <div class="form-group">
            <label for="name">アカウントページ</label>
            <p class="account-link">
              <a href={`${account_url_base}?id=${accountInfo.id}`}
                >{`${account_url_base}?id=${accountInfo.id}`}</a
              >
            </p>
          </div>

          <div class="form-group">
            <label for="name">名前</label>
            <input type="text" id="name" bind:value={accountInfo.name} required />
          </div>

          <div class="form-group">
            <h3 class="label">プロフィール画像</h3>
            <div class="profile-image-container">
              <div class="profile-image">
                {#if imagePreview}
                  <img src={imagePreview} alt="プロフィール画像プレビュー" />
                {:else if accountInfo.icon_id}
                  <img src={`${image_url_base}${accountInfo.icon_id}`} alt="プロフィール画像" />
                {:else}
                  <div class="placeholder-image">
                    {accountInfo.name[0]}
                  </div>
                {/if}
              </div>
              <!-- <div class="image-upload">
                <label for="profile-image" class="upload-button"> 画像を選択 </label>
                <input
                  type="file"
                  id="profile-image"
                  accept="image/*"
                  onchange={handleImageChange}
                  class="hidden"
                />
                <p class="upload-note">推奨: 200x200px以上の正方形の画像</p>
              </div> -->
            </div>
          </div>

          <div class="form-group">
            <label for="bio">自己紹介</label>
            <textarea
              id="bio"
              bind:value={accountInfo.introduction}
              rows="5"
              placeholder="自己紹介文を入力してください"
            ></textarea>
          </div>

          <!-- <div class="form-group">
            <h3 class="label">通知設定</h3>
            <label class="checkbox-label">
              <input type="checkbox" bind:checked={likes_notification} />
              いいねを通知する
            </label>
            <label class="checkbox-label">
              <input type="checkbox" bind:checked={comments_notification} />
              コメントを通知する
            </label>
          </div> -->

          <button type="submit" class="submit">設定を保存</button>
        </form>
      {/if}

      {#if viewMode === 'password'}
        <form onsubmit={handleResetPassword} class="basic settings-form">
          <!-- ToDo: imprement to reset password -->
        </form>
      {/if}

      {#if viewMode === 'delete'}
        <form onsubmit={handleDeleteAccount} class="basic settings-form">
          <p>アカウントを削除すると、作成した全ての棋譜が失われます。</p>
          <p>アカウントを削除しますか？</p>
          <div class="controls">
            <button onclick={() => (viewMode = 'setting')}>キャンセル</button>
            <button type="submit" class="submit warning">削除する</button>
          </div>
        </form>
      {/if}
    </section>

    {#if viewMode === 'setting'}
      <section class="basic">
        <div class="controls">
          <!-- <button onclick={() => (viewMode = 'password')}>パスワード変更</button> -->
          <button onclick={() => (viewMode = 'delete')} class="warning">アカウントを削除</button>
        </div>
      </section>
    {/if}
  {/if}
</div>

<style lang="scss">
  .account-link {
    padding: 0 0 0 0.5rem;

    a {
      color: #369;
      text-decoration: underline;

      &:hover,
      &:focus {
        opacity: 0.9;
      }
    }
  }
  .profile-image-container {
    display: flex;
    gap: 2rem;
    align-items: flex-start;

    .profile-image {
      flex: 0 0 6rem;
      height: 6rem;
      overflow: hidden;
      border-radius: 50%;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .placeholder-image {
        width: 100%;
        height: 100%;
        background-color: var(--secondary-color);
        color: white;
        font-size: 4rem;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .image-upload {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: end;
      gap: 0.2rem;
      height: 6rem;

      .upload-button {
        display: block;
        background-color: var(--secondary-color);
        color: var(--background-color);
        margin: 0;
        padding: 0.4rem 1rem;
        border: none;
        border-radius: 0.5rem;
        width: 10rem;
        font-size: 0.9rem;
        text-align: center;
        cursor: pointer;

        &:hover,
        &:focus {
          opacity: 0.9;
        }
      }

      .upload-note {
        font-size: 0.9rem;
        color: var(--text-color);
        opacity: 0.8;
      }

      .hidden {
        display: none;
      }
    }
  }

  .controls {
    text-align: center;

    button {
      width: 10rem;
      padding: 0.3rem 0.5rem;
      background-color: var(--primary-color);
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;

      &.warning {
        background-color: var(--warning-color);
      }

      &:hover {
        opacity: 0.9;
      }
    }
  }
</style>
