import type { FeatureFormState, VariantForm } from "../types";

export type VariantValidationErrors = Partial<{
  name: string;
  weight: string;
}>;

export interface FeatureValidationErrors {
  name?: string;
  description?: string;
  totalWeight?: string;
  variantsMessage?: string;
  variants: VariantValidationErrors[];
}

export interface FeatureValidationResult {
  errors: FeatureValidationErrors;
  isValid: boolean;
}

const MAX_DESCRIPTION_LENGTH = 500;

const parseWeight = (
  value: VariantForm["weight"],
): { weight: number | null; empty: boolean } => {
  if (typeof value === "number") {
    return {
      weight: Number.isFinite(value) ? value : null,
      empty: false,
    };
  }

  if (value === "") {
    return { weight: null, empty: true };
  }

  const parsed = Number.parseInt(String(value).trim(), 10);

  if (Number.isNaN(parsed)) {
    return { weight: null, empty: false };
  }

  return { weight: parsed, empty: false };
};

export const validateFeatureForm = (
  form: FeatureFormState,
): FeatureValidationResult => {
  const errors: FeatureValidationErrors = {
    variants: form.variants.map(() => ({})),
  };

  let isValid = true;
  const trimmedName = form.name.trim();

  if (!trimmedName) {
    errors.name = "Feature name is required.";
    isValid = false;
  } else if (trimmedName.length < 3) {
    errors.name = "Use at least 3 characters.";
    isValid = false;
  }

  const description = form.description.trim();
  if (description.length > MAX_DESCRIPTION_LENGTH) {
    errors.description = `Keep the description under ${MAX_DESCRIPTION_LENGTH} characters (${description.length} used).`;
    isValid = false;
  }

  if (form.variants.length === 0) {
    errors.variantsMessage = "Add at least one rollout variant.";
    isValid = false;
  }

  const duplicateNames = new Map<string, number[]>();
  let totalWeight = 0;

  form.variants.forEach((variant, index) => {
    const entry: VariantValidationErrors = {};
    const name = variant.name.trim();

    if (!name) {
      entry.name = "Name is required.";
      isValid = false;
    } else {
      const key = name.toLowerCase();
      duplicateNames.set(key, [...(duplicateNames.get(key) ?? []), index]);
    }

    const { weight, empty } = parseWeight(variant.weight);

    if (weight === null) {
      entry.weight = empty ? "Weight is required." : "Weight must be a number.";
      isValid = false;
    } else if (weight < 0 || weight > 100) {
      entry.weight = "Weight must be between 0 and 100.";
      isValid = false;
    } else if (!Number.isInteger(weight)) {
      entry.weight = "Use whole numbers for weights.";
      isValid = false;
    } else {
      totalWeight += weight;
    }

    errors.variants[index] = entry;
  });

  duplicateNames.forEach((indexes) => {
    if (indexes.length < 2) return;
    indexes.forEach((idx) => {
      errors.variants[idx] = {
        ...errors.variants[idx],
        name: "Each variant needs a unique name.",
      };
    });
    isValid = false;
  });

  const hasWeightError = errors.variants.some((variant) => Boolean(variant.weight));

  if (!hasWeightError && form.variants.length > 0 && totalWeight !== 100) {
    errors.totalWeight =
      totalWeight > 0
        ? `Total weight must equal 100 (currently ${totalWeight}).`
        : "Distribute 100 total weight across variants.";
    isValid = false;
  }

  return { errors, isValid };
};
