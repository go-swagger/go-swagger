package generate

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/go-swagger/go-swagger/generator"
	"github.com/stretchr/testify/assert"
)

func TestDeprecated(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)
	s := Server{
		WithContext: true,
	}
	s.apply(new(generator.GenOpts))
	assert.Contains(t, buf.String(), "warning")
}
