package strfmt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testDurationParser(t *testing.T, toParse string, expected time.Duration) {
	r, e := ParseDuration(toParse)
	assert.NoError(t, e)
	assert.Equal(t, expected, r)
}

func TestDurationParser(t *testing.T) {

	// parse the short forms without spaces
	testDurationParser(t, "1ns", 1*time.Nanosecond)
	testDurationParser(t, "1us", 1*time.Microsecond)
	testDurationParser(t, "1µs", 1*time.Microsecond)
	testDurationParser(t, "1ms", 1*time.Millisecond)
	testDurationParser(t, "1s", 1*time.Second)
	testDurationParser(t, "1m", 1*time.Minute)
	testDurationParser(t, "1h", 1*time.Hour)
	testDurationParser(t, "1hr", 1*time.Hour)
	testDurationParser(t, "1d", 24*time.Hour)
	testDurationParser(t, "1w", 7*24*time.Hour)
	testDurationParser(t, "1wk", 7*24*time.Hour)

	// parse the long forms without spaces
	testDurationParser(t, "1nanoseconds", 1*time.Nanosecond)
	testDurationParser(t, "1nanos", 1*time.Nanosecond)
	testDurationParser(t, "1microseconds", 1*time.Microsecond)
	testDurationParser(t, "1micros", 1*time.Microsecond)
	testDurationParser(t, "1millis", 1*time.Millisecond)
	testDurationParser(t, "1milliseconds", 1*time.Millisecond)
	testDurationParser(t, "1second", 1*time.Second)
	testDurationParser(t, "1sec", 1*time.Second)
	testDurationParser(t, "1min", 1*time.Minute)
	testDurationParser(t, "1minute", 1*time.Minute)
	testDurationParser(t, "1hour", 1*time.Hour)
	testDurationParser(t, "1day", 24*time.Hour)
	testDurationParser(t, "1week", 7*24*time.Hour)

	// parse the short forms with spaces
	testDurationParser(t, "1  ns", 1*time.Nanosecond)
	testDurationParser(t, "1  us", 1*time.Microsecond)
	testDurationParser(t, "1  µs", 1*time.Microsecond)
	testDurationParser(t, "1  ms", 1*time.Millisecond)
	testDurationParser(t, "1  s", 1*time.Second)
	testDurationParser(t, "1  m", 1*time.Minute)
	testDurationParser(t, "1  h", 1*time.Hour)
	testDurationParser(t, "1  hr", 1*time.Hour)
	testDurationParser(t, "1  d", 24*time.Hour)
	testDurationParser(t, "1  w", 7*24*time.Hour)
	testDurationParser(t, "1  wk", 7*24*time.Hour)

	// parse the long forms without spaces
	testDurationParser(t, "1  nanoseconds", 1*time.Nanosecond)
	testDurationParser(t, "1  nanos", 1*time.Nanosecond)
	testDurationParser(t, "1  microseconds", 1*time.Microsecond)
	testDurationParser(t, "1  micros", 1*time.Microsecond)
	testDurationParser(t, "1  millis", 1*time.Millisecond)
	testDurationParser(t, "1  milliseconds", 1*time.Millisecond)
	testDurationParser(t, "1  second", 1*time.Second)
	testDurationParser(t, "1  sec", 1*time.Second)
	testDurationParser(t, "1  min", 1*time.Minute)
	testDurationParser(t, "1  minute", 1*time.Minute)
	testDurationParser(t, "1  hour", 1*time.Hour)
	testDurationParser(t, "1  day", 24*time.Hour)
	testDurationParser(t, "1  week", 7*24*time.Hour)
}
