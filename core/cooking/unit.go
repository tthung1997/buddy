package cooking

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnknownUnit       = errors.New("unknown unit")
	ErrIncompatibleUnits = errors.New("incompatible units")
)

type Unit string

const (
	Teaspoon   Unit = "tsp"
	Tablespoon Unit = "tbsp"
	Cup        Unit = "cup"
	Milliliter Unit = "ml"
	Liter      Unit = "l"
	FluidOunce Unit = "fl oz"
	Milligram  Unit = "mg"
	Gram       Unit = "g"
	Kilogram   Unit = "kg"
	Ounce      Unit = "oz"
	Pound      Unit = "lbs"
	Piece      Unit = "piece"
)

type UnitKind string

const (
	UnitKindVolume UnitKind = "VOLUME"
	UnitKindMass   UnitKind = "MASS"
	UnitKindCount  UnitKind = "COUNT"
)

// unitDefinition expresses a unit through the base unit of its kind:
// milliliters for volume, grams for mass, pieces for count.
type unitDefinition struct {
	Kind       UnitKind
	InBaseUnit float64
}

var unitDefinitions = map[Unit]unitDefinition{
	Teaspoon:   {Kind: UnitKindVolume, InBaseUnit: 4.92892159375},
	Tablespoon: {Kind: UnitKindVolume, InBaseUnit: 14.78676478125},
	Cup:        {Kind: UnitKindVolume, InBaseUnit: 236.5882365},
	Milliliter: {Kind: UnitKindVolume, InBaseUnit: 1},
	Liter:      {Kind: UnitKindVolume, InBaseUnit: 1000},
	FluidOunce: {Kind: UnitKindVolume, InBaseUnit: 29.5735295625},
	Milligram:  {Kind: UnitKindMass, InBaseUnit: 0.001},
	Gram:       {Kind: UnitKindMass, InBaseUnit: 1},
	Kilogram:   {Kind: UnitKindMass, InBaseUnit: 1000},
	Ounce:      {Kind: UnitKindMass, InBaseUnit: 28.349523125},
	Pound:      {Kind: UnitKindMass, InBaseUnit: 453.59237},
	Piece:      {Kind: UnitKindCount, InBaseUnit: 1},
}

var orderedUnits = []Unit{
	Teaspoon, Tablespoon, Cup, Milliliter, Liter, FluidOunce,
	Milligram, Gram, Kilogram, Ounce, Pound,
	Piece,
}

// Units returns every supported unit in a stable, display friendly order.
func Units() []Unit {
	units := make([]Unit, len(orderedUnits))
	copy(units, orderedUnits)
	return units
}

// UnitsOfKind returns the supported units belonging to the given kind.
func UnitsOfKind(kind UnitKind) []Unit {
	units := []Unit{}
	for _, unit := range orderedUnits {
		if unitDefinitions[unit].Kind == kind {
			units = append(units, unit)
		}
	}
	return units
}

func (u Unit) IsValid() bool {
	_, exists := unitDefinitions[u]
	return exists
}

func (u Unit) Kind() (UnitKind, error) {
	definition, exists := unitDefinitions[u]
	if !exists {
		return "", fmt.Errorf("%w: %q", ErrUnknownUnit, string(u))
	}
	return definition.Kind, nil
}

// BaseFactor returns how many base units of this unit's kind one unit is worth:
// milliliters for volume, grams for mass, pieces for count.
func (u Unit) BaseFactor() (float64, error) {
	definition, exists := unitDefinitions[u]
	if !exists {
		return 0, fmt.Errorf("%w: %q", ErrUnknownUnit, string(u))
	}
	return definition.InBaseUnit, nil
}

func (u Unit) String() string {
	return string(u)
}

// ParseUnit resolves user supplied text to a supported unit, tolerating case
// and the common spelling variants of the longer units.
func ParseUnit(value string) (Unit, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "_", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")
	switch normalized {
	case "floz", "fluid ounce", "fluid ounces":
		normalized = string(FluidOunce)
	case "lb", "pound", "pounds":
		normalized = string(Pound)
	case "pieces", "pc", "pcs", "unit", "units":
		normalized = string(Piece)
	}

	unit := Unit(normalized)
	if !unit.IsValid() {
		return "", fmt.Errorf("%w: %q", ErrUnknownUnit, value)
	}
	return unit, nil
}

type IUnitConverter interface {
	// Convert restates an amount from one unit into another unit of the same
	// kind. Converting across kinds (volume to mass) returns an error.
	Convert(amount float64, from Unit, to Unit) (float64, error)
	// Compatible reports whether Convert can succeed for the pair.
	Compatible(from Unit, to Unit) bool
}
