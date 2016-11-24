package flagext

import (
	"github.com/docker/go-units"
)

// ByteSize used to pass byte sizes to a go-flags CLI
type ByteSize int

// MarshalFlag implements go-flags Marshaller interface
func (b ByteSize) MarshalFlag() (string, error) {
	return units.HumanSize(float64(b)), nil
}

// UnmarshalFlag implements go-flags Unmarshaller interface
func (b *ByteSize) UnmarshalFlag(value string) error {
	sz, err := units.FromHumanSize(value)
	if err != nil {
		return err
	}
	*b = ByteSize(int(sz))
	return nil
}
