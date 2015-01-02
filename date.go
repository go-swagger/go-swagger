package swagger

import "time"

// RFC3339FullDate represents a full-date as specified by RFC3339
// See: http://goo.gl/xXOvVd
const RFC3339FullDate = "2006-01-02"

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
