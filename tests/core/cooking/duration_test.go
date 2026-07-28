package cooking

import (
	"errors"
	"testing"
	"time"

	"github.com/tthung1997/buddy/core/cooking"
)

func TestDurationUnit_ToDuration(t *testing.T) {
	cases := []struct {
		unit     cooking.DurationUnit
		amount   float64
		expected time.Duration
	}{
		{cooking.Millisecond, 250, 250 * time.Millisecond},
		{cooking.Second, 90, 90 * time.Second},
		{cooking.Minute, 45, 45 * time.Minute},
		{cooking.Minute, 1.5, 90 * time.Second},
		{cooking.Hour, 2, 2 * time.Hour},
		{cooking.Day, 1, 24 * time.Hour},
		{cooking.Week, 1, 7 * 24 * time.Hour},
		{cooking.Month, 1, 30 * 24 * time.Hour},
		{cooking.Year, 1, 365 * 24 * time.Hour},
	}

	for _, testCase := range cases {
		duration, err := testCase.unit.ToDuration(testCase.amount)
		if err != nil {
			t.Fatalf("unexpected error for %v %s: %v", testCase.amount, testCase.unit, err)
		}
		if duration != testCase.expected {
			t.Errorf("expected %v %s to be %v, got %v", testCase.amount, testCase.unit, testCase.expected, duration)
		}
	}
}

func TestDurationUnit_ToDuration_UnknownUnit(t *testing.T) {
	if _, err := cooking.DurationUnit("fortnight").ToDuration(1); !errors.Is(err, cooking.ErrUnknownDurationUnit) {
		t.Fatalf("expected ErrUnknownDurationUnit, got %v", err)
	}
}

func TestDurationUnit_ToDuration_NegativeAmount(t *testing.T) {
	if _, err := cooking.Minute.ToDuration(-1); err == nil {
		t.Fatal("expected an error for a negative duration")
	}
}

func TestDurationUnits(t *testing.T) {
	units := cooking.DurationUnits()
	if len(units) != 8 {
		t.Fatalf("expected 8 duration units, got %d", len(units))
	}
	for _, unit := range units {
		if !unit.IsValid() {
			t.Errorf("expected %s to be valid", unit)
		}
	}

	units[0] = cooking.DurationUnit("mutated")
	if cooking.DurationUnits()[0] == cooking.DurationUnit("mutated") {
		t.Error("expected DurationUnits to return a defensive copy")
	}
}
