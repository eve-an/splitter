package feature

import (
	"errors"
	"fmt"
)

type Feature struct {
	Name         string
	Descritption string
	Active       bool
	Variants     []Variant
}

func (f Feature) Validate() error {
	var errs []error
	if f.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}

	totalWeight := 0
	for _, v := range f.Variants {
		totalWeight += int(v.Weight)
	}

	if totalWeight > 100 {
		errs = append(errs, fmt.Errorf("total weight exceeds 100: %d", totalWeight))
	}

	return errors.Join(errs...)
}
