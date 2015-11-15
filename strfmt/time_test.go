
// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strfmt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {

	p, _ := time.Parse(time.RFC3339Nano, "2011-08-18T19:03:37.000000000+01:00")
	data := []struct {
		in  string
		out time.Time
		str string
	}{
		{"2014-12-15T08:00:00.000Z", time.Date(2014, 12, 15, 8, 0, 0, 0, time.UTC), "2014-12-15T08:00:00.000Z"},
		{"2011-08-18T19:03:37.000000000+01:00", time.Date(2011, 8, 18, 19, 3, 37, 0, p.Location()), "2011-08-18T19:03:37.000+01:00"},
		{"2014-12-15T19:30:20Z", time.Date(2014, 12, 15, 19, 30, 20, 0, time.UTC), "2014-12-15T19:30:20.000Z"},
	}

	for _, example := range data {
		parsed, err := ParseDateTime(example.in)
		assert.NoError(t, err)
		assert.Equal(t, example.out.String(), parsed.Time.String(), "Failed to parse "+example.in)
		assert.Equal(t, example.str, parsed.String())
		mt, err := parsed.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, []byte(example.str), mt)
		pp := DateTime{}
		err = pp.UnmarshalText(mt)
		assert.NoError(t, err)
		assert.Equal(t, example.out.String(), pp.Time.String())

		pp = DateTime{}
		err = pp.Scan(example.in)
		assert.NoError(t, err)
		assert.Equal(t, DateTime{example.out}, pp)
	}

	_, err := ParseDateTime("yada")
	assert.Error(t, err)

	parsed, err := ParseDateTime("")
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Unix(0, 0), parsed.Time, 0)

	pp := DateTime{}
	err = pp.UnmarshalText([]byte{})
	assert.NoError(t, err)
	err = pp.UnmarshalText([]byte("yada"))
	assert.Error(t, err)
}
