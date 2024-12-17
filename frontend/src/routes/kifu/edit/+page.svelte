<script lang="ts">
  import { page } from '$app/stores';
  import type { KifuDetail, KifuMove } from '$lib/types/Kifu';
  import KifuPlayer from '$lib/components/KifuPlayer.svelte';
  import { getKifu, updateKifuInfo } from '$lib/apis/kifu';
  import { account } from '$lib/stores/session';

  // ----------------------------------------
  // 棋譜データの状態管理

  const kifuId = $page.url.searchParams.get('id'); // URLからkifuIDを取得

  let isError = false;
  let kifu: KifuDetail;

  const fetchKifuData = async () => {
    isError = false;

    if (!kifuId) {
      console.error('No kifuId');
      isError = true;
      return;
    }

    const result = await getKifu(kifuId);
    if (result.ok && result.data) {
      kifu = result.data as KifuDetail;
    } else {
      console.error('Failed to fetch kifu detail: ', result);
      isError = true;
      return;
    }

    formData = {
      title: kifu.title,
      isPublic: kifu.is_public,
      gameInfo: {
        先手: kifu.game_info['先手'],
        後手: kifu.game_info['後手'],
        対局日時: kifu.game_info['対局日時'],
        対局場所: kifu.game_info['対局場所'],
        持ち時間: kifu.game_info['持ち時間'],
        秒読み: kifu.game_info['秒読み'],
        秒加算: kifu.game_info['秒加算'],
      },
      tags: kifu.tags,
    };
    initialTimeInput = kifu.game_info['持ち時間']
      ? Math.floor(parseInt(kifu.game_info['持ち時間']) / 60)
      : 0;
    byoyomiTimeInput = kifu.game_info['秒読み']
      ? Math.floor(parseInt(kifu.game_info['秒読み']) / 60)
      : 0;
    incrementTimeInput = kifu.game_info['秒加算']
      ? Math.floor(parseInt(kifu.game_info['秒加算']) / 60)
      : 0;
    moves = kifu.moves;
  };

  // ----------------------------------------
  // 棋譜情報フォーム管理

  interface KifuFormData {
    title: string;
    isPublic: boolean;
    gameInfo: {
      先手: string;
      後手: string;
      対局日時: string;
      対局場所: string;
      持ち時間: string;
      秒読み: string;
      秒加算: string;
    };
    tags: string[];
  }

  let formData: KifuFormData = {
    title: '',
    isPublic: false,
    gameInfo: {
      先手: '',
      後手: '',
      対局日時: '',
      対局場所: '',
      持ち時間: '',
      秒読み: '',
      秒加算: '',
    },
    tags: [],
  };
  let tagInput = '';
  let initialTimeInput = 0;
  let byoyomiTimeInput = 0;
  let incrementTimeInput = 0;

  const addTag = () => {
    if (tagInput && !formData.tags.includes(tagInput)) {
      formData.tags = [...formData.tags, tagInput];
      tagInput = '';
    }
  };
  const removeTag = (index: number) => {
    formData.tags = formData.tags.filter((_, i) => i !== index);
  };

  // 棋譜情報の更新
  async function handleUpdateKifuInfo() {
    if (!kifuId) return;
    formData.gameInfo['持ち時間'] = (initialTimeInput * 60).toString();
    formData.gameInfo['秒読み'] = byoyomiTimeInput.toString();
    formData.gameInfo['秒加算'] = incrementTimeInput.toString();
    const result = await updateKifuInfo(
      kifuId,
      formData.title,
      formData.isPublic,
      formData.gameInfo,
      formData.tags
    );
    if (result.ok) {
      await fetchKifuData();
    } else {
      console.error('Failed to update kifu info: ', result);
      isError = true;
    }
  }

  // ----------------------------------------
  // 棋譜の指し手管理

  let moves: KifuMove[] = [];

  const updateKifuMoves = async () => {
    // ToDo:
    console.debug('update kifu moves');
  };

  // ----------------------------------------
  let preinit = true;
  $: if ($account && preinit) {
    preinit = false;
    fetchKifuData();
  }
</script>

<div class="page">
  <section class="basic">
    <h2>棋譜の編集</h2>
    {#if kifu}
      <form onsubmit={handleUpdateKifuInfo} class="basic">
        <div class="form-group">
          <label for="title">タイトル</label>
          <input
            type="text"
            id="title"
            bind:value={formData.title}
            placeholder="棋譜のタイトル"
            required
          />
        </div>
        <div class="flex-arrange">
          <div class="form-group">
            <label for="black-player">先手</label>
            <input
              type="text"
              id="black-player"
              bind:value={formData.gameInfo.先手}
              placeholder="先手の対局者名"
            />
          </div>
          <div class="form-group">
            <label for="white-player">後手</label>
            <input
              type="text"
              id="white-player"
              bind:value={formData.gameInfo.後手}
              placeholder="後手の対局者名"
            />
          </div>
        </div>
        <div class="flex-arrange">
          <div class="form-group">
            <label for="start-date">対局日時</label>
            <input type="datetime-local" id="start-date" bind:value={formData.gameInfo.対局日時} />
          </div>
          <div class="form-group">
            <label for="place">対局場所</label>
            <input
              type="text"
              id="place"
              bind:value={formData.gameInfo.対局場所}
              placeholder="対局場所"
            />
          </div>
        </div>
        <div class="flex-arrange">
          <div class="form-group">
            <label for="initial-time">持ち時間（分）</label>
            <input type="number" id="initial-time" bind:value={initialTimeInput} min="0" />
          </div>
          <div class="form-group">
            <label for="byoyomi">秒読み（秒）</label>
            <input type="number" id="byoyomi" bind:value={byoyomiTimeInput} min="0" />
          </div>
          <div class="form-group">
            <label for="increment">加算（秒）</label>
            <input type="number" id="increment" bind:value={incrementTimeInput} min="0" />
          </div>
        </div>
        <div class="form-group">
          <h3 class="label">タグ</h3>
          <div class="tag-input-container">
            <input
              type="text"
              bind:value={tagInput}
              placeholder="タグを入力"
              onkeydown={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
            />
            <button type="button" onclick={addTag} class="add-tag-btn">追加</button>
          </div>
          <div class="tags-container">
            {#each formData.tags as tag, i}
              <span class="tag">
                {tag}
                <button type="button" onclick={() => removeTag(i)} class="remove-tag">×</button>
              </span>
            {/each}
          </div>
        </div>
        <div class="form-group">
          <h3 class="label">公開設定</h3>
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={formData.isPublic} />
            公開する
          </label>
        </div>
        <button type="submit" class="submit">棋譜情報を更新</button>
      </form>

      <KifuPlayer {kifu} />
      <button onclick={updateKifuMoves} class="submit">棋譜の指し手を更新</button>
    {/if}
  </section>
</div>

<style lang="scss">
  .tag-input-container {
    display: flex;
    gap: 0.5rem;

    input {
      flex: 1;
    }

    button {
      padding: 0.4rem 1rem;
      border: 1px solid var(--primary-color);
      border-radius: 0.4rem;
      background-color: var(--primary-color);
      color: white;

      &:hover,
      &:focus {
        border-color: var(--secondary-color);
        background-color: var(--secondary-color);
      }
    }
  }

  .tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-top: 0.5rem;

    .tag {
      background-color: var(--secondary-color);
      color: white;
      padding: 0.3rem 0.4rem 0.3rem 0.6rem;
      border-radius: 0.3rem;
      display: flex;
      align-items: center;
      gap: 0.4rem;
      font-size: 0.9rem;
      line-height: 1.4rem;

      .remove-tag {
        background: none;
        border: none;
        color: white;
        cursor: pointer;
      }
    }
  }
</style>
