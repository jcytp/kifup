<!-- src/lib/components/PositionEditor.svelte -->

<script lang="ts">
  import type { KifuMove } from '$lib/types/Kifu';
  import BoardPositionView from './BoardPositionView.svelte';

  // callback
  export let onChange: (position?: string) => void;

  let boardPosition: string | undefined;

  const handleChange = (position?: string, moves?: KifuMove[]) => {
    onChange(position);
  };

  const reset = () => {
    boardPosition = undefined;
    onChange(boardPosition);
  };
</script>

<div class="position-editor">
  <BoardPositionView mode="position" sfen={boardPosition} onChange={handleChange} />
  <div class="controls">
    <button on:click={reset}>初期局面に戻す</button>
  </div>
</div>

<style>
  .controls {
    margin-top: 0.5rem;
    text-align: center;

    button {
      padding: 0.5rem 1rem;
      background-color: var(--primary-color);
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;

      &:hover {
        opacity: 0.9;
      }
    }
  }
</style>
