<!-- src/lib/components/BoardPositionView.svelte -->

<script lang="ts">
  import { BoardPosition } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import { PieceType } from '$lib/types/Piece';
  import Board from './Board.svelte';
  import MoveList from './MoveList.svelte';
  import Piece from './Piece.svelte';
  import PieceBox from './PieceBox.svelte';
  import PieceStand from './PieceStand.svelte';

  // callbacks
  export let onChange: (position?: string, moves?: []) => void;

  // parameters
  export let mode: 'position' | 'moves' = 'position';
  export let sfen: string | undefined;
  export let moveList: KifuMove[] = [];
  export let moveNumber: number = 0;
  $: visibleMoveList = mode === 'moves';
  $: visiblePieceBox = mode === 'position';
  $: position = new BoardPosition(sfen);
  console.debug(position);

  // $: sizeSet = {
  //   borderThin: 5,
  //   borderThic: 10,
  //   piecePadding: 60,
  //   boardMargin: 120,
  //   subBoxMargin: 60,
  //   cell: { w: 240, h: 260 },
  //   counter: { w: 120, h: 120 },
  //   board: {
  //     w: 2220, // 240 * 9 + 5 * 8 + 10 * 2
  //     h: 2400, // 260 * 9 + 5 * 8 + 10 * 2
  //     x: 1100 + (visibleMoveList ? 980 : 0), // 60 + 920 + 120
  //     y: 120,
  //   },
  //   whiteStand: {
  //     w: 920, // (240 + 120) * 2 + 60 * 3 + 10 * 2
  //     h: 1360, // 260 * 4 + 60 * 5 + 10 * 2
  //     x: 60 + (visibleMoveList ? 980 : 0),
  //     y: 120,
  //   },
  //   blackStand: {
  //     w: 920,
  //     h: 1360,
  //     x: 3440 + (visibleMoveList ? 980 : 0), // 60 + 920 + 120 + 2220 + 120
  //     y: 1160, // 120 + 2400 - 1360
  //   },
  //   pieceBox: {
  //     w: 920, // (240 + 120) * 2 + 60 * 3 + 10 * 2
  //     h: 1680, // 260 * 5 + 60 * 6 + 10 * 2
  //     x: 4420 + (visibleMoveList ? 980 : 0), // 60 + 920 + 120 + 2220 + 120 + 920 + 60
  //     y: 840, // 120 + 2400 - 1680
  //   },
  //   moveList: {
  //     w: 920,
  //     h: 2400,
  //     x: 60,
  //     y: 120,
  //   },
  //   viewBox: {
  //     w: 4420 + (visibleMoveList ? 980 : 0) + (visiblePieceBox ? 980 : 0), // 2220 + 920 + 920 + 120 * 2 + 60 * 2 + ? + ?
  //     h: 2640, // 2400 + 120 * 2
  //   },
  // };

  $: viewBoxWidth = 4420 + (visibleMoveList ? 980 : 0) + (visiblePieceBox ? 980 : 0);
  const viewBoxHeight = 2640;
</script>

<svg class="board-position-view-svg" viewBox={`0 0 ${viewBoxWidth} ${viewBoxHeight}`}>
  <rect width="100%" height="100%" fill="#fff" />

  {#if visibleMoveList}
    <MoveList x={60} y={120} {moveList} {moveNumber} />
  {/if}
  <g transform={`translate(${visibleMoveList ? 980 : 0}, 0)`}>
    <PieceStand x={60} y={120} isBlack={false} hands={position.whiteHands} />
    <Board x={1100} y={120} blackBoard={position.blackBoard} whiteBoard={position.whiteBoard} />
    <PieceStand x={3440} y={1160} isBlack={true} hands={position.blackHands} />
    {#if visiblePieceBox}
      <PieceBox x={4420} y={840} pieces={position.pieceBox} />
    {/if}
  </g>
</svg>
