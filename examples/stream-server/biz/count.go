package biz

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-swagger/go-swagger/examples/stream-server/models"
)

// MyCounter is the concrete implementation
type MyCounter struct{}

// Down is the concrete implementation that spits out the JSON bodies
func (mc *MyCounter) Down(max int64, w io.Writer) error {
	if max == 11 {
		return errors.New("we don't *do* elevensies")
	}
	e := json.NewEncoder(w)
	for ix := int64(0); ix <= max; ix++ {
		r := max - ix
		fmt.Printf("Iteration %d\n", r)
		_ = e.Encode(models.Mark{Remains: &r})
		if ix != max {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
