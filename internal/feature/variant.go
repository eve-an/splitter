package feature

import (
	"errors"
)

type Variant struct {
	ID     int64
	Name   string
	Weight uint8
}

type Variants []Variant

func (v Variants) TotalWeight() uint32 {
	totalSum := 0
	for _, vr := range v {
		totalSum += int(vr.Weight)
	}

	return uint32(totalSum)
}

func (v Variants) Names() []string {
	names := make([]string, 0, len(v))
	for _, vs := range v {
		names = append(names, vs.Name)
	}

	return names
}

func NewVariant(name string, weight uint8) (Variant, error) {
	if name == "" {
		return Variant{}, errors.New("name is required")
	}

	if weight > maximumWeight {
		return Variant{}, ErrMaximumWeightExceeded
	}

	return Variant{Name: name, Weight: weight}, nil
}

func NewVariants(vs ...Variant) (Variants, error) {
	variants := Variants(vs)
	if variants.TotalWeight() > maximumWeight {
		return nil, ErrMaximumWeightExceeded
	}

	uniqueNames := make(map[string]struct{}, len(variants))
	for _, variant := range variants {
		if _, exists := uniqueNames[variant.Name]; exists {
			return nil, ErrVariantAlreadyExist
		}
		uniqueNames[variant.Name] = struct{}{}
	}

	return variants, nil
}
