<!-- src/lib/components/PieceStand.svelte -->

<script lang="ts">
  import { PieceType } from '$lib/types/Piece';
  import Piece from './Piece.svelte';

  export let hands: Map<PieceType, number>;
  export let isBlack: boolean;
  export let x: number;
  export let y: number;
  export let style: 'default' | 'pentagon' | undefined = 'default';
  const w = 800;
  const h = 1360;
  const pieceTypeList = [
    PieceType.FU,
    PieceType.VACANCY,
    PieceType.KY,
    PieceType.KE,
    PieceType.GI,
    PieceType.KI,
    PieceType.KA,
    PieceType.HI,
  ];
  const handleClickPiece = () => {};
</script>

<g transform={`translate(${x}, ${y}) ${isBlack ? '' : `rotate(180,${w / 2},${h / 2})`}`}>
  <rect width={w} height={h} fill="#666" />
  <rect x="10" y="10" width={w - 20} height={h - 20} fill="#ffc" />
  {#each pieceTypeList as pieceType, i}
    {@const count = hands.get(pieceType) ?? 0}
    {#if count}
      <g transform={`translate(${70 + (i % 2) * 360}, ${70 + Math.floor(i / 2) * 320})`}>
        <Piece {pieceType} {style} onClick={handleClickPiece} useViewBox={false} />
        {#if count > 1}
          <text
            x={220}
            y={190}
            fill="#333"
            dominant-baseline="middle"
            text-anchor="start"
            font-size={140}
            style:font-size="140px">{hands.get(pieceType)}</text
          >
        {/if}
      </g>
    {/if}
  {/each}
</g>
