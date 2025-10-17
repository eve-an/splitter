<script lang="ts">
    import type { HTMLAttributes } from "svelte/elements";

    type PillVariant = "muted" | "positive" | "warning";

    export let variant: PillVariant = "muted";
    export let className = "";
    const restProps = $$restProps as HTMLAttributes<HTMLElement>;

    const base =
        "inline-flex items-center gap-1 rounded-full border px-3 py-1 text-xs font-medium";
    const variantClasses: Record<PillVariant, string> = {
        muted: "border-slate-700/80 bg-slate-950/40 text-slate-300",
        positive: "border-emerald-500/50 bg-emerald-500/15 text-emerald-100",
        warning: "border-amber-500/50 bg-amber-500/15 text-amber-100",
    };

    $: classes = [base, variantClasses[variant], className]
        .filter(Boolean)
        .join(" ");
</script>

<span {...restProps} class={classes}>
    <slot />
</span>
