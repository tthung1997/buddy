package cooking

import (
	"errors"
	"math"
	"testing"

	"github.com/tthung1997/buddy/core/cooking"
)

const tolerance = 1e-9

func assertClose(t *testing.T, got float64, want float64, context string) {
	t.Helper()
	if math.Abs(got-want) > tolerance*math.Max(1, math.Abs(want)) {
		t.Errorf("%s: expected %v, got %v", context, want, got)
	}
}

func TestUnit_Kind(t *testing.T) {
	cases := map[cooking.Unit]cooking.UnitKind{
		cooking.Teaspoon:   cooking.UnitKindVolume,
		cooking.Cup:        cooking.UnitKindVolume,
		cooking.Liter:      cooking.UnitKindVolume,
		cooking.FluidOunce: cooking.UnitKindVolume,
		cooking.Milligram:  cooking.UnitKindMass,
		cooking.Kilogram:   cooking.UnitKindMass,
		cooking.Pound:      cooking.UnitKindMass,
		cooking.Piece:      cooking.UnitKindCount,
	}

	for unit, expected := range cases {
		kind, err := unit.Kind()
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", unit, err)
		}
		if kind != expected {
			t.Errorf("expected %s to be %s, got %s", unit, expected, kind)
		}
	}
}

func TestUnit_Kind_UnknownUnit(t *testing.T) {
	if _, err := cooking.Unit("furlong").Kind(); !errors.Is(err, cooking.ErrUnknownUnit) {
		t.Fatalf("expected ErrUnknownUnit, got %v", err)
	}
	if _, err := cooking.Unit("furlong").BaseFactor(); !errors.Is(err, cooking.ErrUnknownUnit) {
		t.Fatalf("expected ErrUnknownUnit, got %v", err)
	}
	if cooking.Unit("furlong").IsValid() {
		t.Error("expected an unknown unit to be invalid")
	}
}

func TestUnit_BaseFactor(t *testing.T) {
	factor, err := cooking.Kilogram.BaseFactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertClose(t, factor, 1000, "kilogram base factor")

	factor, err = cooking.Milliliter.BaseFactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertClose(t, factor, 1, "milliliter base factor")
}

func TestParseUnit(t *testing.T) {
	cases := map[string]cooking.Unit{
		"tsp":         cooking.Teaspoon,
		"  TBSP ":     cooking.Tablespoon,
		"Cup":         cooking.Cup,
		"fl oz":       cooking.FluidOunce,
		"floz":        cooking.FluidOunce,
		"FL_OZ":       cooking.FluidOunce,
		"fluid ounce": cooking.FluidOunce,
		"lb":          cooking.Pound,
		"pounds":      cooking.Pound,
		"pieces":      cooking.Piece,
		"unit":        cooking.Piece,
	}

	for input, expected := range cases {
		unit, err := cooking.ParseUnit(input)
		if err != nil {
			t.Fatalf("unexpected error parsing %q: %v", input, err)
		}
		if unit != expected {
			t.Errorf("expected %q to parse as %s, got %s", input, expected, unit)
		}
	}
}

func TestParseUnit_Unknown(t *testing.T) {
	if _, err := cooking.ParseUnit("smidge"); !errors.Is(err, cooking.ErrUnknownUnit) {
		t.Fatalf("expected ErrUnknownUnit, got %v", err)
	}
}

func TestUnits_AreAllValidAndOrdered(t *testing.T) {
	units := cooking.Units()
	if len(units) == 0 {
		t.Fatal("expected at least one unit")
	}
	for _, unit := range units {
		if !unit.IsValid() {
			t.Errorf("expected %s to be valid", unit)
		}
	}

	// Mutating the returned slice must not affect later calls.
	units[0] = cooking.Unit("mutated")
	if cooking.Units()[0] == cooking.Unit("mutated") {
		t.Error("expected Units to return a defensive copy")
	}
}

func TestUnitsOfKind(t *testing.T) {
	count := cooking.UnitsOfKind(cooking.UnitKindCount)
	if len(count) != 1 || count[0] != cooking.Piece {
		t.Fatalf("expected count units to be [piece], got %v", count)
	}

	mass := cooking.UnitsOfKind(cooking.UnitKindMass)
	for _, unit := range mass {
		kind, err := unit.Kind()
		if err != nil || kind != cooking.UnitKindMass {
			t.Errorf("expected %s to be a mass unit", unit)
		}
	}
	if len(mass) != 5 {
		t.Errorf("expected 5 mass units, got %d", len(mass))
	}
}
