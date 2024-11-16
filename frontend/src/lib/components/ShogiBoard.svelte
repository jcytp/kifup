<!-- src/lib/components/ShogiBoard.svelte -->

<script lang="ts">
  import type { BoardPosition, PieceType } from '$lib/types/Kifu';

  // プロパティの定義
  export let position: BoardPosition | undefined = undefined;
  export let reversed: boolean = false;  // 盤面の向き（true: 後手視点）

  // 将棋の駒の文字表示用マッピング
  const pieceKanji: { [K in PieceType]: string } = {
    '歩': '歩', '香': '香', '桂': '桂', '銀': '銀', '金': '金',
    '角': '角', '飛': '飛', '玉': '玉',
    'と': 'と', '成香': '杏', '成桂': '圭', '成銀': '全',
    '馬': '馬', '龍': '龍'
  };

  // 盤面のサイズ設定
  const boardSize = 400;  // 盤面の幅（px）
  const cellSize = boardSize / 9;  // マスの幅
  const komadaiWidth = cellSize * 1.5;  // 駒台の幅

  // 駒台の持ち駒を集計する関数
  function countPieces(hands: { [K in PieceType]?: number } | undefined) {
    if (!hands) return [];
    return Object.entries(hands).map(([piece, count]) => ({
      piece: piece as PieceType,
      count: count || 0
    }));
  }

  // 座標変換関数（reversed対応）
  function transformX(x: number): number {
    return reversed ? 10 - x : x;
  }

  function transformY(y: number): number {
    return reversed ? 10 - y : y;
  }
</script>

<div class="board-container">
  <svg
    viewBox={`0 0 ${boardSize + komadaiWidth * 2} ${boardSize}`}
    class="board-svg"
  >
    <!-- 後手の駒台 -->
    <g transform={`translate(0, 0)`}>
      <rect
        x="0"
        y="0"
        width={komadaiWidth}
        height={boardSize}
        class="komadai"
      />
      {#if position?.hands?.white}
        {#each countPieces(position.hands.white) as {piece, count}, i}
          <g transform={`translate(${komadaiWidth/2}, ${cellSize * (1 + i)})`}>
            <text class="piece-text white">{pieceKanji[piece]}</text>
            {#if count > 1}
              <text class="count-text">x{count}</text>
            {/if}
          </g>
        {/each}
      {/if}
    </g>

    <!-- 盤面 -->
    <g transform={`translate(${komadaiWidth}, 0)`}>
      <!-- マス目 -->
      {#each Array(9) as _, i}
        {#each Array(9) as _, j}
          <rect
            x={cellSize * i}
            y={cellSize * j}
            width={cellSize}
            height={cellSize}
            class="cell"
          />
        {/each}
      {/each}

      <!-- 駒の配置 -->
      {#if position?.pieces}
        {#each position.pieces as {position: pos, piece, isBlack}}
          <g transform={`translate(
            ${cellSize * (transformX(pos.x) - 0.5)},
            ${cellSize * (transformY(pos.y) - 0.5)}
          )`}>
            <text
              class="piece-text"
              class:black={isBlack}
              class:white={!isBlack}
              transform={reversed ? "rotate(180)" : ""}
            >
              {pieceKanji[piece]}
            </text>
          </g>
        {/each}
      {/if}
    </g>

    <!-- 先手の駒台 -->
    <g transform={`translate(${boardSize + komadaiWidth}, 0)`}>
      <rect
        x="0"
        y="0"
        width={komadaiWidth}
        height={boardSize}
        class="komadai"
      />
      {#if position?.hands?.black}
        {#each countPieces(position.hands.black) as {piece, count}, i}
          <g transform={`translate(${komadaiWidth/2}, ${cellSize * (1 + i)})`}>
            <text class="piece-text black">{pieceKanji[piece]}</text>
            {#if count > 1}
              <text class="count-text">x{count}</text>
            {/if}
          </g>
        {/each}
      {/if}
    </g>
  </svg>
</div>

<style lang="scss">
  .board-container {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
  }

  .board-svg {
    width: 100%;
    height: auto;
  }

  .cell {
    fill: #fff;
    stroke: #000;
    stroke-width: 1;
  }

  .komadai {
    fill: #f0f0f0;
    stroke: #000;
    stroke-width: 1;
  }

  .piece-text {
    text-anchor: middle;
    // dominant-baseline: central;
    font-size: 24px;
    font-weight: bold;
    user-select: none;

    &.black {
      fill: #000;
    }

    &.white {
      fill: #666;
      transform: rotate(180deg);
    }
  }

  .count-text {
    text-anchor: middle;
    // dominant-baseline: central;
    font-size: 12px;
    fill: #666;
    transform: translate(15px, 0);
  }
</style>