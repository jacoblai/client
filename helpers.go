// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"sync"
)

// a futureStore is used to store active Futures
type futureStore struct {
	sync.Mutex

	store map[uint16]Future
}

// newFutureStore will create a new futureStore
func newFutureStore() *futureStore {
	return &futureStore{
		store: make(map[uint16]Future),
	}
}

// put will save a Future to the store
func (s *futureStore) put(id uint16, future Future) {
	s.Lock()
	defer s.Unlock()

	s.store[id] = future
}

// get will retrieve a Future from the store
func (s *futureStore) get(id uint16) Future {
	s.Lock()
	defer s.Unlock()

	return s.store[id]
}

// del will remove a Future from the store
func (s *futureStore) del(id uint16) {
	delete(s.store, id)
}

// a counter keeps track of packet ids
type counter struct {
	sync.Mutex

	id uint16
}

// newCounter will return a new counter
func newCounter() *counter {
	return &counter{}
}

// next will generate the next packet id
func (c *counter) next() uint16 {
	c.Lock()
	defer func(){
		c.id++
		c.Unlock()
	}()

	return c.id
}
