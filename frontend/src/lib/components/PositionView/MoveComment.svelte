<!-- src/lib/components/MoveComment.svelte -->

<script lang="ts">
  import type { ChangeEventHandler } from 'svelte/elements';

  export let x: number;
  export let y: number;
  export let comment: string;
  export let editable: boolean = true;
  export let closed: boolean = true;
  export let onCloseOpen: () => void = () => {};
  export let onChangeComment: (comment: string) => void = (comment) => {};
  const w = 3800; // 790 + 2220 + 790
  const h = 380; // 10 + 20 + 80*4 + 20 + 10

  const handleChangeComment: ChangeEventHandler<HTMLTextAreaElement> = (event: Event) => {
    comment = (event.target as HTMLTextAreaElement).value;
    onChangeComment(comment);
  };
</script>

<g>
  {#if closed}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <g
      transform={`translate(${x + 3300}, ${y})`}
      class="toggle-view-button"
      onclick={() => onCloseOpen && onCloseOpen()}
      role="button"
      tabindex="0"
    >
      <rect width="500" height="80" rx="20" ry="20" fill="#999" />
      <rect x="10" y="10" width="480" height="60" rx="10" ry="10" fill="#cfc" />
      <text
        x="250"
        y="40"
        dominant-baseline="middle"
        text-anchor="middle"
        fill="#666"
        font-size="60"
        style:font-size="60px">↓ show comment</text
      >
    </g>
  {:else}
    <foreignObject {x} {y} width={w} height={h}>
      <!-- xmlns属性がTypeScriptの型チェックで警告の対象となるため回避 -->
      <textarea
        {...{ xmlns: 'http://www.w3.org/1999/xhtml' } as any}
        class={'move-comment' + (editable ? ' editable' : '')}
        readonly={!editable}
        placeholder={editable ? '<comment>' : ''}
        onchange={handleChangeComment}>{comment}</textarea
      >
    </foreignObject>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <g
      transform={`translate(${x + 3720}, ${y})`}
      class="toggle-view-button"
      onclick={() => onCloseOpen && onCloseOpen()}
      role="button"
      tabindex="0"
    >
      <rect width="80" height="80" rx="20" ry="20" fill="#999" />
      <rect x="10" y="10" width="60" height="60" rx="10" ry="10" fill="#fff" />
      <text
        x="40"
        y="40"
        dominant-baseline="middle"
        text-anchor="middle"
        fill="#666"
        font-size="60"
        style:font-size="60px">x</text
      >
    </g>
  {/if}
</g>

<style lang="scss">
  .move-comment {
    width: 100%;
    height: 100%;
    background: #fff;
    border: 10px solid #ccc;
    border-radius: 20px;
    padding: 20px 50px;
    overflow-y: auto;
    line-height: 80px;
    font-size: 72px;
    color: #333;

    &::placeholder {
      color: #999;
    }

    &.editable {
      &:focus {
        border-color: #696;
        background-color: #ffe;
      }
    }
  }
  .toggle-view-button {
    cursor: pointer;
  }
  .toggle-view-button:hover rect:nth-child(1) {
    fill: #696;
  }
</style>
