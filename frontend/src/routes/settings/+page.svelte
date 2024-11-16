<!-- src/routes/settings/+page.svelte -->

<script lang="ts">
	import { page } from '$app/stores';
  import type { Account } from '$lib/types/Account';
  import { onMount } from 'svelte';

  // アカウント情報の状態管理
  let account: Account | null = null;
  let isLoading = true;
  let error: string | null = null;

  // 設定フォームの状態
  let profileForm = {
    name: '',
    bio: '',
    profileImage: null as File | null,
    notifications: {
      likes: true,
      comments: true
    }
  };

  // 現在のプロフィール画像のプレビューURL
  let imagePreview: string | null = null;

  // アカウント情報の取得
  async function fetchAccountData() {
    isLoading = true;
    try {
      // TODO: API実装後に実際のデータ取得に置き換え
      await new Promise(resolve => setTimeout(resolve, 500));
      
      account = {
        id: 'user-1',
        name: 'サンプルユーザー',
        email: 'sample@example.com'
      };

      // フォームの初期値を設定
      profileForm.name = account.name;
      profileForm.bio = ''; // TODO: Account型に追加後に実装
      profileForm.notifications = {
        likes: true,
        comments: true
      };
      
    } catch (e) {
      error = 'アカウント情報の取得に失敗しました。';
    } finally {
      isLoading = false;
    }
  }

  // プロフィール画像の処理
  function handleImageChange(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    
    if (file) {
      profileForm.profileImage = file;
      // プレビュー用URLを生成
      imagePreview = URL.createObjectURL(file);
    }
  }

  // プロフィール画像のプレビューのクリーンアップ
  function cleanupImagePreview() {
    if (imagePreview) {
      URL.revokeObjectURL(imagePreview);
    }
  }

  // 設定の保存
  async function handleSubmit() {
    try {
      // TODO: API実装後に実際の保存処理に置き換え
      console.log('Saving settings:', profileForm);
      // 成功メッセージの表示など
    } catch (e) {
      error = '設定の保存に失敗しました。';
    }
  }

  onMount(() => {
    fetchAccountData();
  });

  // コンポーネントのアンマウント時にプレビューをクリーンアップ
  onMount(() => {
    return () => cleanupImagePreview();
  });
</script>

<div class="container">
  {#if isLoading}
    <div class="loading">
      <p>設定を読み込んでいます...</p>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
    </div>
  {:else if account}
    <h1>設定</h1>

    <form on:submit|preventDefault={handleSubmit} class="settings-form">
      <section class="settings-section">
        <h2>通知設定</h2>
        <div class="notification-options">
          <label class="checkbox-label">
            <input
              type="checkbox"
              bind:checked={profileForm.notifications.likes}
            />
            いいねを通知する
          </label>
          <label class="checkbox-label">
            <input
              type="checkbox"
              bind:checked={profileForm.notifications.comments}
            />
            コメントを通知する
          </label>
        </div>
      </section>

      <section class="settings-section">
        <h2>プロフィール設定</h2>
        
        <div class="form-group">
          <h3>プロフィール画像</h3>
          <div class="profile-image-container">
            <div class="profile-image">
              {#if imagePreview}
                <img src={imagePreview} alt="プロフィール画像プレビュー" />
              {:else}
                <div class="placeholder-image">
                  {profileForm.name[0]}
                </div>
              {/if}
            </div>
            <div class="image-upload">
              <label for="profile-image" class="upload-button">
                画像を選択
              </label>
              <input
                type="file"
                id="profile-image"
                accept="image/*"
                on:change={handleImageChange}
                class="hidden"
              />
              <p class="upload-note">
                推奨: 500x500px以上の正方形の画像
              </p>
            </div>
          </div>
        </div>

        <div class="form-group">
          <label for="name">名前</label>
          <input
            type="text"
            id="name"
            bind:value={profileForm.name}
            required
          />
        </div>

        <div class="form-group">
          <label for="bio">自己紹介</label>
          <textarea
            id="bio"
            bind:value={profileForm.bio}
            rows="5"
            placeholder="自己紹介文を入力してください"
          ></textarea>
        </div>
      </section>

      <div class="form-actions">
        <button type="submit" class="save-button">
          設定を保存
        </button>
      </div>
    </form>
  {/if}
</div>

<style lang="scss">
  .container {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;

    h1 {
      color: var(--primary-color);
      margin-bottom: 2rem;
    }
  }

  .loading, .error {
    text-align: center;
    padding: 2rem;
  }

  .settings-form {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .settings-section {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    h2 {
      color: var(--primary-color);
      margin-bottom: 1.5rem;
      font-size: 1.2rem;
    }

    h3 {
      color: var(--text-color);
      margin-bottom: 1rem;
      font-size: 1rem;
    }
  }

  .notification-options {
    display: flex;
    flex-direction: column;
    gap: 1rem;

    .checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      cursor: pointer;

      input[type="checkbox"] {
        width: 1.2rem;
        height: 1.2rem;
        cursor: pointer;
      }
    }
  }

  .form-group {
    margin-bottom: 1.5rem;

    label {
      display: block;
      margin-bottom: 0.5rem;
      color: var(--text-color);
    }

    input[type="text"],
    textarea {
      width: 100%;
      padding: 0.75rem;
      border: 1px solid var(--border-color);
      border-radius: 4px;
      cursor: text;

      &:focus {
        border-color: var(--secondary-color);
      }
    }
  }

  .profile-image-container {
    display: flex;
    gap: 2rem;
    align-items: flex-start;

    .profile-image {
      flex: 0 0 120px;
      height: 120px;
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
        font-size: 3rem;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .image-upload {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;

      .upload-button {
        display: inline-block;
        padding: 0.5rem 1rem;
        background-color: var(--secondary-color);
        color: white;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.9rem;

        &:hover {
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

  .form-actions {
    display: flex;
    justify-content: center;

    .save-button {
      padding: 0.75rem 2rem;
      background-color: var(--primary-color);
      color: white;
      border-radius: 4px;
      font-weight: bold;
      cursor: pointer;

      &:hover {
        opacity: 0.9;
      }
    }
  }
</style>
