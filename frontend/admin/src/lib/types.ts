export interface Variant {
  id?: number;
  name: string;
  weight: number;
}

export interface Feature {
  id: number;
  name: string;
  description: string | null;
  active: boolean;
  variants: Variant[];
}

export interface VariantForm {
  id?: number;
  name: string;
  weight: number | "";
}

export interface FeatureFormState {
  id?: number;
  name: string;
  description: string;
  active: boolean;
  variants: VariantForm[];
}

export const createEmptyVariant = (): VariantForm => ({
  name: "",
  weight: "",
});

export const createEmptyFeatureForm = (): FeatureFormState => ({
  name: "",
  description: "",
  active: true,
  variants: [createEmptyVariant()],
});

export interface FeatureRequestPayload {
  name: string;
  description: string | null;
  active: boolean;
  variants: Array<{
    name: string;
    weight: number;
  }>;
}
