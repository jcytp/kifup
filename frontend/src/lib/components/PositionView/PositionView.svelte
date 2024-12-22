<!-- src/lib/components/PositionView/PositionView.svelte -->

<script lang="ts">
  import { BoardPosition } from '$lib/types/BoardPosition';
  import type { KifuMove } from '$lib/types/Kifu';
  import { PieceType, type PieceClickEvent } from '$lib/types/Piece';
  import Board from './Board.svelte';
  import MoveList from './MoveList.svelte';
  import PieceBox from './PieceBox.svelte';
  import PieceStand from './PieceStand.svelte';
  import TurnIndicator from './TurnIndicator.svelte';

  // callbacks
  export let onChange: (position?: string, moves?: KifuMove[]) => void;

  // parameters
  export let mode: 'position' | 'moves' = 'position';
  export let sfen: string | undefined;
  export let moveList: KifuMove[] = [];
  export let moveNumber: number = 0;
  $: visibleMoveList = mode === 'moves';
  $: visiblePieceBox = mode === 'position';
  $: position = new BoardPosition(sfen);

  // $: sizeSet = {
  //   borderThin: 5,
  //   borderThic: 10,
  //   piecePadding: 60,
  //   boardMargin: 120,
  //   subBoxMargin: 60,
  //   cell: { w: 240, h: 260 },
  //   counter: { w: 120, h: 120 },
  //   board: {
  //     w: 2220, // 240 * 9 + 5 * 8 + 10 * 2
  //     h: 2400, // 260 * 9 + 5 * 8 + 10 * 2
  //     x: 980 + (visibleMoveList ? 860 : 0), // 60 + 800 + 120
  //     y: 120,
  //   },
  //   whiteStand: {
  //     w: 800, // (240 + 60) * 2 + 60 * 3 + 10 * 2
  //     h: 1360, // 260 * 4 + 60 * 5 + 10 * 2
  //     x: 60 + (visibleMoveList ? 860 : 0),
  //     y: 120,
  //   },
  //   blackStand: {
  //     w: 800,
  //     h: 1360,
  //     x: 3320 + (visibleMoveList ? 980 : 0), // 60 + 800 + 120 + 2220 + 120
  //     y: 1160, // 120 + 2400 - 1360
  //   },
  //   pieceBox: {
  //     w: 800, // (240 + 60) * 2 + 60 * 3 + 10 * 2
  //     h: 1680, // 260 * 5 + 60 * 6 + 10 * 2
  //     x: 4180 + (visibleMoveList ? 860 : 0), // 60 + 800 + 120 + 2220 + 120 + 800 + 60
  //     y: 840, // 120 + 2400 - 1680
  //   },
  //   moveList: {
  //     w: 800,
  //     h: 2400,
  //     x: 60,
  //     y: 120,
  //   },
  //   viewBox: {
  //     w: 4180 + (visibleMoveList ? 860 : 0) + (visiblePieceBox ? 860 : 0), // 2220 + 800 + 800 + 120 * 2 + 60 * 2 + ? + ?
  //     h: 2640, // 2400 + 120 * 2
  //   },
  // };

  $: viewBoxWidth = 4180 + (visibleMoveList ? 860 : 0) + (visiblePieceBox ? 860 : 0);
  const viewBoxHeight = 2640;

  const handleToggleTurn = () => {
    const newPosition = position.copy();
    newPosition.isBlackTurn = !newPosition.isBlackTurn;
    const newSfen = newPosition.toSfen(1); // 手数は1固定
    onChange(newSfen, moveList);
  };

  const handleRightClick = (event: PieceClickEvent) => {
    if (event.source.type !== 'board' || event.pieceType === PieceType.VACANCY) {
      return;
    }

    const i = event.source.row;
    const j = event.source.col;
    const newPosition = position.copy();
    if (event.source.isBlack) {
      const piece = newPosition.blackBoard[i][j];
      if (piece & PieceType.PROMOTE || piece === PieceType.KI || piece === PieceType.OU) {
        // 先手の成駒or金or王 -> 後手の表駒
        newPosition.whiteBoard[i][j] = piece & ~PieceType.PROMOTE;
        newPosition.blackBoard[i][j] = PieceType.VACANCY;
      } else {
        // 先手の表駒 -> 先手の成駒
        newPosition.blackBoard[i][j] = piece | PieceType.PROMOTE;
      }
    } else {
      const piece = newPosition.whiteBoard[i][j];
      if (piece & PieceType.PROMOTE || piece === PieceType.KI || piece === PieceType.OU) {
        // 後手の成駒or金or王 -> 先手の表駒
        newPosition.blackBoard[i][j] = piece & ~PieceType.PROMOTE;
        newPosition.whiteBoard[i][j] = PieceType.VACANCY;
      } else {
        // 後手の表駒 -> 後手の成駒
        newPosition.whiteBoard[i][j] = piece | PieceType.PROMOTE;
      }
    }

    // 新しいSFENを生成してコールバックを実行
    const newSfen = newPosition.toSfen(1); // 手数は1固定
    onChange(newSfen, moveList);
  };

  let pickedPiece: PieceClickEvent | undefined;
  const handlePieceClick = (newPickedPiece: PieceClickEvent) => {
    if (pickedPiece === undefined) {
      // 駒を持ちあげる
      if (newPickedPiece.pieceType !== PieceType.VACANCY) {
        pickedPiece = newPickedPiece;
      }
    } else {
      // 駒を移動させる
      const moving = pickedPiece; // 移動中の情報
      const target = newPickedPiece; // 移動先の情報
      pickedPiece = undefined;

      // 移動元の移動先が同じならキャンセル
      if (moving.source.type === 'board' && target.source.type === 'board') {
        if (moving.source.row === target.source.row && moving.source.col === target.source.col) {
          return;
        }
      }

      const newPosition = position.copy();

      // 移動先に駒がある場合、持ち駒へ
      if (target.source.type === 'board') {
        if (target.pieceType !== PieceType.VACANCY) {
          const originalType = target.pieceType & ~PieceType.PROMOTE;
          if (originalType === PieceType.OU) {
            // targetが玉の場合は、targetは駒箱へ
            const num = newPosition.pieceBox.get(originalType) ?? 0;
            newPosition.pieceBox.set(originalType, num + 1);
          } else if (moving.source.type === 'box' || moving.source.isBlack) {
            // 移動中の駒が、駒箱からor先手の場合、targetは先手の持ち駒へ
            const num = newPosition.blackHands.get(originalType) ?? 0;
            newPosition.blackHands.set(originalType, num + 1);
          } else {
            // 移動中の駒が後手の場合、targetは後手の持ち駒へ
            const num = newPosition.whiteHands.get(originalType) ?? 0;
            newPosition.whiteHands.set(originalType, num + 1);
          }
        }
      }

      // 移動先に駒を配置
      if (target.source.type === 'board') {
        if (moving.source.type === 'box' || moving.source.isBlack) {
          newPosition.blackBoard[target.source.row][target.source.col] = moving.pieceType;
          newPosition.whiteBoard[target.source.row][target.source.col] = PieceType.VACANCY;
        } else {
          newPosition.blackBoard[target.source.row][target.source.col] = PieceType.VACANCY;
          newPosition.whiteBoard[target.source.row][target.source.col] = moving.pieceType;
        }
      } else if (target.source.type === 'stand' || target.source.type === 'box') {
        const targetContainer =
          target.source.type === 'stand'
            ? target.source.isBlack
              ? newPosition.blackHands
              : newPosition.whiteHands
            : newPosition.pieceBox;
        const originalType = moving.pieceType & ~PieceType.PROMOTE;
        const num = targetContainer.get(originalType) ?? 0;
        targetContainer.set(originalType, num + 1);
      }

      // 移動元の駒を削除
      if (moving.source.type === 'board') {
        newPosition.blackBoard[moving.source.row][moving.source.col] = PieceType.VACANCY;
        newPosition.whiteBoard[moving.source.row][moving.source.col] = PieceType.VACANCY;
      } else if (moving.source.type === 'stand' || moving.source.type === 'box') {
        const list =
          moving.source.type === 'stand'
            ? moving.source.isBlack
              ? newPosition.blackHands
              : newPosition.whiteHands
            : newPosition.pieceBox;
        const num = list.get(moving.pieceType) || 0;
        if (num > 1) {
          list.set(moving.pieceType, num - 1);
        } else {
          list.delete(moving.pieceType);
        }
      }

      // 新しいSFENを生成してコールバックを実行
      const newSfen = newPosition.toSfen(1); // 手数は1固定
      onChange(newSfen, moveList);
    }
  };
  $: if (sfen) {
    pickedPiece = undefined;
    console.debug(sfen);
  }
</script>

<svg class="board-position-view-svg" viewBox={`0 0 ${viewBoxWidth} ${viewBoxHeight}`}>
  <rect width="100%" height="100%" fill="#fff" />

  {#if visibleMoveList}
    <MoveList x={60} y={120} {moveList} {moveNumber} />
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
      x={980}
      y={120}
      blackBoard={position.blackBoard}
      whiteBoard={position.whiteBoard}
      onCellClick={handlePieceClick}
      onRightClick={handleRightClick}
      {pickedPiece}
    />
    <PieceStand
      x={3320}
      y={1160}
      isBlack={true}
      hands={position.blackHands}
      onAreaClick={handlePieceClick}
      onPieceClick={handlePieceClick}
      {pickedPiece}
    />
    <TurnIndicator x={3320} y={120} isBlackTurn={position.isBlackTurn} onClick={handleToggleTurn} />
    {#if visiblePieceBox}
      <PieceBox
        x={4180}
        y={840}
        pieces={position.pieceBox}
        onAreaClick={handlePieceClick}
        onPieceClick={handlePieceClick}
        {pickedPiece}
      />
    {/if}
  </g>
</svg>
