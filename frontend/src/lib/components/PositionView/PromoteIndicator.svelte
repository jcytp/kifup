<!-- src/lib/components/PositionView/TurnIndicator.svelte -->

<script lang="ts">
  import { PieceType } from '$lib/types/Piece';
  import Piece from './Piece.svelte';

  export let x: number;
  export let y: number;
  export let pieceType: PieceType;
  export let onSelect: ((promote: boolean) => void) | undefined = undefined;
  const style = 'pentagon';
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={900} height={620} rx={30} ry={30} fill="#666" />
  <rect x={10} y={10} width={880} height={600} rx={30} ry={30} fill="#cfc" />

  <text
    x={450}
    y={150}
    dominant-baseline="middle"
    text-anchor="middle"
    fill="#333"
    font-size="90"
    style:font-size="90px">成りますか？</text
  >

  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <g transform={`translate(${150}, ${280})`}>
    <Piece {pieceType} {style} useViewBox={false} onClick={() => onSelect && onSelect(false)} />
  </g>
  <g transform={`translate(${510}, ${280})`}>
    <Piece
      pieceType={pieceType | PieceType.PROMOTE}
      {style}
      useViewBox={false}
      onClick={() => onSelect && onSelect(true)}
    />
  </g>
</g>
