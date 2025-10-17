import type { Feature, FeatureRequestPayload, Variant } from "../types";
import { apiBaseUrl, apiClient } from "./client";
import type { components } from "./swagger/schema";

type ApiFeature = components["schemas"]["Feature"];
type ApiVariant = components["schemas"]["Variant"];
type ApiFeatureRequest = components["schemas"]["FeatureRequest"];

const mapVariant = (variant: ApiVariant): Variant => ({
  id: variant.id,
  name: variant.name ?? "",
  weight: Number.isFinite(variant.weight) ? variant.weight : 0,
});

const mapFeature = (feature: ApiFeature): Feature => {
  return {
    id: feature.id,
    name: feature.name ?? "",
    description: feature.description ?? null,
    active: Boolean(feature.active),
    variants: Array.isArray(feature.variants)
      ? feature.variants.map(mapVariant)
      : [],
  };
};

const toRequestBody = (payload: FeatureRequestPayload): ApiFeatureRequest => ({
  name: payload.name,
  description: payload.description,
  active: payload.active,
  variants: payload.variants.map((variant) => ({
    name: variant.name,
    weight: variant.weight,
  })),
});

const errorMessage = (error: unknown, fallback: string) => {
  if (error instanceof Error && error.message.trim().length > 0) {
    return error.message;
  }
  if (
    error &&
    typeof error === "object" &&
    "message" in error &&
    typeof (error as { message?: unknown }).message === "string"
  ) {
    const message = (error as { message: string }).message.trim();
    if (message) return message;
  }
  return fallback;
};

export const listFeatures = async (): Promise<Feature[]> => {
  const response = await apiClient.GET("/api/v1/features");
  if (response.error) {
    throw new Error(errorMessage(response.error, "Failed to load features."));
  }

  if (response.response.status !== 200 || !response.data) {
    throw new Error(
      `Failed to load features. Server responded with ${response.response.status}.`,
    );
  }
  return response.data.map(mapFeature);
};

export const createFeature = async (
  payload: FeatureRequestPayload,
): Promise<Feature> => {
  const response = await apiClient.POST("/api/v1/features", {
    body: toRequestBody(payload),
  });

  if (response.error) {
    throw new Error(errorMessage(response.error, "Failed to create feature."));
  }
  if (response.response.status !== 201 || !response.data) {
    throw new Error(
      `Failed to create feature. Server responded with ${response.response.status}.`,
    );
  }

  return mapFeature(response.data);
};

export const updateFeature = async (
  featureId: number,
  payload: FeatureRequestPayload,
): Promise<Feature> => {
  const response = await apiClient.PUT("/api/v1/features/{featureID}", {
    params: {
      path: {
        featureID: featureId,
      },
    },
    body: toRequestBody(payload),
  });

  if (response.error) {
    throw new Error(errorMessage(response.error, "Failed to update feature."));
  }

  if (response.response.status !== 200 || !response.data) {
    throw new Error(
      `Failed to update feature. Server responded with ${response.response.status}.`,
    );
  }

  return mapFeature(response.data);
};

export const deleteFeature = async (featureId: number): Promise<void> => {
  const response = await fetch(`${apiBaseUrl}/api/v1/features/${featureId}`, {
    method: "DELETE",
    credentials: "same-origin",
  });

  if (!response.ok) {
    throw new Error(
      `Failed to delete feature. Server responded with ${response.status}.`,
    );
  }
};
