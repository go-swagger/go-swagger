package generate

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeprecated(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stdout)
	s := Server{
		WithContext: true,
	}
	_, err := s.getOpts()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "warning")
}
