<!-- src/lib/components/MovesEditor.svelte -->

<script lang="ts">
  import type { KifuMove } from '$lib/types/Kifu';
  import BoardPositionView from './PositionView/PositionView.svelte';

  // callback
  export let onChange: (moves: KifuMove[]) => void;

  // parameters
  export let initial: string | undefined;
  export let moveList: KifuMove[];

  let moveNumber = 0;

  let handleToStart = () => {
    moveNumber = 0;
  };
  let handleToPrev = () => {
    if (moveNumber > 0) moveNumber--;
  };
  let handleToNext = () => {
    if (moveNumber < moveList.length) moveNumber++;
  };
  let handleToEnd = () => {
    moveNumber = moveList.length;
  };

  const handleChange = (position?: string, moves?: KifuMove[]) => {
    if (!moves) return;
    onChange(moves);
  };
</script>

<div class="position-editor">
  <BoardPositionView mode="moves" sfen={initial} {moveList} {moveNumber} onChange={handleChange} />
  <div class="controls">
    <button onclick={handleToStart}>|◀</button>
    <button onclick={handleToPrev}>◀</button>
    <button onclick={handleToNext}>▶</button>
    <button onclick={handleToEnd}>▶|</button>
  </div>
</div>

<style>
  .controls {
    margin-top: 0.5rem;
    text-align: center;

    button {
      width: 6rem;
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
