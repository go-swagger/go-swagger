/*
Copyright 2014 Zachary Klippenstein

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package regen

/*
The default Source implementation is very slow to seed. Replaced with a
64-bit xor-shift source from http://vigna.di.unimi.it/ftp/papers/xorshift.pdf.
This source seeds very quickly, and only uses a single variable, so concurrent
modification by multiple goroutines is possible.

To create a seeded source:
	randSource := xorShift64Source(mySeed)

To create a source with the default seed:
	var randSource xorShift64Source
*/
type xorShift64Source uint64

func (src *xorShift64Source) Seed(seed int64) {
	*src = xorShift64Source(seed)
}

func (src *xorShift64Source) Int63() int64 {
	// A zero seed will only generate zeros.
	if *src == 0 {
		*src = 1
	}

	*src ^= *src >> 12 // a
	*src ^= *src << 25 // b
	*src ^= *src >> 27 // c

	return int64((*src * 2685821657736338717) >> 1)
}
