<!-- src/lib/components/PieceStand.svelte -->

<script lang="ts">
  import type { PieceType } from '$lib/types/Piece';

  export let pieces: Map<PieceType, number>;
  export let x: number;
  export let y: number;
  // const x = 4420;
  // const y = 840;
  const w = 920;
  const h = 1680;
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={w} height={h} fill="#666" />
  <rect x="10" y="10" width={w - 20} height={h - 20} fill="#ffc" />
</g>
