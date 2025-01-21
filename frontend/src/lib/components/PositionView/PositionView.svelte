<!-- src/lib/components/PositionView/PositionView.svelte -->

<!-- 棋譜作成時の局面編集、棋譜更新での指し手の編集、棋譜詳細での棋譜再生、各ページで共通 -->

<script lang="ts">
  import { BoardPosition } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import { PieceMovables, PiecePlace, PieceType, type PieceClickEvent } from '$lib/types/Piece';
  import Board from './Board.svelte';
  import MoveList from './MoveList.svelte';
  import PieceBox from './PieceBox.svelte';
  import PieceStand from './PieceStand.svelte';
  import PromoteIndicator from './PromoteIndicator.svelte';
  import TurnIndicator from './TurnIndicator.svelte';

  // callbacks
  export let onToggleTurn: () => void = () => {};
  export let onRotatePiece: (row: number, col: number, isBlack: boolean) => void = (
    row,
    col,
    isBlack
  ) => {};
  export let onMovePiece: (moving: PieceClickEvent, target: PieceClickEvent) => void = (
    moving,
    target
  ) => {};
  export let onAppendMove: (num: number, move: KifuMove) => void = (num, move) => {};
  export let onPromote: (promote: boolean) => void = (promote) => {};

  // parameters
  export let mode: 'position' | 'moves' | 'replay' = 'position';
  export let showNumbers: boolean = true;
  export let isBlackFirst: boolean | undefined = undefined;
  export let moveList: KifuMove[] = [];
  export let moveNumber: number = 0; // 現在の手数
  export let sfen: string | undefined; // 現在の局面
  $: visibleMoveList = mode === 'moves' || mode === 'replay';
  $: visiblePieceBox = mode === 'position';
  $: position = new BoardPosition(sfen);
  let promoteChoice: PieceType = PieceType.VACANCY;

  $: viewBoxWidth = 3920 + (visibleMoveList ? 860 : 0) + (visiblePieceBox ? 860 : 0);
  const viewBoxHeight = 2640;

  const handleToggleTurn = () => {
    if (mode === 'position') {
      onToggleTurn();
    }
  };

  const handleRightClick = (event: PieceClickEvent) => {
    if (mode === 'position') {
      if (event.source.type !== 'board' || event.pieceType === PieceType.VACANCY) {
        return;
      }
      const row = event.source.row;
      const col = event.source.col;
      const isBlack = event.source.isBlack || false;
      onRotatePiece(row, col, isBlack);
    }
  };

  let pickedPiece: PieceClickEvent | undefined;
  const handlePieceClick = (newPickedPiece: PieceClickEvent) => {
    if (mode === 'replay') {
      return;
    }

    if (pickedPiece === undefined) {
      if (newPickedPiece.pieceType !== PieceType.VACANCY) {
        if (mode === 'moves') {
          if (promoteChoice != PieceType.VACANCY) {
            return;
          }

          // movesモードでは、自分の駒以外は持ちあげられない
          if (
            newPickedPiece.source.type === 'box' ||
            newPickedPiece.source.isBlack != position.isBlackTurn
          ) {
            return;
          }
        }

        // 駒を持ちあげる
        pickedPiece = newPickedPiece;
      }
    } else {
      const moving = pickedPiece; // 移動中の情報
      const target = newPickedPiece; // 移動先の情報
      pickedPiece = undefined; // 駒の持ちあげをリセット

      // 移動元と移動先が同じなら何もしない
      if (moving.source.type === 'board' && target.source.type === 'board') {
        if (moving.source.row === target.source.row && moving.source.col === target.source.col) {
          return;
        }
      }
      if (moving.source.type === 'stand' && target.source.type === 'stand') {
        if (moving.source.isBlack === target.source.isBlack) {
          return;
        }
      }
      if (moving.source.type === 'box' && target.source.type === 'box') {
        return;
      }

      // 駒を移動させる
      if (mode === 'position') {
        onMovePiece(moving, target);
      } else if (mode === 'moves' && isLegalMove(moving, target)) {
        if (isPossibleToPromote(moving, target)) {
          promoteChoice = moving.pieceType;
        }
        const move: KifuMove = {
          number: moveNumber + 1,
          piece: moving.pieceType,
          from_place: getPlaceOfEvent(moving),
          to_place: getPlaceOfEvent(target),
          catch_piece: target.pieceType == PieceType.VACANCY ? undefined : target.pieceType,
        };
        onAppendMove(moveNumber + 1, move);
      }
    }
  };

  const handleSelectPoromote = (promote: boolean) => {
    promoteChoice = PieceType.VACANCY;
    onPromote(promote);
  };

  // 駒の動きの合法性を確認する
  const isLegalMove = (moving: PieceClickEvent, target: PieceClickEvent): boolean => {
    console.debug('isLegalMove');
    // 駒台or盤上から盤上へのみが許可される
    if (moving.source.type === 'box' || target.source.type !== 'board') {
      return false;
    }
    // 行き先に自分の駒があってはいけない
    if (target.source.isBlack === position.isBlackTurn) {
      return false;
    }

    if (moving.source.type === 'stand') {
      // 駒を打つときは、targetに駒があってはいけない
      if (target.pieceType !== PieceType.VACANCY) {
        return false;
      }
    } else {
      // 盤上の移動では、from-toの位置関係が合法か
      let isRegalLocation = false;
      const x = target.source.col - moving.source.col;
      const y = target.source.row - moving.source.row;
      const probables = PieceMovables.get(moving.pieceType) || [];
      for (const probable of probables) {
        if (probable.x === x && (position.isBlackTurn ? probable.y : -probable.y) === y) {
          isRegalLocation = true;
          break;
        }
      }
      if (!isRegalLocation) {
        return false;
      }

      // 間に駒があってはいけない
      if (moving.pieceType !== PieceType.KE && (x > 1 || x < -1 || y > 1 || y < -1)) {
        const dx = x === 0 ? 0 : x > 0 ? 1 : -1;
        const dy = y === 0 ? 0 : y > 0 ? 1 : -1;
        let mx = moving.source.col + dx;
        let my = moving.source.row + dy;
        for (let i = 0; i < 8; i++) {
          if (mx < 0 || mx > 8 || my < 0 || my > 8) {
            break;
          }
          if (mx === target.source.col && my === target.source.row) {
            break;
          }
          if (
            position.blackBoard[my][mx] !== PieceType.VACANCY ||
            position.whiteBoard[my][mx] !== PieceType.VACANCY
          ) {
            return false;
          }
          mx += dx;
          my += dy;
        }
      }
    }

    console.debug('isLegalMove: legal');
    return true;
  };

  const isPossibleToPromote = (moving: PieceClickEvent, target: PieceClickEvent): boolean => {
    if (moving.source.type !== 'board' || target.source.type !== 'board') {
      return false;
    }

    const piece = moving.pieceType;
    if (piece & PieceType.PROMOTE || piece == PieceType.KI || piece == PieceType.OU) {
      return false;
    }

    if (position.isBlackTurn && (moving.source.row < 3 || target.source.row < 3)) {
      return true;
    }
    if (!position.isBlackTurn && (moving.source.row > 5 || target.source.row > 5)) {
      return true;
    }
    return false;
  };

  const getPlaceOfEvent = (event: PieceClickEvent): number => {
    const place = new PiecePlace();
    if (event.source.type === 'board') {
      place.setRowCol(event.source.row, event.source.col);
    }
    return place.val;
  };

  $: if (sfen) {
    pickedPiece = undefined;
    console.debug(sfen);
  }
</script>

<svg class="board-position-view-svg" viewBox={`0 0 ${viewBoxWidth} ${viewBoxHeight}`}>
  <rect width="100%" height="100%" fill="#fff" />

  {#if visibleMoveList}
    <MoveList x={60} y={120} {isBlackFirst} {moveList} {moveNumber} />
  {/if}
  <g transform={`translate(${visibleMoveList ? 860 : 0}, 0)`}>
    <PieceStand
      x={60}
      y={120}
      isBlack={false}
      hands={position.whiteHands}
      onAreaClick={handlePieceClick}
      onPieceClick={handlePieceClick}
      {pickedPiece}
    />
    <Board
      x={850}
      y={120}
      blackBoard={position.blackBoard}
      whiteBoard={position.whiteBoard}
      onCellClick={handlePieceClick}
      onRightClick={handleRightClick}
      {pickedPiece}
    />
    {#if showNumbers}
      {#each '１２３４５６７８９' as file, i}
        <text
          x={2940 - i * 245}
          y={85}
          dominant-baseline="middle"
          text-anchor="middle"
          fill="#666"
          font-size={70}
          style:font-size="70px"
        >
          {file}
        </text>
      {/each}
      {#each '一二三四五六七八九' as file, i}
        <text
          x={3105}
          y={260 + i * 265}
          dominant-baseline="middle"
          text-anchor="middle"
          fill="#666"
          font-size={70}
          style:font-size="70px"
        >
          {file}
        </text>
      {/each}
    {/if}
    <!-- y = 120 + 10 + 260*5 + 5*5 - 10 -->
    <PieceStand
      x={3190}
      y={1445}
      isBlack={true}
      hands={position.blackHands}
      onAreaClick={handlePieceClick}
      onPieceClick={handlePieceClick}
      {pickedPiece}
    />
    <TurnIndicator x={3190} y={120} isBlackTurn={position.isBlackTurn} onClick={handleToggleTurn} />
    {#if visiblePieceBox}
      <PieceBox
        x={3920}
        y={840}
        pieces={position.pieceBox}
        onAreaClick={handlePieceClick}
        onPieceClick={handlePieceClick}
        {pickedPiece}
      />
    {/if}
    <!-- 850 + (2220 / 2) - (900 / 2) -->
    {#if promoteChoice !== PieceType.VACANCY}
      <PromoteIndicator
        x={1310}
        y={1010}
        pieceType={promoteChoice}
        onSelect={handleSelectPoromote}
      />
    {/if}
  </g>
</svg>
