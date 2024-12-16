<!-- src/lib/components/Board.svelte -->

<script lang="ts">
  import { PieceType } from '$lib/types/Piece';
  import Piece from './Piece.svelte';

  export let blackBoard: PieceType[][];
  export let whiteBoard: PieceType[][];
  export let x: number;
  export let y: number;

  $: console.debug('blackBoard:', blackBoard);
  $: console.debug('whiteBoard:', whiteBoard);
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={2220} height={2400} fill="#666" />

  {#each blackBoard as row, i}
    {#each row as piece, j}
      <rect x={10 + j * 245} y={10 + i * 265} width="240" height="260" fill="#ffc" />
      {#if piece !== PieceType.VACANCY}
        <g transform={`translate(${10 + j * 245}, ${10 + i * 265})`}>
          <Piece
            pieceType={piece}
            style="pentagon"
            reverse={false}
            onClick={undefined}
            useViewBox={false}
          />
        </g>
      {/if}
      {#if whiteBoard[i][j] !== PieceType.VACANCY}
        <g transform={`translate(${10 + j * 245}, ${10 + i * 265})`}>
          <Piece
            pieceType={whiteBoard[i][j]}
            style="default"
            reverse={true}
            onClick={undefined}
            useViewBox={false}
          />
        </g>
      {/if}
    {/each}
  {/each}
</g>
