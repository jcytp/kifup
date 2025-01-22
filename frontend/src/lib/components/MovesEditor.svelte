<!-- src/lib/components/MovesEditor.svelte -->

<script lang="ts">
  import { generateMovedSfen, isBlackOfSfen } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import { PieceType } from '$lib/types/Piece';
  import PositionView from './PositionView/PositionView.svelte';

  // callback
  export let onChange: (moves: KifuMove[]) => void;

  // parameters
  export let initialSfen: string | undefined;
  export let moveList: KifuMove[];
  $: moveNumber = moveList.length;
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

  $: if (moveList) {
    console.debug('MovesEditor moveList:', moveList);
  }

  const handleAppendMove = (num: number, move: KifuMove) => {
    const newMoveList = moveList.slice(0, num - 1);
    newMoveList.push(move);
    moveNumber++;
    onChange(newMoveList);
  };
  const handlePromote = (promote: boolean) => {
    if (!moveList.length) {
      return;
    }
    moveList[moveList.length - 1].promote = promote;
    if (promote) {
      moveList[moveList.length - 1].piece = moveList[moveList.length - 1].piece | PieceType.PROMOTE;
    }
  };
  const handleChangeComment = (comment: string) => {
    if (moveNumber >= 1) {
      moveList[moveNumber - 1].comment = comment;
    } else {
      // ToDo: initial comment
    }
  };
</script>

<div class="position-editor">
  <PositionView
    mode="moves"
    {isBlackFirst}
    sfen={currentSfen}
    {comment}
    {moveList}
    {moveNumber}
    onAppendMove={handleAppendMove}
    onPromote={handlePromote}
    onChangeComment={handleChangeComment}
  />
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
