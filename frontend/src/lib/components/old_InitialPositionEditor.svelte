<!-- src/lib/components/InitialPositionEditor.svelte -->

<script lang="ts">
  import type { BoardPosition, PieceType, CellPosition } from '$lib/types/Kifu';

  // コールバックプロパティ
  export let change: (position: BoardPosition) => void;

  // 将棋の駒の文字表示用マッピング
  const pieceKanji: { [K in PieceType]: string } = {
    歩: '歩',
    香: '香',
    桂: '桂',
    銀: '銀',
    金: '金',
    角: '角',
    飛: '飛',
    玉: '玉',
    と: 'と',
    成香: '杏',
    成桂: '圭',
    成銀: '全',
    馬: '馬',
    龍: '龍',
  };

  // 成ることができる駒のマッピング（元の駒 → 成り駒）
  const promotionMap: { [K in PieceType]?: PieceType } = {
    歩: 'と',
    香: '成香',
    桂: '成桂',
    銀: '成銀',
    角: '馬',
    飛: '龍',
  };

  // 成り駒から元の駒へのマッピング（成り駒 → 元の駒）
  const unpromoteMap: { [K in PieceType]?: PieceType } = {
    と: '歩',
    成香: '香',
    成桂: '桂',
    成銀: '銀',
    馬: '角',
    龍: '飛',
  };

  // 駒の循環を管理する関数（先手→先手成→後手→後手成→先手）
  function cyclepiece(piece: PieceType, isBlack: boolean): [PieceType, boolean] {
    // 金・玉は成れないので、先手⇔後手のみ
    if (piece === '金' || piece === '玉') {
      return [piece, !isBlack];
    }

    // 現在の駒が成り駒かどうかをチェック
    const unpromoted = unpromoteMap[piece];
    if (unpromoted) {
      // 成り駒の場合は、反対の手番の元の駒に変化
      return [unpromoted, !isBlack];
    }

    // 現在の駒が元の駒の場合
    const promoted = promotionMap[piece];
    if (promoted) {
      // 同じ手番の成り駒に変化
      return [promoted, isBlack];
    }

    // ここには来ないはずだが、型安全のため
    return [piece, !isBlack];
  }

  // 局面のディープコピーを行う関数
  function deepCopyPosition(pos: BoardPosition): BoardPosition {
    return {
      pieces: pos.pieces.map((piece) => ({
        position: { ...piece.position },
        piece: piece.piece,
        isBlack: piece.isBlack,
      })),
      hands: {
        black: pos.hands ? { ...pos.hands.black } : {},
        white: pos.hands ? { ...pos.hands.white } : {},
      },
    };
  }

  // サイズ設定
  const boardSize = 400;
  const cellSize = boardSize / 9;
  const komadaiWidth = cellSize * 2;
  const pieceboxWidth = cellSize * 2;

  // 初期配置（平手）
  const initialPosition: BoardPosition = {
    pieces: [
      // 後手の駒
      { position: { x: 1, y: 1 }, piece: '香', isBlack: false },
      { position: { x: 2, y: 1 }, piece: '桂', isBlack: false },
      { position: { x: 3, y: 1 }, piece: '銀', isBlack: false },
      { position: { x: 4, y: 1 }, piece: '金', isBlack: false },
      { position: { x: 5, y: 1 }, piece: '玉', isBlack: false },
      { position: { x: 6, y: 1 }, piece: '金', isBlack: false },
      { position: { x: 7, y: 1 }, piece: '銀', isBlack: false },
      { position: { x: 8, y: 1 }, piece: '桂', isBlack: false },
      { position: { x: 9, y: 1 }, piece: '香', isBlack: false },
      { position: { x: 2, y: 2 }, piece: '飛', isBlack: false },
      { position: { x: 8, y: 2 }, piece: '角', isBlack: false },
      ...Array(9)
        .fill(null)
        .map((_, i) => ({
          position: { x: i + 1, y: 3 },
          piece: '歩' as PieceType,
          isBlack: false,
        })),
      // 先手の駒
      { position: { x: 1, y: 9 }, piece: '香', isBlack: true },
      { position: { x: 2, y: 9 }, piece: '桂', isBlack: true },
      { position: { x: 3, y: 9 }, piece: '銀', isBlack: true },
      { position: { x: 4, y: 9 }, piece: '金', isBlack: true },
      { position: { x: 5, y: 9 }, piece: '玉', isBlack: true },
      { position: { x: 6, y: 9 }, piece: '金', isBlack: true },
      { position: { x: 7, y: 9 }, piece: '銀', isBlack: true },
      { position: { x: 8, y: 9 }, piece: '桂', isBlack: true },
      { position: { x: 9, y: 9 }, piece: '香', isBlack: true },
      { position: { x: 2, y: 8 }, piece: '角', isBlack: true },
      { position: { x: 8, y: 8 }, piece: '飛', isBlack: true },
      ...Array(9)
        .fill(null)
        .map((_, i) => ({
          position: { x: i + 1, y: 7 },
          piece: '歩' as PieceType,
          isBlack: true,
        })),
    ],
    hands: {
      black: {},
      white: {},
    },
  };

  // 全駒の定義（各駒の枚数）
  const ALL_PIECES: { [K in PieceType]: number } = {
    歩: 18,
    香: 4,
    桂: 4,
    銀: 4,
    金: 4,
    角: 2,
    飛: 2,
    玉: 2,
    と: 0,
    成香: 0,
    成桂: 0,
    成銀: 0,
    馬: 0,
    龍: 0,
  };

  // 使用されていない駒を計算する関数
  function calculateUnusedPieces(): { [K in PieceType]?: number } {
    const unused = { ...ALL_PIECES };

    // 盤上の駒を集計
    position.pieces.forEach((p) => {
      const basicPiece = unpromoteMap[p.piece] || p.piece;
      unused[basicPiece]--;
    });

    // 駒台の駒を集計
    if (position.hands) {
      Object.entries(position.hands.black).forEach(([piece, count]) => {
        unused[piece as PieceType] -= count || 0;
      });
      Object.entries(position.hands.white).forEach(([piece, count]) => {
        unused[piece as PieceType] -= count || 0;
      });
    }

    // 0枚以下の駒を除外
    return Object.fromEntries(Object.entries(unused).filter(([_, count]) => count > 0)) as {
      [K in PieceType]?: number;
    };
  }

  // 駒箱からの駒の取り出し
  function takeFromPiecebox(piece: PieceType) {
    const unusedPieces = calculateUnusedPieces();
    if (unusedPieces[piece]) {
      startDragFromKomadai(piece, true); // デフォルトで先手の駒として取り出す
    }
  }

  // 現在の局面状態
  let position = deepCopyPosition(initialPosition);

  // ドラッグ中の駒の状態
  let draggedPiece: {
    piece: PieceType;
    isBlack: boolean;
    fromHand?: boolean;
  } | null = null;

  // マウスの位置
  let mouseX = 0;
  let mouseY = 0;

  // マウス位置の更新
  function handleMouseMove(event: MouseEvent) {
    const rect = (event.currentTarget as SVGElement).getBoundingClientRect();
    mouseX = event.clientX - rect.left;
    mouseY = event.clientY - rect.top;
  }

  // 局面変更時の処理
  function updatePosition() {
    position = { ...position };
    change(position); // dispatchの代わりにコールバックを直接呼び出し
  }

  // 駒のダブルクリック処理
  function handlePieceDoubleClick(index: number) {
    draggedPiece = null;
    const piece = position.pieces[index];
    const [newPiece, newIsBlack] = cyclepiece(piece.piece, piece.isBlack);

    position.pieces[index] = {
      ...piece,
      piece: newPiece,
      isBlack: newIsBlack,
    };

    updatePosition();
  }

  // 駒のドラッグ開始
  function startDragFromBoard(index: number) {
    const piece = position.pieces[index];
    position.pieces.splice(index, 1);
    draggedPiece = { piece: piece.piece, isBlack: piece.isBlack, fromHand: false };
  }
  function startDragFromKomadai(piece: PieceType, isBlack: boolean) {
    draggedPiece = { piece, isBlack, fromHand: true };
  }

  // 駒の配置
  function placePiece(x: number, y: number) {
    if (!draggedPiece) return;

    // 既存の駒を取る
    const existingPieceIndex = position.pieces.findIndex(
      (p) => p.position.x === x && p.position.y === y
    );
    console.debug(existingPieceIndex);

    if (existingPieceIndex !== -1) {
      const capturedPiece = position.pieces[existingPieceIndex];
      position.pieces.splice(existingPieceIndex, 1);

      // 取った駒を駒台へ
      if (!position.hands) {
        position.hands = { black: {}, white: {} };
      }
      const hand = draggedPiece.isBlack ? position.hands.black : position.hands.white;
      const basicPiece = capturedPiece.piece.replace(/と|成|馬|龍/g, '') as PieceType;
      hand[basicPiece] = (hand[basicPiece] || 0) + 1;
    }

    // 新しい駒を配置
    if (draggedPiece.fromHand) {
      if (!position.hands) {
        position.hands = { black: {}, white: {} };
      }
      const hand = draggedPiece.isBlack ? position.hands.black : position.hands.white;
      hand[draggedPiece.piece] = Math.max(0, (hand[draggedPiece.piece] || 0) - 1);
    }

    position.pieces.push({
      position: { x, y },
      piece: draggedPiece.piece,
      isBlack: draggedPiece.isBlack,
    });

    draggedPiece = null;
    updatePosition();
  }

  // 持ち駒の表示用に集計
  function countPieces(hand: { [K in PieceType]?: number } = {}) {
    return Object.entries(hand).filter(([_, count]) => count > 0);
  }

  // リセットボタン
  function reset() {
    console.debug(initialPosition);
    position = deepCopyPosition(initialPosition);
    updatePosition();
  }
</script>

<div class="position-editor">
  <svg
    viewBox={`0 0 ${boardSize + komadaiWidth * 2 + pieceboxWidth} ${boardSize}`}
    on:mousemove={handleMouseMove}
    role="application"
    aria-label="将棋盤エディタ"
  >
    <!-- 後手の駒台 -->
    <g transform={`translate(0, 0)`}>
      <rect x="0" y="0" width={komadaiWidth} height={boardSize} class="komadai" />
      {#if position.hands}
        {#each countPieces(position.hands.white) as [piece, count], i}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <g
            transform={`translate(${komadaiWidth / 2}, ${cellSize * (1 + i)})`}
            on:click={() => startDragFromKomadai(piece as PieceType, false)}
            role="button"
            aria-label={`後手の${pieceKanji[piece as PieceType]}${count > 1 ? ` ${count}枚` : ''}`}
            tabindex="0"
          >
            <text class="piece-text white">{pieceKanji[piece as PieceType]}</text>
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
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <rect
            x={cellSize * i}
            y={cellSize * j}
            width={cellSize}
            height={cellSize}
            class="cell"
            on:click={() => placePiece(i + 1, j + 1)}
            role="button"
            aria-label={`${9 - i}${j + 1}のマス`}
            tabindex="0"
          />
        {/each}
      {/each}

      <!-- 駒の配置 -->
      {#each position.pieces as piece, i}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <g
          transform={`translate(
            ${cellSize * (piece.position.x - 0.5)},
            ${cellSize * (piece.position.y - 0.5)}
          )`}
          on:dblclick={() => handlePieceDoubleClick(i)}
          on:click={() => startDragFromBoard(i)}
          role="button"
          aria-label={`${piece.isBlack ? '先手' : '後手'}の${pieceKanji[piece.piece]}`}
          tabindex="0"
        >
          <text class="piece-text" class:black={piece.isBlack} class:white={!piece.isBlack}>
            {pieceKanji[piece.piece]}
          </text>
        </g>
      {/each}
    </g>

    <!-- 先手の駒台 -->
    <g transform={`translate(${boardSize + komadaiWidth}, 0)`}>
      <rect x="0" y="0" width={komadaiWidth} height={boardSize} class="komadai" />
      {#if position.hands}
        {#each countPieces(position.hands.black) as [piece, count], i}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <g
            transform={`translate(${komadaiWidth / 2}, ${cellSize * (1 + i)})`}
            on:click={() => startDragFromKomadai(piece as PieceType, true)}
            role="button"
            aria-label={`先手の${pieceKanji[piece as PieceType]}${count > 1 ? ` ${count}枚` : ''}`}
            tabindex="0"
          >
            <text class="piece-text black">{pieceKanji[piece as PieceType]}</text>
            {#if count > 1}
              <text class="count-text">x{count}</text>
            {/if}
          </g>
        {/each}
      {/if}
    </g>

    <!-- 駒箱（先手の駒台の右側） -->
    <g transform={`translate(${boardSize + komadaiWidth * 2}, 0)`}>
      <rect x="0" y="0" width={pieceboxWidth} height={boardSize} class="piecebox" />
      {#each Object.entries(calculateUnusedPieces()) as [piece, count], i}
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <g
          transform={`translate(${pieceboxWidth / 2}, ${cellSize * (1 + i)})`}
          on:click={() => takeFromPiecebox(piece as PieceType)}
          role="button"
          aria-label={`駒箱の${pieceKanji[piece as PieceType]}${count > 1 ? ` ${count}枚` : ''}`}
          tabindex="0"
        >
          <text class="piece-text unused">{pieceKanji[piece as PieceType]}</text>
          {#if count > 1}
            <text class="count-text">x{count}</text>
          {/if}
        </g>
      {/each}
    </g>

    <!-- ドラッグ中の駒 -->
    <!-- {#if draggedPiece}
      <g transform={`translate(${mouseX}, ${mouseY})`}>
        <text
          class="piece-text dragging"
          class:black={draggedPiece.isBlack}
          class:white={!draggedPiece.isBlack}
        >
          {pieceKanji[draggedPiece.piece]}
        </text>
      </g>
    {/if} -->
  </svg>

  <div class="controls">
    <button on:click={reset}>初期局面に戻す</button>
  </div>
</div>

<style lang="scss">
  .position-editor {
    width: 100%;
    max-width: 800px;
    margin: 0 auto;
  }

  svg {
    width: 100%;
    height: auto;
    user-select: none;
  }

  .cell {
    fill: #fff;
    stroke: #000;
    stroke-width: 1;

    &:hover {
      fill: #f0f0f0;
      cursor: pointer;
    }
  }

  .komadai {
    fill: #f0f0f0;
    stroke: #000;
    stroke-width: 1;
  }

  .piecebox {
    fill: #e8e8e8; // 駒台よりも少し暗めの色
    stroke: #000;
    stroke-width: 1;
  }

  .piece-text {
    // 既存のスタイルに追加
    &.unused {
      fill: #444; // 未使用の駒は少し薄めの色
    }
  }

  .piece-text {
    text-anchor: middle;
    font-size: 24px;
    font-weight: bold;
    cursor: pointer;

    &.black {
      fill: #000;
    }

    &.white {
      fill: #666;
      transform: rotate(180deg);
    }

    // &.dragging {
    //   pointer-events: none;
    // }
  }

  .count-text {
    text-anchor: middle;
    font-size: 12px;
    fill: #666;
    transform: translate(15px, 0);
  }

  .controls {
    margin-top: 1rem;
    text-align: center;

    button {
      padding: 0.5rem 1rem;
      background-color: var(--primary-color);
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;

      &:hover {
        opacity: 0.9;
      }
    }
  }
</style>
