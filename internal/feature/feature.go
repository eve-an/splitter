package feature

import (
	"errors"
	"slices"
)

const maximumWeight = 100

var ErrMaximumWeightExceeded = errors.New("maximum weight exceeded")
var ErrVariantAlreadyExist = errors.New("variant with the same name exist")

type Feature struct {
	Name         string
	Descritption string
	Active       bool
	Variants     Variants
}

func NewFeature(name string, description string, active bool, variants *Variants) (*Feature, error) {
	f := &Feature{
		Name:         name,
		Descritption: description,
		Active:       active,
		Variants:     *variants,
	}

	return f, f.Validate()
}

func (f *Feature) AddVariant(v *Variant) error {
	if f.Variants.TotalWeight()+uint32(v.Weight) > maximumWeight {
		return ErrMaximumWeightExceeded
	}

	if slices.Contains(f.Variants.Names(), v.Name) {
		return ErrVariantAlreadyExist
	}

	return nil
}

func (f *Feature) Validate() error {
	var errs []error
	if f.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	if f.Variants.TotalWeight() > maximumWeight {
		errs = append(errs, ErrMaximumWeightExceeded)
	}

	uniqueNames := make(map[string]struct{}, len(f.Variants))
	for _, name := range f.Variants.Names() {
		if _, found := uniqueNames[name]; !found {
			uniqueNames[name] = struct{}{}
			continue
		}

		errs = append(errs, ErrVariantAlreadyExist)
		break
	}

	return errors.Join(errs...)
}
