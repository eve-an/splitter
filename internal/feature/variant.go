package feature

import (
	"errors"
)

type Variant struct {
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
