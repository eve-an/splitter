export interface paths {
  "/api/v1/features": {
    get: operations["listFeatures"];
    post: operations["createFeature"];
  };
  "/api/v1/features/{featureID}": {
    get: operations["getFeature"];
    put: operations["updateFeature"];
    parameters: {
      path: {
        featureID: components["parameters"]["FeatureId"];
      };
    };
  };
}

export interface operations {
  listFeatures: {
    responses: {
      200: {
        content: {
          "application/json": components["schemas"]["Feature"][];
        };
      };
    };
  };
  getFeature: {
    parameters: {
      path: {
        featureID: components["parameters"]["FeatureId"];
      };
    };
    responses: {
      200: {
        content: {
          "application/json": components["schemas"]["Feature"];
        };
      };
    };
  };
  createFeature: {
    requestBody: {
      content: {
        "application/json": components["schemas"]["FeatureRequest"];
      };
    };
    responses: {
      201: {
        content: {
          "application/json": components["schemas"]["Feature"];
        };
      };
    };
  };
  updateFeature: {
    parameters: {
      path: {
        featureID: components["parameters"]["FeatureId"];
      };
    };
    requestBody: {
      content: {
        "application/json": components["schemas"]["FeatureRequest"];
      };
    };
    responses: {
      200: {
        content: {
          "application/json": components["schemas"]["Feature"];
        };
      };
    };
  };
}

export interface components {
  parameters: {
    FeatureId: number;
  };
  schemas: {
    Feature: {
      id: number;
      name: string;
      description?: string | null;
      active: boolean;
      variants: components["schemas"]["Variant"][];
    };
    Variant: {
      id?: number;
      name: string;
      weight: number;
    };
    FeatureRequest: {
      name: string;
      description?: string | null;
      active: boolean;
      variants?: components["schemas"]["VariantRequest"][];
    };
    VariantRequest: {
      name: string;
      weight: number;
    };
  };
}

export type webhooks = Record<string, never>;
