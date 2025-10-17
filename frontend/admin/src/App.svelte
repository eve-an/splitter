<script lang="ts">
    import { onMount } from "svelte";
    import {
        createFeature,
        deleteFeature,
        listFeatures,
        updateFeature,
    } from "./lib/api/features";
    import {
        createEmptyFeatureForm,
        createEmptyVariant,
        type Feature,
        type FeatureFormState,
        type VariantForm,
    } from "./lib/types";
    import { validateFeatureForm } from "./lib/validation/feature";
    import Button from "./lib/components/Button.svelte";
    import Panel from "./lib/components/Panel.svelte";
    import Pill from "./lib/components/Pill.svelte";

    const ensureVariantRows = (value: FeatureFormState): FeatureFormState =>
        value.variants.length > 0
            ? value
            : {
                  ...value,
                  variants: [createEmptyVariant()],
              };

    const toFormState = (feature: Feature): FeatureFormState =>
        ensureVariantRows({
            id: feature.id,
            name: feature.name,
            description: feature.description ?? "",
            active: feature.active,
            variants: feature.variants.map((variant) => ({
                id: variant.id,
                name: variant.name,
                weight: variant.weight,
            })),
        });

    const toRequestPayload = (form: FeatureFormState) => ({
        name: form.name.trim(),
        description:
            form.description.trim() === "" ? null : form.description.trim(),
        active: form.active,
        variants: form.variants.map((variant) => ({
            name: variant.name.trim(),
            weight:
                typeof variant.weight === "number"
                    ? variant.weight
                    : Number.parseInt(String(variant.weight || "0"), 10) || 0,
        })),
    });

    const numericWeight = (value: VariantForm["weight"]) =>
        typeof value === "number"
            ? value
            : value === ""
              ? 0
              : Number.parseInt(String(value), 10) || 0;

    const listTileBaseClass =
        "flex w-full min-w-0 flex-col items-start gap-2 rounded-xl border border-transparent bg-slate-950/40 p-4 text-left transition hover:border-slate-700/80 hover:bg-slate-950/60 focus-visible:outline focus-visible:outline-2 focus-visible:outline-indigo-400/40 overflow-hidden";
    const listTileActiveClass =
        "border-indigo-400/80 bg-indigo-500/15 text-indigo-100 border-slate-700/80 shadow-lg shadow-indigo-500/25";
    const fieldBaseClass =
        "rounded-lg border border-slate-800 bg-slate-950/70 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 transition focus:outline-none focus-visible:border-indigo-400 focus-visible:ring-2 focus-visible:ring-indigo-500/40";
    const invalidFieldClass =
        "border-rose-500/70 focus-visible:border-rose-400 focus-visible:ring-rose-500/40";

    let features: Feature[] = [];
    let filtered: Feature[] = [];
    let selectedId: number | null = null;
    let search = "";

    let form: FeatureFormState = ensureVariantRows(createEmptyFeatureForm());
    let showErrors = false;

    let loading = false;
    let saving = false;
    let deleting = false;
    let loadMessage = "";
    let saveMessage = "";
    let successMessage = "";

    $: validation = validateFeatureForm(form);
    $: totalWeight = form.variants.reduce(
        (sum, variant) => sum + numericWeight(variant.weight),
        0,
    );
    $: filtered = (() => {
        const term = search.trim().toLowerCase();
        if (!term) return features;
        return features.filter((feature) => {
            const description = feature.description ?? "";
            return (
                feature.name.toLowerCase().includes(term) ||
                description.toLowerCase().includes(term)
            );
        });
    })();
    $: activeCount = filtered.filter((feature) => feature.active).length;
    $: inactiveCount = filtered.length - activeCount;

    const setForm = (next: FeatureFormState) => {
        form = ensureVariantRows({
            ...next,
            variants: next.variants.map((variant) => ({ ...variant })),
        });
        showErrors = false;
        saveMessage = "";
        successMessage = "";
    };

    const emptyForm = () => ensureVariantRows(createEmptyFeatureForm());

    const setFeatures = (items: Feature[]) => {
        features = [...items].sort((a, b) => a.name.localeCompare(b.name));
    };

    const currentFeature = () =>
        selectedId == null
            ? null
            : (features.find((feature) => feature.id === selectedId) ?? null);

    const selectFirst = () => {
        if (features.length === 0) {
            selectedId = null;
            setForm(emptyForm());
            return;
        }

        const first = features[0];
        selectedId = first.id;
        setForm(toFormState(first));
    };

    const loadFeatures = async () => {
        loading = true;
        loadMessage = "";
        try {
            const data = await listFeatures();
            console.log(data);
            setFeatures(data);

            if (selectedId != null) {
                const match = features.find(
                    (feature) => feature.id === selectedId,
                );
                if (match) {
                    setForm(toFormState(match));
                } else {
                    selectFirst();
                }
            } else {
                selectFirst();
            }
        } catch (error) {
            loadMessage =
                error instanceof Error
                    ? error.message
                    : "Could not load features. Please retry.";
        } finally {
            loading = false;
        }
    };

    onMount(loadFeatures);

    const startCreate = () => {
        selectedId = null;
        setForm(emptyForm());
    };

    const selectFeature = (id: number) => {
        const item = features.find((feature) => feature.id === id);
        if (!item) return;
        selectedId = id;
        setForm(toFormState(item));
    };

    const updateField = <K extends keyof FeatureFormState>(
        key: K,
        value: FeatureFormState[K],
    ) => {
        form = { ...form, [key]: value };
    };

    const updateVariant = (index: number, patch: Partial<VariantForm>) => {
        form = {
            ...form,
            variants: form.variants.map((variant, idx) =>
                idx === index ? { ...variant, ...patch } : variant,
            ),
        };
    };

    const addVariant = () => {
        form = {
            ...form,
            variants: [...form.variants, createEmptyVariant()],
        };
    };

    const removeVariant = (index: number) => {
        const remaining = form.variants.filter((_, idx) => idx !== index);
        form = {
            ...form,
            variants: remaining.length > 0 ? remaining : [createEmptyVariant()],
        };
    };

    const resetForm = () => {
        const existing = currentFeature();
        if (existing) {
            setForm(toFormState(existing));
            return;
        }
        startCreate();
    };

    const deleteCurrentFeature = async () => {
        const feature = currentFeature();
        if (!feature) return;

        const promptMessage = `Delete the feature “${feature.name}”? This action cannot be undone.`;
        if (typeof window !== "undefined" && !window.confirm(promptMessage)) {
            return;
        }

        deleting = true;
        saveMessage = "";
        successMessage = "";

        try {
            await deleteFeature(feature.id);
            const remaining = features.filter((item) => item.id !== feature.id);
            setFeatures(remaining);
            selectedId = null;
            selectFirst();
            successMessage = "Feature deleted.";
        } catch (error) {
            saveMessage =
                error instanceof Error
                    ? error.message
                    : "Delete failed. Please retry.";
        } finally {
            deleting = false;
        }
    };

    const upsertFeature = (item: Feature) => {
        const exists = features.some((feature) => feature.id === item.id);
        if (exists) {
            setFeatures(
                features.map((feature) =>
                    feature.id === item.id ? item : feature,
                ),
            );
        } else {
            setFeatures([...features, item]);
        }
    };

    const save = async () => {
        showErrors = true;

        if (!validation.isValid) {
            saveMessage = "Please fix the highlighted fields.";
            successMessage = "";
            return;
        }

        saving = true;
        saveMessage = "";
        successMessage = "";

        try {
            const payload = toRequestPayload(form);
            const isUpdate = selectedId != null;
            const result = isUpdate
                ? await updateFeature(selectedId!, payload)
                : await createFeature(payload);

            upsertFeature(result);
            selectedId = result.id;
            setForm(toFormState(result));

            await loadFeatures();

            successMessage = isUpdate ? "Feature updated." : "Feature created.";
        } catch (error) {
            saveMessage =
                error instanceof Error
                    ? error.message
                    : "Save failed. Please retry.";
        } finally {
            saving = false;
        }
    };
</script>

<main class="px-4 py-10 sm:px-6">
    <div class="mx-auto flex w-full max-w-6xl flex-col gap-8">
        <Panel
            tag="header"
            className="flex flex-col gap-4 p-6 sm:flex-row sm:items-center sm:justify-between"
        >
            <div>
                <h1 class="text-2xl font-semibold text-white">Feature flags</h1>
                <p class="text-sm text-slate-400">
                    Review existing flags and make quick edits.
                </p>
            </div>
            <Button variant="soft" on:click={loadFeatures} disabled={loading}>
                {loading ? "Loading…" : "Reload"}
            </Button>
        </Panel>

        {#if loadMessage}
            <p
                class="rounded-xl border border-rose-500/50 bg-rose-500/10 px-4 py-3 text-sm text-rose-200"
            >
                {loadMessage}
            </p>
        {/if}
        {#if successMessage}
            <p
                class="rounded-xl border border-emerald-500/40 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-100"
            >
                {successMessage}
            </p>
        {/if}
        {#if saveMessage}
            <p
                class="rounded-xl border border-rose-500/50 bg-rose-500/10 px-4 py-3 text-sm text-rose-200"
            >
                {saveMessage}
            </p>
        {/if}

        <section class="grid gap-6 lg:grid-cols-[320px_minmax(0,1fr)] xl:gap-8">
            <Panel
                tag="aside"
                className="flex h-full flex-col gap-5 p-5 lg:sticky lg:top-8 lg:self-start"
            >
                <div
                    class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between md:gap-3 lg:flex-col lg:items-stretch lg:gap-4"
                >
                    <input
                        class="flex-1 min-w-0 rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 transition focus:outline-none focus-visible:border-indigo-400 focus-visible:ring-2 focus-visible:ring-indigo-400/40"
                        type="search"
                        placeholder="Search features"
                        value={search}
                        on:input={(event) =>
                            (search = event.currentTarget.value)}
                    />
                    <Button
                        variant="primary"
                        className="w-full md:flex-none md:w-auto md:justify-center lg:w-full"
                        on:click={startCreate}
                        disabled={saving || deleting}
                    >
                        New feature
                    </Button>
                </div>

                <div class="flex flex-wrap gap-2 text-xs text-slate-400">
                    <Pill variant="muted">
                        Total {features.length}
                    </Pill>
                    <Pill variant="positive">
                        Active {activeCount}
                    </Pill>
                    <Pill variant="warning">
                        Inactive {inactiveCount}
                    </Pill>
                </div>

                <ul
                    class="flex min-h-0 flex-1 flex-col gap-2 overflow-y-auto overflow-x-hidden"
                >
                    {#if loading}
                        <li
                            class="w-full rounded-xl border border-transparent bg-slate-950/40 p-4 text-sm text-slate-400"
                        >
                            Loading…
                        </li>
                    {:else if filtered.length === 0}
                        <li
                            class="w-full rounded-xl border border-transparent bg-slate-950/40 p-4 text-sm text-slate-400"
                        >
                            {search.trim()
                                ? "No matches found."
                                : "No features yet. Create one to get started."}
                        </li>
                    {:else}
                        {#each filtered as feature, index (feature.id ?? `temp-${index}`)}
                            <li>
                                <button
                                    type="button"
                                    class={`${listTileBaseClass} ${feature.id === selectedId ? listTileActiveClass : ""}`.trim()}
                                    on:click={() => selectFeature(feature.id)}
                                >
                                    <div
                                        class="flex w-full min-w-0 items-center justify-between gap-3"
                                    >
                                        <span
                                            class="text-sm font-semibold text-white"
                                        >
                                            {feature.name}
                                        </span>
                                        <Pill
                                            variant={feature.active
                                                ? "positive"
                                                : "muted"}
                                        >
                                            {feature.active
                                                ? "Active"
                                                : "Inactive"}
                                        </Pill>
                                    </div>
                                    {#if feature.description}
                                        <p
                                            class="w-full truncate text-xs text-slate-400"
                                            title={feature.description}
                                        >
                                            {feature.description}
                                        </p>
                                    {/if}
                                    {#if feature.variants.length}
                                        {@const variantSummary =
                                            feature.variants
                                                .map(
                                                    (variant) =>
                                                        `${variant.name} (${variant.weight})`,
                                                )
                                                .join(", ")}
                                        <p
                                            class="w-full truncate text-[11px] text-slate-400"
                                            title={variantSummary}
                                        >
                                            {variantSummary}
                                        </p>
                                    {/if}
                                </button>
                            </li>
                        {/each}
                    {/if}
                </ul>
            </Panel>

            <Panel tag="section" className="flex flex-col gap-6 p-6 lg:p-7">
                <h2 class="text-lg font-semibold text-white">
                    {selectedId ? "Edit feature" : "Create feature"}
                </h2>

                <label class="flex flex-col gap-2 text-sm">
                    <span
                        class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400"
                    >
                        Name
                    </span>
                    <input
                        type="text"
                        class={`${fieldBaseClass} ${showErrors && validation.errors.name ? invalidFieldClass : ""}`.trim()}
                        value={form.name}
                        on:input={(event) =>
                            updateField("name", event.currentTarget.value)}
                        placeholder="checkout-button"
                    />
                    {#if showErrors && validation.errors.name}
                        <span class="text-xs text-rose-300"
                            >{validation.errors.name}</span
                        >
                    {/if}
                </label>

                <label class="flex items-center gap-2 text-sm text-slate-200">
                    <input
                        type="checkbox"
                        class="h-4 w-4 rounded border border-slate-700 bg-slate-950 text-indigo-400 transition focus-visible:ring-2 focus-visible:ring-indigo-400/40"
                        checked={form.active}
                        on:change={(event) =>
                            updateField("active", event.currentTarget.checked)}
                    />
                    <span>Active</span>
                </label>

                <label class="flex flex-col gap-2 text-sm">
                    <span
                        class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400"
                    >
                        Description
                    </span>
                    <textarea
                        class={`min-h-[120px] ${fieldBaseClass} ${showErrors && validation.errors.description ? invalidFieldClass : ""}`.trim()}
                        value={form.description}
                        on:input={(event) =>
                            updateField(
                                "description",
                                event.currentTarget.value,
                            )}
                    ></textarea>
                    {#if showErrors && validation.errors.description}
                        <span class="text-xs text-rose-300">
                            {validation.errors.description}
                        </span>
                    {/if}
                </label>

                <section
                    class="flex flex-col gap-4 rounded-xl border border-slate-800/70 bg-slate-950/40 p-4"
                >
                    <header
                        class="flex flex-col gap-2 md:flex-row md:items-center md:justify-between"
                    >
                        <div>
                            <h3 class="text-sm font-semibold text-white">
                                Variants
                            </h3>
                            <p class="text-xs text-slate-400">
                                Give each variant a unique name and send the
                                weights to 100 total.
                            </p>
                        </div>
                        <Button
                            variant="soft"
                            size="sm"
                            className="w-full md:w-auto"
                            on:click={addVariant}
                            disabled={saving || deleting}
                        >
                            Add variant
                        </Button>
                    </header>

                    {#if showErrors && validation.errors.variantsMessage}
                        <p class="text-xs text-rose-300">
                            {validation.errors.variantsMessage}
                        </p>
                    {/if}

                    {#each form.variants as variant, index}
                        <div
                            class="grid gap-4 rounded-xl border border-slate-800/60 bg-slate-950/30 p-4 md:grid-cols-[minmax(0,1fr)_minmax(120px,0.3fr)_auto] md:items-end"
                        >
                            <label class="flex flex-col gap-2 text-sm">
                                <span
                                    class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400"
                                >
                                    Variant name
                                </span>
                                <input
                                    type="text"
                                    class={`${fieldBaseClass} ${
                                        showErrors &&
                                        Boolean(
                                            validation.errors.variants[index]
                                                ?.name,
                                        )
                                            ? invalidFieldClass
                                            : ""
                                    }`.trim()}
                                    value={variant.name}
                                    on:input={(event) =>
                                        updateVariant(index, {
                                            name: event.currentTarget.value,
                                        })}
                                    placeholder={`Variant ${index + 1}`}
                                />
                                {#if showErrors && validation.errors.variants[index]?.name}
                                    <span class="text-xs text-rose-300">
                                        {validation.errors.variants[index]
                                            ?.name}
                                    </span>
                                {/if}
                            </label>

                            <label class="flex flex-col gap-2 text-sm">
                                <span
                                    class="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400"
                                >
                                    Weight
                                </span>
                                <input
                                    type="number"
                                    min="0"
                                    max="100"
                                    step="1"
                                    class={`${fieldBaseClass} ${
                                        showErrors &&
                                        Boolean(
                                            validation.errors.variants[index]
                                                ?.weight,
                                        )
                                            ? invalidFieldClass
                                            : ""
                                    }`.trim()}
                                    value={variant.weight}
                                    on:input={(event) => {
                                        const value = event.currentTarget.value;
                                        updateVariant(index, {
                                            weight:
                                                value === ""
                                                    ? ""
                                                    : Number.parseInt(
                                                          value,
                                                          10,
                                                      ),
                                        });
                                    }}
                                />
                                {#if showErrors && validation.errors.variants[index]?.weight}
                                    <span class="text-xs text-rose-300">
                                        {validation.errors.variants[index]
                                            ?.weight}
                                    </span>
                                {/if}
                            </label>

                            <Button
                                variant="danger"
                                size="sm"
                                className="justify-start px-0 w-full md:w-auto"
                                on:click={() => removeVariant(index)}
                                disabled={saving || deleting}
                            >
                                Delete variant
                            </Button>
                        </div>
                    {/each}

                    <footer
                        class="flex flex-wrap items-center justify-between gap-2 text-sm text-slate-400"
                    >
                        <span>Variant count: {form.variants.length}</span>
                        <span>
                            Total weight: {totalWeight}
                            {#if showErrors && validation.errors.totalWeight}
                                <span class="ml-1 text-xs text-rose-300">
                                    {validation.errors.totalWeight}
                                </span>
                            {/if}
                        </span>
                    </footer>
                </section>

                <footer
                    class="flex flex-col gap-3 border-t border-slate-800/70 pt-4 md:flex-row md:items-center md:justify-between"
                >
                    <div class="flex w-full md:w-auto">
                        {#if selectedId}
                            <Button
                                variant="danger"
                                className="w-full md:mr-auto md:w-auto"
                                on:click={deleteCurrentFeature}
                                disabled={saving || deleting}
                            >
                                {deleting ? "Deleting…" : "Delete feature"}
                            </Button>
                        {/if}
                    </div>
                    <div
                        class="flex flex-col gap-3 md:flex-row md:items-center md:justify-end md:gap-3 w-full md:w-auto"
                    >
                        <Button
                            variant="soft"
                            on:click={resetForm}
                            disabled={saving || deleting}
                            className="w-full md:w-auto"
                        >
                            Reset
                        </Button>
                        <Button
                            variant="primary"
                            on:click={save}
                            disabled={saving || deleting}
                            className="w-full md:w-auto"
                        >
                            {saving ? "Saving…" : "Save changes"}
                        </Button>
                    </div>
                </footer>
            </Panel>
        </section>
    </div>
</main>
