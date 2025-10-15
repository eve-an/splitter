package feature

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"slices"
	"strconv"
)

const maximumWeight = 100

var (
	ErrMaximumWeightExceeded = errors.New("maximum weight exceeded")
	ErrVariantAlreadyExist   = errors.New("variant with the same name exist")
)

type Feature struct {
	ID           int64
	Name         string
	Descritption string
	Active       bool
	Variants     Variants
}

func NewFeature(
	name string,
	description string,
	active bool,
	variants *Variants,
) (*Feature, error) {
	var variantList Variants
	if variants != nil {
		variantList = *variants
	}

	f := &Feature{
		Name:         name,
		Descritption: description,
		Active:       active,
		Variants:     variantList,
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

	f.Variants = append(f.Variants, *v)

	return nil
}

func (f *Feature) Validate() error {
	var errs []error // nolint: prealloc
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

func featureHashForUser(u *User, f *Feature) uint64 {
	buf := make([]byte, 8+len(f.Name))

	binary.LittleEndian.PutUint64(buf, uint64(u.ID()))
	copy(buf[8:], []byte(f.Name))

	sum := md5.Sum(buf)
	return binary.LittleEndian.Uint64(sum[:8])
}

// VariantForUser determines the deterministically computed variant assignment for a user
func VariantForUser(u *User, feature *Feature) *Variant {
	bucket := uint8(featureHashForUser(u, feature) % 100)

	var cumulative uint8
	for i, v := range feature.Variants {
		if v.Weight == 0 {
			continue
		}

		cumulative += v.Weight

		if bucket < cumulative || i == len(feature.Variants)-1 {
			return &v
		}
	}

	panic("variants dont add up to " + strconv.FormatInt(maximumWeight, 10))
}
