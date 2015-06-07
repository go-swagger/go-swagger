package strfmt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	d := Duration(0)
	Default.Add("duration", &d, IsDuration)
}

var (
	timeUnits = [][]string{
		[]string{"ns", "nano"},
		[]string{"us", "µs", "micro"},
		[]string{"ms", "milli"},
		[]string{"s", "sec"},
		[]string{"m", "min"},
		[]string{"h", "hr", "hour"},
		[]string{"d", "day"},
		[]string{"w", "wk", "week"},
	}

	timeMultiplier = map[string]time.Duration{
		"ns": time.Nanosecond,
		"us": time.Microsecond,
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  24 * time.Hour,
		"w":  7 * 24 * time.Hour,
	}

	durationMatcher = regexp.MustCompile(`((\d+)\s*([A-Za-zµ]+))`)
)

// IsDuration returns true if the provided string is a valid duration
func IsDuration(str string) bool {
	_, err := ParseDuration(str)
	return err == nil
}

// Duration represents a duration
//
// swagger:strfmt duration
type Duration time.Duration

// MarshalText turns this instance into text
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// UnmarshalText hydrates this instance from text
func (d *Duration) UnmarshalText(data []byte) error { // validation is performed later on
	dd, err := ParseDuration(string(data))
	if err != nil {
		return err
	}
	*d = Duration(dd)
	return nil
}

// ParseDuration parses a duration from a string, compatible with scala duration syntax
func ParseDuration(cand string) (time.Duration, error) {
	if dur, err := time.ParseDuration(cand); err == nil {
		return dur, nil
	}

	var dur time.Duration
	ok := false
	for _, match := range durationMatcher.FindAllStringSubmatch(cand, -1) {

		factor, err := strconv.Atoi(match[2]) // converts string to int
		if err != nil {
			return 0, err
		}
		unit := strings.ToLower(strings.TrimSpace(match[3]))

		for _, variants := range timeUnits {
			last := len(variants) - 1
			multiplier := timeMultiplier[variants[0]]

			for i, variant := range variants {
				if (last == i && strings.HasPrefix(unit, variant)) || strings.EqualFold(variant, unit) {
					ok = true
					dur += (time.Duration(factor) * multiplier)
				}
			}
		}
	}

	if ok {
		return dur, nil
	}
	return 0, fmt.Errorf("Unable to parse %s as duration", cand)
}
