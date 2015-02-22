package strfmt

import (
	"regexp"
	"time"
)

func init() {
	d := Date{}
	Default.Add("date", &d, IsDate)
}

// IsDate returns true when the string is a valid date
func IsDate(str string) bool {
	matches := rxDate.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return false
	}
	m := matches[0]
	return !(m[2] < "01" || m[2] > "12" || m[3] < "01" || m[3] > "31")
}

const (
	// RFC3339FullDate represents a full-date as specified by RFC3339
	// See: http://goo.gl/xXOvVd
	RFC3339FullDate = "2006-01-02"
	// DatePattern pattern to match for the date format from http://tools.ietf.org/html/rfc3339#section-5.6
	DatePattern = `^([0-9]{4})-([0-9]{2})-([0-9]{2})`
)

var (
	rxDate = regexp.MustCompile(DatePattern)
)

// Date represents a date from the API
type Date struct {
	time.Time
}

// String converts this date into a string
func (d Date) String() string {
	return d.Format(RFC3339FullDate)
}

// UnmarshalText parses a text representation into a date type
func (d *Date) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	dd, err := time.Parse(RFC3339FullDate, string(text))
	if err != nil {
		return err
	}
	*d = Date{Time: dd}
	return nil
}

// MarshalText serializes this date type to string
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}
