package cooking

import (
	"errors"
	"fmt"
	"time"
)

var ErrUnknownDurationUnit = errors.New("unknown duration unit")

type DurationUnit string

const (
	Millisecond DurationUnit = "millisecond"
	Second      DurationUnit = "second"
	Minute      DurationUnit = "minute"
	Hour        DurationUnit = "hour"
	Day         DurationUnit = "day"
	Week        DurationUnit = "week"
	Month       DurationUnit = "month"
	Year        DurationUnit = "year"
)

// durationUnitLengths approximates calendar units so a step duration can always
// be reduced to a wall clock duration: a month is 30 days, a year is 365 days.
var durationUnitLengths = map[DurationUnit]time.Duration{
	Millisecond: time.Millisecond,
	Second:      time.Second,
	Minute:      time.Minute,
	Hour:        time.Hour,
	Day:         24 * time.Hour,
	Week:        7 * 24 * time.Hour,
	Month:       30 * 24 * time.Hour,
	Year:        365 * 24 * time.Hour,
}

var orderedDurationUnits = []DurationUnit{
	Millisecond, Second, Minute, Hour, Day, Week, Month, Year,
}

// DurationUnits returns every supported duration unit, shortest first.
func DurationUnits() []DurationUnit {
	units := make([]DurationUnit, len(orderedDurationUnits))
	copy(units, orderedDurationUnits)
	return units
}

func (d DurationUnit) IsValid() bool {
	_, exists := durationUnitLengths[d]
	return exists
}

func (d DurationUnit) String() string {
	return string(d)
}

// ToDuration converts an amount expressed in this unit into a time.Duration.
func (d DurationUnit) ToDuration(amount float64) (time.Duration, error) {
	length, exists := durationUnitLengths[d]
	if !exists {
		return 0, fmt.Errorf("%w: %q", ErrUnknownDurationUnit, string(d))
	}
	if amount < 0 {
		return 0, fmt.Errorf("duration amount must not be negative, got %v", amount)
	}
	return time.Duration(amount * float64(length)), nil
}
