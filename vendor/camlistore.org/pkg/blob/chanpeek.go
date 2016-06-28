/*
Copyright 2011 Google Inc.

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

package blob

// ChanPeeker wraps a channel receiving SizedRefs and adds a 1 element
// buffer on it, with Peek and Take methods.
type ChanPeeker struct {
	Ch <-chan SizedRef

	// Invariant after a call to Peek: either peekok or closed.
	closed bool
	peekok bool // whether peek is valid
	peek   SizedRef
}

// MustPeek returns the next SizedRef or panics if none is available.
func (cp *ChanPeeker) MustPeek() SizedRef {
	sr, ok := cp.Peek()
	if !ok {
		panic("No Peek value available")
	}
	return sr
}

// Peek returns the next SizedRef and whether one was available.
func (cp *ChanPeeker) Peek() (sr SizedRef, ok bool) {
	if cp.closed {
		return
	}
	if cp.peekok {
		return cp.peek, true
	}
	v, ok := <-cp.Ch
	if !ok {
		cp.closed = true
		return
	}
	cp.peek = v
	cp.peekok = true
	return cp.peek, true
}

// Closed reports true if no more SizedRef values are available.
func (cp *ChanPeeker) Closed() bool {
	cp.Peek()
	return cp.closed
}

// MustTake returns the next SizedRef, else panics if none is available.
func (cp *ChanPeeker) MustTake() SizedRef {
	sr, ok := cp.Take()
	if !ok {
		panic("MustTake called on empty channel")
	}
	return sr
}

// Take returns the next SizedRef and whether one was available for the taking.
func (cp *ChanPeeker) Take() (sr SizedRef, ok bool) {
	v, ok := cp.Peek()
	if !ok {
		return
	}
	cp.peekok = false
	return v, true
}

// ConsumeAll drains the channel of all items.
func (cp *ChanPeeker) ConsumeAll() {
	for !cp.Closed() {
		cp.Take()
	}
}
