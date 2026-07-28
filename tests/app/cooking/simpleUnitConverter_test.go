package cooking

import (
	"errors"
	"math"
	"testing"

	appCooking "github.com/tthung1997/buddy/app/cooking"
	coreCooking "github.com/tthung1997/buddy/core/cooking"
)

var _ coreCooking.IUnitConverter = (*appCooking.SimpleUnitConverter)(nil)

func assertClose(t *testing.T, got float64, want float64, context string) {
	t.Helper()
	if math.Abs(got-want) > 1e-6*math.Max(1, math.Abs(want)) {
		t.Errorf("%s: expected %v, got %v", context, want, got)
	}
}

func TestSimpleUnitConverter_Convert_SameKind(t *testing.T) {
	converter := appCooking.NewSimpleUnitConverter()

	cases := []struct {
		amount   float64
		from     coreCooking.Unit
		to       coreCooking.Unit
		expected float64
	}{
		{1, coreCooking.Kilogram, coreCooking.Gram, 1000},
		{2500, coreCooking.Gram, coreCooking.Kilogram, 2.5},
		{1, coreCooking.Liter, coreCooking.Milliliter, 1000},
		{3, coreCooking.Teaspoon, coreCooking.Tablespoon, 1},
		{1, coreCooking.Cup, coreCooking.Milliliter, 236.5882365},
		{1, coreCooking.Pound, coreCooking.Ounce, 16},
		{4, coreCooking.Piece, coreCooking.Piece, 4},
	}

	for _, testCase := range cases {
		got, err := converter.Convert(testCase.amount, testCase.from, testCase.to)
		if err != nil {
			t.Fatalf("unexpected error converting %v %s to %s: %v", testCase.amount, testCase.from, testCase.to, err)
		}
		assertClose(t, got, testCase.expected, string(testCase.from)+" to "+string(testCase.to))
	}
}

func TestSimpleUnitConverter_Convert_RoundTrip(t *testing.T) {
	converter := appCooking.NewSimpleUnitConverter()

	converted, err := converter.Convert(750, coreCooking.Milliliter, coreCooking.FluidOunce)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	back, err := converter.Convert(converted, coreCooking.FluidOunce, coreCooking.Milliliter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertClose(t, back, 750, "round trip")
}

func TestSimpleUnitConverter_Convert_IncompatibleKinds(t *testing.T) {
	converter := appCooking.NewSimpleUnitConverter()

	if _, err := converter.Convert(1, coreCooking.Cup, coreCooking.Gram); !errors.Is(err, coreCooking.ErrIncompatibleUnits) {
		t.Fatalf("expected ErrIncompatibleUnits, got %v", err)
	}
	if _, err := converter.Convert(1, coreCooking.Piece, coreCooking.Gram); !errors.Is(err, coreCooking.ErrIncompatibleUnits) {
		t.Fatalf("expected ErrIncompatibleUnits, got %v", err)
	}
}

func TestSimpleUnitConverter_Convert_UnknownUnit(t *testing.T) {
	converter := appCooking.NewSimpleUnitConverter()

	if _, err := converter.Convert(1, coreCooking.Unit("furlong"), coreCooking.Gram); !errors.Is(err, coreCooking.ErrUnknownUnit) {
		t.Fatalf("expected ErrUnknownUnit, got %v", err)
	}
	if _, err := converter.Convert(1, coreCooking.Unit("furlong"), coreCooking.Unit("furlong")); !errors.Is(err, coreCooking.ErrUnknownUnit) {
		t.Fatalf("expected ErrUnknownUnit for an unknown identity conversion, got %v", err)
	}
}

func TestSimpleUnitConverter_Compatible(t *testing.T) {
	converter := appCooking.NewSimpleUnitConverter()

	if !converter.Compatible(coreCooking.Gram, coreCooking.Pound) {
		t.Error("expected mass units to be compatible")
	}
	if converter.Compatible(coreCooking.Gram, coreCooking.Liter) {
		t.Error("expected mass and volume units to be incompatible")
	}
	if converter.Compatible(coreCooking.Unit("furlong"), coreCooking.Gram) {
		t.Error("expected an unknown unit to be incompatible")
	}
}
