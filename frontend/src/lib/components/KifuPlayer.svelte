<!-- src/lib/components/KifuPlayer.svelte -->

<script lang="ts">
  import type { Kifu, Move } from '$lib/types/Kifu';
  import ShogiBoard from './old_ShogiBoard.svelte';
  import { generatePosition } from '$lib/utils/positionControl';

  export let kifu: Kifu;

  // 現在の手数
  let currentMoveIndex = -1; // -1は初期局面を表す

  // 現在の手の情報を取得
  $: currentMove = currentMoveIndex >= 0 ? kifu.moves[currentMoveIndex] : null;

  // 現在の局面を生成
  $: currentPosition = kifu.initialPosition
    ? generatePosition(kifu.initialPosition, kifu.moves, currentMoveIndex)
    : undefined;

  // 操作関数
  function goToFirst() {
    currentMoveIndex = -1;
  }

  function goToPrev() {
    if (currentMoveIndex > -1) {
      currentMoveIndex--;
    }
  }

  function goToNext() {
    if (currentMoveIndex < kifu.moves.length - 1) {
      currentMoveIndex++;
    }
  }

  function goToLast() {
    currentMoveIndex = kifu.moves.length - 1;
  }

  // 指し手のテキスト表示用
  function formatMove(move: Move): string {
    const moveNum = String(move.moveNumber).padStart(3, ' ');
    return `${moveNum} ${move.piece} (${move.to.x}, ${move.to.y})`;
  }
</script>

<div class="kifu-player">
  <!-- 左側：指し手リスト -->
  <div class="moves-list">
    <div
      class="move-item initial"
      class:current={currentMoveIndex === -1}
      on:click={() => (currentMoveIndex = -1)}
      on:keydown={(e) => e.key === 'Enter' && (currentMoveIndex = -1)}
      tabindex="0"
      role="button"
    >
      開始局面
    </div>
    {#each kifu.moves as move, index}
      <div
        class="move-item"
        class:current={currentMoveIndex === index}
        on:click={() => (currentMoveIndex = index)}
        on:keydown={(e) => e.key === 'Enter' && (currentMoveIndex = index)}
        tabindex="0"
        role="button"
      >
        {formatMove(move)}
      </div>
    {/each}
  </div>

  <!-- 右側：局面表示エリア -->
  <div class="board-area">
    <!-- 局面表示 -->
    <div class="board-placeholder">
      <div class="board-area">
        <ShogiBoard position={currentPosition} />
      </div>
    </div>

    <!-- 局面コメント -->
    <div class="comment-area">
      {#if currentMove?.comment}
        <p>{currentMove.comment}</p>
      {:else}
        <p class="no-comment">この局面にコメントはありません</p>
      {/if}
    </div>

    <!-- 操作ボタン -->
    <div class="control-buttons">
      <button on:click={goToFirst} disabled={currentMoveIndex === -1} class="control-button">
        |◀
      </button>
      <button on:click={goToPrev} disabled={currentMoveIndex === -1} class="control-button">
        ◀
      </button>
      <button
        on:click={goToNext}
        disabled={currentMoveIndex === kifu.moves.length - 1}
        class="control-button"
      >
        ▶
      </button>
      <button
        on:click={goToLast}
        disabled={currentMoveIndex === kifu.moves.length - 1}
        class="control-button"
      >
        ▶|
      </button>
    </div>
  </div>
</div>

<style lang="scss">
  .kifu-player {
    display: flex;
    gap: 2rem;
    background: white;
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .moves-list {
    flex: 0 0 200px;
    height: 600px;
    overflow-y: auto;
    border: 1px solid var(--border-color);
    border-radius: 4px;

    .move-item {
      padding: 0.5rem 1rem;
      cursor: pointer;
      border-bottom: 1px solid var(--border-color);
      font-family: monospace;

      &:hover {
        background-color: rgba(0, 0, 0, 0.05);
      }

      &.current {
        background-color: var(--primary-color);
        color: white;
      }

      &:last-child {
        border-bottom: none;
      }
    }
  }

  .board-area {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .comment-area {
    min-height: 100px;
    padding: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: #f8f8f8;

    .no-comment {
      color: #666;
      font-style: italic;
    }
  }

  .control-buttons {
    display: flex;
    gap: 0.5rem;
    justify-content: center;

    .control-button {
      padding: 0.5rem 1rem;
      border: 1px solid var(--border-color);
      border-radius: 4px;
      background: white;
      font-family: monospace;

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }

      &:not(:disabled):hover {
        background-color: var(--secondary-color);
        color: white;
        border-color: var(--secondary-color);
      }
    }
  }
</style>
