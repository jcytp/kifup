<!-- src/lib/components/PositionView/Board.svelte -->

<script lang="ts">
  import { PieceType, type PieceClickEvent } from '$lib/types/Piece';
  import Piece from './Piece.svelte';

  export let blackBoard: PieceType[][];
  export let whiteBoard: PieceType[][];
  export let x: number;
  export let y: number;
  export let style: 'default' | 'pentagon' | undefined = 'default';
  export let onCellClick: ((e: PieceClickEvent) => void) | undefined = undefined;
  export let onRightClick: ((e: PieceClickEvent) => void) | undefined = undefined;
  export let pickedPiece: PieceClickEvent | undefined = undefined;

  const isPickedPiece = (i: number, j: number): boolean => {
    if (!pickedPiece) return false;
    if (pickedPiece.source.type !== 'board') return false;
    if (pickedPiece.source.row !== i || pickedPiece.source.col !== j) return false;
    return true;
  };
  const pieceState = (i: number, j: number): 'normal' | 'picked' => {
    if (isPickedPiece(i, j)) {
      return 'picked';
    }
    return 'normal';
  };

  const buildPieceClickEvent = (i: number, j: number): PieceClickEvent => {
    let pieceType = PieceType.VACANCY;
    let isBlack: boolean | undefined = undefined;
    if (blackBoard[i][j] !== PieceType.VACANCY) {
      pieceType = blackBoard[i][j];
      isBlack = true;
    } else if (whiteBoard[i][j] !== PieceType.VACANCY) {
      pieceType = whiteBoard[i][j];
      isBlack = false;
    }
    const event = {
      pieceType: pieceType,
      source: {
        type: 'board',
        row: i,
        col: j,
        isBlack: isBlack,
      },
    } as PieceClickEvent;
    return event;
  };
</script>

<g transform={`translate(${x}, ${y})`}>
  <rect width={2220} height={2400} fill="#666" />

  {#each blackBoard as row, i}
    {#each row as blackPiece, j}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <rect
        x={10 + j * 245}
        y={10 + i * 265}
        width="240"
        height="260"
        fill="#ffc"
        onclick={onCellClick ? () => onCellClick(buildPieceClickEvent(i, j)) : undefined}
        oncontextmenu={onRightClick
          ? (e) => {
              e.preventDefault();
              onRightClick(buildPieceClickEvent(i, j));
            }
          : undefined}
        role="button"
        aria-label={`board-${i}-${j}`}
        tabindex="0"
      />
      {#if blackPiece !== PieceType.VACANCY}
        <g transform={`translate(${10 + j * 245}, ${10 + i * 265})`}>
          <Piece
            pieceType={blackPiece}
            {style}
            reverse={false}
            useViewBox={false}
            state={pieceState(i, j)}
          />
        </g>
      {/if}
      {#if whiteBoard[i][j] !== PieceType.VACANCY}
        <g transform={`translate(${10 + j * 245}, ${10 + i * 265})`}>
          <Piece
            pieceType={whiteBoard[i][j]}
            {style}
            reverse={true}
            useViewBox={false}
            state={pieceState(i, j)}
          />
        </g>
      {/if}
    {/each}
  {/each}
</g>
