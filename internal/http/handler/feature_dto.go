package handler

import "github.com/eve-an/splitter/internal/feature"

type variantPayload struct {
	Name   string `json:"name"`
	Weight uint8  `json:"weight"`
}

type featureRequest struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Active      bool             `json:"active"`
	Variants    []variantPayload `json:"variants"`
}

type eventRequest struct {
	UserID  string `json:"user_id"`
	Variant string `json:"variant"`
	Type    string `json:"type"`
}

type variantResponse struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Weight uint8  `json:"weight"`
}

type featureResponse struct {
	ID          int32             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Active      bool              `json:"active"`
	Variants    []variantResponse `json:"variants"`
}

func mapFeatureResponse(feature *feature.Feature) featureResponse {
	return featureResponse{
		ID:          feature.ID,
		Name:        feature.Name,
		Description: feature.Descritption,
		Active:      feature.Active,
		Variants:    mapVariantsResponse(feature.Variants),
	}
}

func mapVariantsResponse(variants []feature.Variant) []variantResponse {
	variantResponses := make([]variantResponse, len(variants))
	for i, variant := range variants {
		variantResponses[i] = variantResponse{
			ID:     variant.ID,
			Name:   variant.Name,
			Weight: variant.Weight,
		}
	}
	return variantResponses
}

func buildFeatureFromRequest(req *featureRequest) (*feature.Feature, error) {
	variants := make([]feature.Variant, 0, len(req.Variants))
	for _, v := range req.Variants {
		variant, err := feature.NewVariant(v.Name, v.Weight)
		if err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}

	domainVariants, err := feature.NewVariants(variants...)
	if err != nil {
		return nil, err
	}

	return feature.NewFeature(req.Name, req.Description, req.Active, &domainVariants)
}
