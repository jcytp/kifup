<!-- src/lib/components/KifuPlayer.svelte -->

<script lang="ts">
  import { generateMovedSfen, isBlackOfSfen } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import PositionView from './PositionView/PositionView.svelte';

  export let initialSfen: string | undefined;
  export let moveList: KifuMove[];
  let moveNumber: number = 0;
  $: comment = moveNumber >= 1 ? moveList[moveNumber - 1].comment || '' : '';
  $: currentSfen = generateMovedSfen(initialSfen, moveList, moveNumber);
  $: isBlackFirst = isBlackOfSfen(initialSfen);

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
</script>

<div class="kifu-player">
  <PositionView mode="replay" {isBlackFirst} sfen={currentSfen} {comment} {moveList} {moveNumber} />
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
