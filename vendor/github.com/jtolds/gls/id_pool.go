
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

package gls

// though this could probably be better at keeping ids smaller, the goal of
// this class is to keep a registry of the smallest unique integer ids
// per-process possible

import (
	"sync"
)

type idPool struct {
	mtx      sync.Mutex
	released []uint
	max_id   uint
}

func (p *idPool) Acquire() (id uint) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if len(p.released) > 0 {
		id = p.released[len(p.released)-1]
		p.released = p.released[:len(p.released)-1]
		return id
	}
	id = p.max_id
	p.max_id++
	return id
}

func (p *idPool) Release(id uint) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.released = append(p.released, id)
}
