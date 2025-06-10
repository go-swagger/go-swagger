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
func (mc *MyCounter) Down(maximum int64, w io.Writer) error {
	if maximum == 11 {
		return errors.New("we don't *do* elevensies")
	}
	e := json.NewEncoder(w)
	for ix := int64(0); ix <= maximum; ix++ {
		r := maximum - ix
		fmt.Printf("Iteration %d\n", r)
		_ = e.Encode(models.Mark{Remains: &r})
		if ix != maximum {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
