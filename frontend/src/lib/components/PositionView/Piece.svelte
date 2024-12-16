<!-- src/lib/components/PositionView/Piece.svelte -->

<script lang="ts">
  import { PieceChar, PieceType } from '$lib/types/Piece';

  export let pieceType: PieceType;
  export let style: 'default' | 'pentagon' = 'default';
  export let reverse: boolean = false;
  export let onClick: (() => void) | undefined = undefined;
  export let useViewBox = true;
  export let state: 'normal' | 'picked' = 'normal';

  const L_SIZE_PATH = 'M120,10 L210,40 L230,250 L10,250, L30,40 Z';
  const M_SIZE_PATH = 'M120,20 L200,50 L230,250 L10,250, L40,50 Z';
  const S_SIZE_PATH = 'M120,30 L190,60 L220,250 L20,250, L50,60 Z';
  const SS_SIZE_PATH = 'M120,30 L190,60 L210,240 L30,240, L50,60 Z';
  const piecePath = new Map<PieceType, string>([
    [PieceType.FU, SS_SIZE_PATH],
    [PieceType.KY, S_SIZE_PATH],
    [PieceType.KE, S_SIZE_PATH],
    [PieceType.GI, M_SIZE_PATH],
    [PieceType.KI, M_SIZE_PATH],
    [PieceType.KA, L_SIZE_PATH],
    [PieceType.HI, L_SIZE_PATH],
    [PieceType.OU, L_SIZE_PATH],
    [PieceType.TO, SS_SIZE_PATH],
    [PieceType.NY, S_SIZE_PATH],
    [PieceType.NK, S_SIZE_PATH],
    [PieceType.NG, M_SIZE_PATH],
    [PieceType.UM, L_SIZE_PATH],
    [PieceType.RY, L_SIZE_PATH],
  ]);
  $: path = piecePath.get(pieceType);
  $: pieceChar = PieceChar.get(pieceType) ?? '?';
  $: transform = reverse ? 'rotate(180,120,130)' : undefined;
  $: opacity = state == 'picked' ? 0.5 : 1;
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<svg
  xmlns="http://www.w3.org/2000/svg"
  viewBox={useViewBox ? '0 0 240 260' : undefined}
  class="piece-svg"
  on:click={onClick}
  role="button"
  aria-label={`${pieceChar}`}
  tabindex="0"
  pointer-events={onClick ? 'all' : 'none'}
>
  <g {transform}>
    {#if style == 'pentagon'}
      <path d={path} fill="#ffd" stroke="#666" stroke-width="10" {opacity} />
      <text
        x={120}
        y={165}
        dominant-baseline="middle"
        text-anchor="middle"
        fill="#333"
        {opacity}
        font-size={140}
        style:font-size="140px"
      >
        {pieceChar}
      </text>
    {:else}
      <text
        x={120}
        y={165}
        dominant-baseline="middle"
        text-anchor="middle"
        fill="#333"
        {opacity}
        font-size={180}
        style:font-size="180px"
      >
        {pieceChar}
      </text>
    {/if}
  </g>
</svg>

<style lang="scss">
  .piece-svg {
    cursor: pointer;
    &:hover {
      opacity: 0.8;
    }
  }
</style>
