<!-- src/lib/components/PositionView/PieceStand.svelte -->

<script lang="ts">
  import { PieceType, type PieceClickEvent } from '$lib/types/Piece';
  import Piece from './Piece.svelte';

  export let hands: Map<PieceType, number>;
  export let isBlack: boolean;
  export let x: number;
  export let y: number;
  export let style: 'default' | 'pentagon' | undefined = 'default';
  export let onAreaClick: ((e: PieceClickEvent) => void) | undefined = undefined;
  export let onPieceClick: ((e: PieceClickEvent) => void) | undefined = undefined;
  export let pickedPiece: PieceClickEvent | undefined;
  const w = 670; // 5 + 30 + (240 + 60)*2 + 30 + 5
  const h = 1075; // 260*4 + 5*3 + 10*2
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
  const buildClickEvent = (pieceType: PieceType): PieceClickEvent => {
    const event = {
      pieceType: pieceType,
      source: {
        type: 'stand',
        isBlack: isBlack,
      },
    } as PieceClickEvent;
    return event;
  };

  const isPickedPiece = (pieceType: PieceType): boolean => {
    if (!pickedPiece) return false;
    if (pickedPiece.source.type !== 'stand') return false;
    if (pickedPiece.source.isBlack !== isBlack) return false;
    if (pickedPiece.pieceType !== pieceType) return false;
    return true;
  };
  const pieceState = (pieceType: PieceType): 'normal' | 'picked' => {
    if (isPickedPiece(pieceType)) {
      return 'picked';
    }
    return 'normal';
  };
</script>

<g transform={`translate(${x}, ${y}) ${isBlack ? '' : `rotate(180,${w / 2},${h / 2})`}`}>
  <rect width={w} height={h} fill="#666" />
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <rect
    x="10"
    y="10"
    width={w - 20}
    height={h - 20}
    fill="#ffc"
    onclick={() => {
      // 駒を駒台に移動させる用のイベント
      onAreaClick && onAreaClick(buildClickEvent(PieceType.VACANCY));
    }}
    role="button"
    aria-label={`piece-stand-${isBlack ? 'black' : 'white'}`}
    tabindex="0"
  />
  {#each pieceTypeList as pieceType, i}
    {@const count = hands.get(pieceType) ?? 0}
    {#if count}
      <g transform={`translate(${35 + (i % 2) * 300}, ${5 + Math.floor(i / 2) * 265})`}>
        <Piece
          {pieceType}
          {style}
          onClick={pickedPiece
            ? undefined
            : () => onPieceClick && onPieceClick(buildClickEvent(pieceType))}
          useViewBox={false}
          state={pieceState(pieceType)}
        />
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
