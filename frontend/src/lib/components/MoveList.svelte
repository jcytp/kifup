<!-- src/lib/components/PieceStand.svelte -->

<script lang="ts">
  import type { KifuMove } from '$lib/types/Kifu';

  export let moveList: KifuMove[];
  export let moveNumber: number;
  export let x: number;
  export let y: number;
  const w = 920;
  const h = 2400;
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={w} height={h} fill="#666" />
  <rect x="10" y="10" width={w - 20} height={h - 20} fill="#ffc" />
</g>
