package cooking

import (
	"fmt"

	"github.com/tthung1997/buddy/core/cooking"
)

// SimpleUnitConverter converts between units of the same kind by way of the
// base unit of that kind. Volume to mass conversion needs an ingredient
// density, so it is reported as incompatible instead of being guessed.
type SimpleUnitConverter struct{}

func NewSimpleUnitConverter() *SimpleUnitConverter {
	return &SimpleUnitConverter{}
}

func (c *SimpleUnitConverter) Convert(amount float64, from cooking.Unit, to cooking.Unit) (float64, error) {
	if from == to {
		if !from.IsValid() {
			return 0, fmt.Errorf("%w: %q", cooking.ErrUnknownUnit, from.String())
		}
		return amount, nil
	}

	fromKind, err := from.Kind()
	if err != nil {
		return 0, err
	}
	toKind, err := to.Kind()
	if err != nil {
		return 0, err
	}
	if fromKind != toKind {
		return 0, fmt.Errorf("%w: cannot convert %s (%s) to %s (%s)",
			cooking.ErrIncompatibleUnits, from, fromKind, to, toKind)
	}

	fromFactor, err := from.BaseFactor()
	if err != nil {
		return 0, err
	}
	toFactor, err := to.BaseFactor()
	if err != nil {
		return 0, err
	}
	if toFactor == 0 {
		return 0, fmt.Errorf("%w: %q has no base factor", cooking.ErrUnknownUnit, to.String())
	}

	return amount * fromFactor / toFactor, nil
}

func (c *SimpleUnitConverter) Compatible(from cooking.Unit, to cooking.Unit) bool {
	fromKind, err := from.Kind()
	if err != nil {
		return false
	}
	toKind, err := to.Kind()
	if err != nil {
		return false
	}
	return fromKind == toKind
}
