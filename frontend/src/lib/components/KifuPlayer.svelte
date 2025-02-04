<!-- src/lib/components/KifuPlayer.svelte -->

<script lang="ts">
  import { generateMovedSfen, isBlackOfSfen } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import {
    CameraIcon,
    SkipBackIcon,
    SkipForwardIcon,
    StepBackIcon,
    StepForwardIcon,
  } from 'lucide-svelte';
  import PositionView from './PositionView/PositionView.svelte';
  import { viewport } from '$lib/stores/viewport';

  export let initialSfen: string | undefined;
  export let moveList: KifuMove[];
  let moveNumber: number = 0;
  $: comment = moveNumber >= 1 ? moveList[moveNumber - 1].comment || '' : '';
  $: currentSfen = generateMovedSfen(initialSfen, moveList, moveNumber);
  $: isBlackFirst = isBlackOfSfen(initialSfen);
  let positionView: PositionView;
  let iconSize = $viewport.isMobile ? 16 : 20;

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
  let handleDownloadPng = async () => {
    if (positionView) {
      await positionView.exportAsPng();
    }
  };
</script>

<div class="kifu-player">
  <PositionView
    bind:this={positionView}
    mode="replay"
    {isBlackFirst}
    sfen={currentSfen}
    {comment}
    {moveList}
    {moveNumber}
  />
  <div class="controls">
    <button onclick={handleToStart}><SkipBackIcon size={iconSize} /></button>
    <button onclick={handleToPrev}><StepBackIcon size={iconSize} /></button>
    <button onclick={handleToNext}><StepForwardIcon size={iconSize} /></button>
    <button onclick={handleToEnd}><SkipForwardIcon size={iconSize} /></button>
    <button onclick={handleDownloadPng}><CameraIcon size={iconSize} /></button>
  </div>
</div>

<style lang="scss">
  @import '../styles/mixins.scss';

  .controls {
    display: flex;
    justify-content: center;
    gap: 0.4rem;
    margin-top: 0.5rem;

    button {
      @include sp {
        width: 4rem;
        font-size: small;
        padding: 0.2rem 0;
      }
      @include pc {
        width: 6rem;
        padding: 0.4rem 0;
      }
      display: flex;
      justify-content: center;
      align-items: center;

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
