<!-- src/lib/components/PieceStand.svelte -->

<script lang="ts">
  import type { PieceType } from '$lib/types/Piece';

  export let hands: Map<PieceType, number>;
  export let isBlack: boolean;
  export let x: number;
  export let y: number;
  // const x = isBlack ? 3440 : 60;
  // const y = isBlack ? 1160 : 120;
  const w = 920;
  const h = 1360;
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={w} height={h} fill="#666" />
  <rect x="10" y="10" width={w - 20} height={h - 20} fill="#ffc" />
</g>
