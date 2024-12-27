<!-- src/lib/components/MoveList.svelte -->

<script lang="ts">
  import type { KifuMove } from '$lib/types/Kifu';
  import { PiecePlace, PieceChar } from '$lib/types/Piece';

  export let isBlackFirst: boolean = true; //ToDo: set in upper component
  export let moveList: KifuMove[];
  export let moveNumber: number;
  export let x: number;
  export let y: number;
  const w = 800;
  const h = 2400;

  const testMoveList: KifuMove[] = [
    { number: 1, piece: 0x0, from_place: 0x67, to_place: 0x57 },
    { number: 2, piece: 0x0, from_place: 0x21, to_place: 0x31 },
  ];

  const formatMove = (move: KifuMove, num: number): string => {
    const turnMark = (isBlackFirst && num % 2 == 1) || (!isBlackFirst && num % 2 == 0) ? '▲' : '△';
    const toPlace = new PiecePlace(move.to_place);
    const toText = `${toPlace.fileChar()}${toPlace.rankChar()}`;
    const pieceChar = PieceChar.get(move.promote ? move.piece & ~0x8 : move.piece);
    const promoteText = move.promote ? '成' : move.promote === false ? '不成' : ''; // 3 variations: true, false, undefined
    const fromPlace = new PiecePlace(move.from_place);
    const fromText = fromPlace.isInHand() ? '打' : `(${fromPlace.fileNum()}${fromPlace.rankNum()})`;
    return `${turnMark}${toText}${pieceChar}${promoteText}${fromText}`;
  };
</script>

<foreignObject {x} {y} width={w} height={h}>
  <!-- xmlns属性がTypeScriptの型チェックで警告の対象となるため回避 -->
  <div {...{ xmlns: 'http://www.w3.org/1999/xhtml' } as any} class="moves-container">
    <div class="move-item" class:current={0 === moveNumber}>
      <span class="move-number">0</span>
      <span class="move-text">開始局面</span>
    </div>
    {#each moveList as move, i}
      <div class="move-item" class:current={i + 1 === moveNumber}>
        <span class="move-number">{i + 1}</span>
        <span class="move-text">{formatMove(move, i + 1)}</span>
      </div>
    {/each}
  </div>
</foreignObject>

<style>
  .moves-container {
    width: 100%;
    height: 100%;
    background: #fff;
    border: 10px solid #ccc;
    overflow-y: auto;

    .move-item {
      padding: 20px;
      display: flex;
      gap: 30px;
      align-items: baseline;
      font-size: 68px;
      line-height: 100px;
      cursor: pointer;

      &.current {
        background: #e0f0ff;
      }

      .move-number {
        width: 120px;
        text-align: right;
        color: #666;
      }

      &:hover {
        background-color: #ffe0f0;
      }
    }
  }
</style>
