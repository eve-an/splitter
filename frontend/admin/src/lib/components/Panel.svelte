<script lang="ts">
    import type { HTMLAttributes } from "svelte/elements";

    export let tag: keyof HTMLElementTagNameMap = "div";
    export let className = "";
    const restProps = $$restProps as HTMLAttributes<HTMLElement>;
</script>

<svelte:element
    this={tag}
    {...restProps}
    class={`w-full rounded-2xl border border-slate-800/70 bg-slate-950/60 shadow-[0_20px_45px_-30px_rgba(12,16,45,0.9)] backdrop-blur ${className}`.trim()}
>
    <slot />
</svelte:element>
