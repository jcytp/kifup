<!-- src/lib/components/MovesEditor.svelte -->

<script lang="ts">
  import { BoardPosition } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import { PieceType } from '$lib/types/Piece';
  import BoardPositionView from './PositionView/PositionView.svelte';

  // callback
  export let onChange: (moves: KifuMove[]) => void;

  // parameters
  export let initialSfen: string | undefined;
  export let moveList: KifuMove[];
  $: moveNumber = moveList.length;
  $: currentSfen = generateMovedSfen(initialSfen, moveList, moveNumber);
  $: isBlackFirst = isBlackOfSfen(initialSfen);

  const generateMovedSfen = (
    sfen: string | undefined,
    moves: KifuMove[],
    num: number
  ): string | undefined => {
    const position = new BoardPosition(sfen);
    for (let i = 0; i < num; i++) {
      if (i >= moves.length) {
        console.error('move number is over moves count');
      }
      position.next(moves[i]);
    }
    return position.toSfen(1); // 手数は1固定
  };

  const isBlackOfSfen = (sfen?: string) => {
    if (!sfen) return true;
    const parts = sfen.split(' ');
    return parts[1] === 'b';
  };

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
</script>

<div class="position-editor">
  <BoardPositionView
    mode="moves"
    {isBlackFirst}
    sfen={currentSfen}
    {moveList}
    {moveNumber}
    onAppendMove={handleAppendMove}
    onPromote={handlePromote}
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
