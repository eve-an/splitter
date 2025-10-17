<script lang="ts">
    import type { HTMLButtonAttributes } from "svelte/elements";

    type ButtonVariant = "primary" | "soft" | "ghost" | "danger";
    type ButtonSize = "md" | "sm";

    export let type: "button" | "submit" | "reset" = "button";
    export let variant: ButtonVariant = "soft";
    export let size: ButtonSize = "md";
    export let disabled = false;
    export let className = "";
    const restProps = $$restProps as HTMLButtonAttributes;

    const base =
        "inline-flex items-center justify-center gap-2 rounded-lg font-semibold transition focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 disabled:cursor-not-allowed disabled:opacity-60";
    const sizeClasses: Record<ButtonSize, string> = {
        md: "px-4 py-2 text-sm",
        sm: "px-3 py-1.5 text-xs",
    };
    const variantClasses: Record<ButtonVariant, string> = {
        primary:
            "border border-indigo-500/70 bg-indigo-500 text-white shadow-lg shadow-indigo-500/20 hover:bg-indigo-400 focus-visible:outline-indigo-400",
        soft: "border border-slate-700/80 bg-slate-900/70 text-slate-200 hover:border-indigo-400 hover:text-indigo-200 focus-visible:outline-slate-400/40",
        ghost: "text-rose-200 hover:text-rose-100 focus-visible:outline-rose-300/40",
        danger:
            "border border-rose-500/70 bg-rose-500/10 text-rose-200 hover:bg-rose-500/20 focus-visible:outline-rose-400",
    };

    $: classes = [
        base,
        sizeClasses[size],
        variantClasses[variant],
        className,
    ]
        .filter(Boolean)
        .join(" ");
</script>

<button
    {...restProps}
    type={type}
    class={classes}
    {disabled}
    on:click
>
    <slot />
</button>
