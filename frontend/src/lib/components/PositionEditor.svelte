<!-- src/lib/components/PositionEditor.svelte -->

<script lang="ts">
  import { BoardPosition } from '$lib/types/BoardPosition';
  import { PieceType, type PieceClickEvent } from '$lib/types/Piece';
  import BoardPositionView from './PositionView/PositionView.svelte';

  // callback
  export let onChange: (newSfen?: string) => void;

  let sfen: string | undefined = undefined;

  const handleToggleTurn = () => {
    const position = new BoardPosition(sfen);
    position.isBlackTurn = !position.isBlackTurn;
    sfen = position.toSfen(1); // 手数は1固定
    onChange(sfen);
  };

  const handleRotatePiece = (row: number, col: number, isBlack: boolean) => {
    const position = new BoardPosition(sfen);

    // 先手の表駒 -> 先手の成駒 -> 後手の表駒 -> 後手の成駒
    if (isBlack) {
      const piece = position.blackBoard[row][col];
      if (piece & PieceType.PROMOTE || piece === PieceType.KI || piece === PieceType.OU) {
        // 先手の成駒or金or王 -> 後手の表駒
        position.whiteBoard[row][col] = piece & ~PieceType.PROMOTE;
        position.blackBoard[row][col] = PieceType.VACANCY;
      } else {
        // 先手の表駒 -> 先手の成駒
        position.blackBoard[row][col] = piece | PieceType.PROMOTE;
      }
    } else {
      const piece = position.whiteBoard[row][col];
      if (piece & PieceType.PROMOTE || piece === PieceType.KI || piece === PieceType.OU) {
        // 後手の成駒or金or王 -> 先手の表駒
        position.blackBoard[row][col] = piece & ~PieceType.PROMOTE;
        position.whiteBoard[row][col] = PieceType.VACANCY;
      } else {
        // 後手の表駒 -> 後手の成駒
        position.whiteBoard[row][col] = piece | PieceType.PROMOTE;
      }
    }

    sfen = position.toSfen(1); // 手数は1固定
    onChange(sfen);
  };

  const handleMovePiece = (moving: PieceClickEvent, target: PieceClickEvent) => {
    const position = new BoardPosition(sfen);

    // 移動先に駒がある場合、持ち駒へ
    if (target.source.type === 'board') {
      if (target.pieceType !== PieceType.VACANCY) {
        const originalType = target.pieceType & ~PieceType.PROMOTE;
        if (originalType === PieceType.OU) {
          // targetが玉の場合は、targetは駒箱へ
          const num = position.pieceBox.get(originalType) ?? 0;
          position.pieceBox.set(originalType, num + 1);
        } else if (moving.source.type === 'box' || moving.source.isBlack) {
          // 移動中の駒が、駒箱からor先手の場合、targetは先手の持ち駒へ
          const num = position.blackHands.get(originalType) ?? 0;
          position.blackHands.set(originalType, num + 1);
        } else {
          // 移動中の駒が後手の場合、targetは後手の持ち駒へ
          const num = position.whiteHands.get(originalType) ?? 0;
          position.whiteHands.set(originalType, num + 1);
        }
      }
    }

    // 移動先に駒を配置
    if (target.source.type === 'board') {
      if (moving.source.type === 'box' || moving.source.isBlack) {
        position.blackBoard[target.source.row][target.source.col] = moving.pieceType;
        position.whiteBoard[target.source.row][target.source.col] = PieceType.VACANCY;
      } else {
        position.blackBoard[target.source.row][target.source.col] = PieceType.VACANCY;
        position.whiteBoard[target.source.row][target.source.col] = moving.pieceType;
      }
    } else if (target.source.type === 'stand' || target.source.type === 'box') {
      const targetContainer =
        target.source.type === 'stand'
          ? target.source.isBlack
            ? position.blackHands
            : position.whiteHands
          : position.pieceBox;
      const originalType = moving.pieceType & ~PieceType.PROMOTE;
      const num = targetContainer.get(originalType) ?? 0;
      targetContainer.set(originalType, num + 1);
    }

    // 移動元の駒を削除
    if (moving.source.type === 'board') {
      position.blackBoard[moving.source.row][moving.source.col] = PieceType.VACANCY;
      position.whiteBoard[moving.source.row][moving.source.col] = PieceType.VACANCY;
    } else if (moving.source.type === 'stand' || moving.source.type === 'box') {
      const list =
        moving.source.type === 'stand'
          ? moving.source.isBlack
            ? position.blackHands
            : position.whiteHands
          : position.pieceBox;
      const num = list.get(moving.pieceType) || 0;
      if (num > 1) {
        list.set(moving.pieceType, num - 1);
      } else {
        list.delete(moving.pieceType);
      }
    }

    sfen = position.toSfen(1); // 手数は1固定
    onChange(sfen);
  };

  const handleReset = () => {
    sfen = undefined;
    onChange(sfen);
  };

  const handleAllInBox = () => {
    sfen = '9/9/9/9/9/9/9/9/9 b - 1';
    onChange(sfen);
  };
</script>

<div class="position-editor">
  <BoardPositionView
    mode="position"
    {sfen}
    onToggleTurn={handleToggleTurn}
    onRotatePiece={handleRotatePiece}
    onMovePiece={handleMovePiece}
  />
  <div class="controls">
    <button onclick={handleReset}>平手初形</button>
    <button onclick={handleAllInBox}>全て駒箱</button>
  </div>
</div>

<style>
  .controls {
    margin-top: 0.5rem;
    text-align: center;

    button {
      width: 10rem;
      padding: 0.3rem 0.5rem;
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
