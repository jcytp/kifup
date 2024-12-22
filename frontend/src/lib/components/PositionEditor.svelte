<!-- src/lib/components/PositionEditor.svelte -->

<script lang="ts">
  import type { KifuMove } from '$lib/types/Kifu';
  import BoardPositionView from './PositionView/PositionView.svelte';

  // callback
  export let onChange: (position?: string) => void;

  let boardPosition: string | undefined;

  const handleChange = (position?: string, moves?: KifuMove[]) => {
    boardPosition = position;
    onChange(position);
  };

  const reset = () => {
    // TEST POSITION
    // boardPosition = 'lnsg5/9/ppppp4/9/9/9/PPPPPPP2/9/4KGSNL b 2PLNSGBR2plnsgbr 1';
    boardPosition = undefined;
    onChange(boardPosition);
  };

  const allRemove = () => {
    boardPosition = '9/9/9/9/9/9/9/9/9 b - 1';
    onChange(boardPosition);
  };
</script>

<div class="position-editor">
  <BoardPositionView mode="position" sfen={boardPosition} onChange={handleChange} />
  <div class="controls">
    <button onclick={reset}>平手初形</button>
    <button onclick={allRemove}>全て駒箱</button>
  </div>
</div>

<style>
  .controls {
    margin-top: 0.5rem;
    text-align: center;

    button {
      width: 10rem;
      padding: 0.3rem 0.5rem;
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
