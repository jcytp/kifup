<!-- src/lib/components/MoveList.svelte -->

<script lang="ts">
  import type { KifuMove } from '$lib/types/Kifu';
  import {
    fileChar,
    fileNum,
    PIECE_PLACE_IN_HAND,
    PieceChar,
    rankChar,
    rankNum,
  } from '$lib/types/Piece';

  export let moveList: KifuMove[];
  export let moveNumber: number;
  export let x: number;
  export let y: number;
  const w = 800;
  const h = 2400;

  $: if (moveList) {
    console.debug(moveList);
  }
  const testMoveList: KifuMove[] = [
    { number: 1, piece: 0x0, from_place: 0x67, to_place: 0x57 },
    { number: 2, piece: 0x0, from_place: 0x21, to_place: 0x31 },
  ];

  const formatMove = (move: KifuMove): string => {
    const pieceChar = PieceChar.get(move.promote ? move.piece & ~0x8 : move.piece);
    const toFileChar = fileChar(move.to_place);
    const toRankChar = rankChar(move.to_place);
    const promoteText = move.promote ? '成' : move.promote === false ? '不成' : ''; // 3 variations: true, false, undefined
    const fromText = ((place: number): string => {
      if (place === PIECE_PLACE_IN_HAND) return '打';
      const fromFile = fileNum(place);
      const fromRank = rankNum(place);
      return `(${fromFile}${fromRank})`;
    })(move.from_place);
    return `${toFileChar}${toRankChar}${pieceChar}${promoteText}${fromText}`;
  };
</script>

<foreignObject {x} {y} width={w} height={h}>
  <!-- xmlns属性がTypeScriptの型チェックで警告の対象となるため回避 -->
  <div {...{ xmlns: 'http://www.w3.org/1999/xhtml' } as any} class="moves-container">
    <div class="move-item" class:current={0 === moveNumber}>
      <span class="move-number">0</span>
      <span class="move-text">開始局面</span>
    </div>
    {#each testMoveList as move, i}
      <div class="move-item" class:current={i + 1 === moveNumber}>
        <span class="move-number">{i + 1}</span>
        <span class="move-text">{formatMove(move)}</span>
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
      padding: 20px 40px;
      display: flex;
      gap: 40px;
      align-items: baseline;
      font-size: 80px;
      line-height: 120px;

      &.current {
        background: #e0f0ff;
      }

      .move-number {
        width: 140px;
        text-align: right;
        color: #666;
      }

      .move-text {
        font-weight: bold;
      }
    }
  }
</style>
