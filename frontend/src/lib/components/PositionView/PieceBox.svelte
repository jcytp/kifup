<!-- src/lib/components/PositionView/PieceBox.svelte -->

<script lang="ts">
  import { PieceType, type PieceClickEvent } from '$lib/types/Piece';
  import Piece from './Piece.svelte';

  export let pieces: Map<PieceType, number>;
  export let x: number;
  export let y: number;
  export let style: 'default' | 'pentagon' | undefined = 'default';
  export let onAreaClick: ((e: PieceClickEvent) => void) | undefined = undefined;
  export let onPieceClick: ((e: PieceClickEvent) => void) | undefined = undefined;
  export let pickedPiece: PieceClickEvent | undefined;
  const w = 800;
  const h = 1680;
  const pieceTypeList = [
    PieceType.FU,
    PieceType.VACANCY,
    PieceType.KY,
    PieceType.KE,
    PieceType.GI,
    PieceType.KI,
    PieceType.KA,
    PieceType.HI,
    PieceType.OU,
  ];
  const buildClickEvent = (pieceType: PieceType): PieceClickEvent => {
    const event = {
      pieceType: pieceType,
      source: {
        type: 'box',
      },
    } as PieceClickEvent;
    return event;
  };
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={w} height={h} fill="#666" />
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <rect
    x="10"
    y="10"
    width={w - 20}
    height={h - 20}
    fill="#ccc"
    onclick={() => {
      // 駒を駒箱に移動させる用のイベント
      onAreaClick && onAreaClick(buildClickEvent(PieceType.VACANCY));
    }}
    role="button"
    aria-label={`piece-box`}
    tabindex="0"
  />
  {#each pieceTypeList as pieceType, i}
    {@const count = pieces.get(pieceType) ?? 0}
    {#if count}
      <g transform={`translate(${70 + (i % 2) * 360}, ${70 + Math.floor(i / 2) * 320})`}>
        <Piece
          {pieceType}
          {style}
          onClick={pickedPiece
            ? undefined
            : () => onPieceClick && onPieceClick(buildClickEvent(pieceType))}
          useViewBox={false}
        />
        {#if count > 1}
          <text
            x={220}
            y={190}
            fill="#333"
            dominant-baseline="middle"
            text-anchor="start"
            font-size={140}
            style:font-size="140px">{pieces.get(pieceType)}</text
          >
        {/if}
      </g>
    {/if}
  {/each}
</g>
