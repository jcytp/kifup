<!-- src/lib/components/ShogiPiece.svelte -->

<script lang="ts">
  import { PieceChar, PieceType } from '$lib/types/KifuPiece';

  export let pieceType: PieceType;
  export let style: 'default' | 'pentagon' = 'default';
  export let reverse: boolean;
  export let onClick: () => void | undefined;

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
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<svg
  xmlns="http://www.w3.org/2000/svg"
  viewBox="0 0 240 260"
  class="piece-svg"
  on:click={onClick}
  role="button"
  aria-label={`${pieceChar}`}
  tabindex="0"
>
  <g {transform}>
    {#if style == 'pentagon'}
      <path d={path} fill="#fff" fill-opacity="0.6" stroke="#333" stroke-width="10" />
      <text x="120" y="160" text-anchor="middle" fill="#333" font-size="120">
        {pieceChar}
      </text>
    {:else}
      <text x="120" y="140" text-anchor="middle" fill="#333" font-size="120">
        {pieceChar}
      </text>
    {/if}
  </g>
</svg>

{#if onClick != undefined}
  <style lang="scss">
    .piece-svg {
      cursor: pointer;
      &:hover {
        opacity: 0.8;
      }
    }
  </style>
{/if}

<style lang="scss">
  .piece-svg {
    width: 100%;
    height: 100%;
  }
</style>
