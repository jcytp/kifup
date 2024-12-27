<!-- src/lib/components/PositionView/TurnIndicator.svelte -->

<script lang="ts">
  export let x: number;
  export let y: number;
  export let isBlackTurn: boolean;
  export let onClick: (() => void) | undefined = undefined;
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<g
  transform={`translate(${x}, ${y})`}
  class="turn-indicator"
  on:click={() => onClick && onClick()}
  role="button"
  tabindex="0"
>
  <rect width="800" height="160" rx="30" ry="30" fill="#666" />
  <rect x="10" y="10" width="780" height="140" rx="20" ry="20" fill="#cfc" />
  <text
    x="400"
    y="90"
    dominant-baseline="middle"
    text-anchor="middle"
    fill="#333"
    font-size="90"
    style:font-size="90px"
  >
    {isBlackTurn ? '▲先手番' : '△後手番'}
  </text>
</g>

<style>
  .turn-indicator {
    cursor: pointer;
  }
  .turn-indicator:hover rect:nth-child(1) {
    fill: #696;
  }
</style>
