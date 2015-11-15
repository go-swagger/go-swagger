
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

package flags

import (
	"strconv"
)

type multiTag struct {
	value string
	cache map[string][]string
}

func newMultiTag(v string) multiTag {
	return multiTag{
		value: v,
	}
}

func (x *multiTag) scan() (map[string][]string, error) {
	v := x.value

	ret := make(map[string][]string)

	// This is mostly copied from reflect.StructTag.Get
	for v != "" {
		i := 0

		// Skip whitespace
		for i < len(v) && v[i] == ' ' {
			i++
		}

		v = v[i:]

		if v == "" {
			break
		}

		// Scan to colon to find key
		i = 0

		for i < len(v) && v[i] != ' ' && v[i] != ':' && v[i] != '"' {
			i++
		}

		if i >= len(v) {
			return nil, newErrorf(ErrTag, "expected `:' after key name, but got end of tag (in `%v`)", x.value)
		}

		if v[i] != ':' {
			return nil, newErrorf(ErrTag, "expected `:' after key name, but got `%v' (in `%v`)", v[i], x.value)
		}

		if i+1 >= len(v) {
			return nil, newErrorf(ErrTag, "expected `\"' to start tag value at end of tag (in `%v`)", x.value)
		}

		if v[i+1] != '"' {
			return nil, newErrorf(ErrTag, "expected `\"' to start tag value, but got `%v' (in `%v`)", v[i+1], x.value)
		}

		name := v[:i]
		v = v[i+1:]

		// Scan quoted string to find value
		i = 1

		for i < len(v) && v[i] != '"' {
			if v[i] == '\n' {
				return nil, newErrorf(ErrTag, "unexpected newline in tag value `%v' (in `%v`)", name, x.value)
			}

			if v[i] == '\\' {
				i++
			}
			i++
		}

		if i >= len(v) {
			return nil, newErrorf(ErrTag, "expected end of tag value `\"' at end of tag (in `%v`)", x.value)
		}

		val, err := strconv.Unquote(v[:i+1])

		if err != nil {
			return nil, newErrorf(ErrTag, "Malformed value of tag `%v:%v` => %v (in `%v`)", name, v[:i+1], err, x.value)
		}

		v = v[i+1:]

		ret[name] = append(ret[name], val)
	}

	return ret, nil
}

func (x *multiTag) Parse() error {
	vals, err := x.scan()
	x.cache = vals

	return err
}

func (x *multiTag) cached() map[string][]string {
	if x.cache == nil {
		cache, _ := x.scan()

		if cache == nil {
			cache = make(map[string][]string)
		}

		x.cache = cache
	}

	return x.cache
}

func (x *multiTag) Get(key string) string {
	c := x.cached()

	if v, ok := c[key]; ok {
		return v[len(v)-1]
	}

	return ""
}

func (x *multiTag) GetMany(key string) []string {
	c := x.cached()
	return c[key]
}

func (x *multiTag) Set(key string, value string) {
	c := x.cached()
	c[key] = []string{value}
}

func (x *multiTag) SetMany(key string, value []string) {
	c := x.cached()
	c[key] = value
}
