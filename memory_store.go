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
	"strings"

	"github.com/gomqtt/packet"
)

// MemoryStore organizes packets in memory.
type MemoryStore struct {
	store map[string]packet.Packet
	mutex sync.Mutex
}

// NewMemoryStore returns a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]packet.Packet),
	}
}

// Put will store the specified packet in the store.
func (s *MemoryStore) Put(dir string, pkt packet.Packet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id, ok := packet.PacketID(pkt)
	if ok {
		s.store[s.key(dir, id)] = pkt
	}

	return nil
}

// Get will retrieve and return a packet by its id.
func (s *MemoryStore) Get(dir string, id uint16) (packet.Packet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.store[s.key(dir, id)], nil
}

// Del will remove a packet using its id.
func (s *MemoryStore) Del(dir string, id uint16) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, s.key(dir, id))

	return nil
}

// All will return all stored packets.
func (s *MemoryStore) All(dir string) ([]packet.Packet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	all := make([]packet.Packet, 0)

	for key, pkt := range s.store {
		if strings.HasPrefix(key, dir) {
			all = append(all, pkt)
		}
	}

	return all, nil
}

// Reset will wipe all packets currently stored.
func (s *MemoryStore) Reset() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store = make(map[string]packet.Packet)

	return nil
}

// return a string key based on direction and id
func (s *MemoryStore) key(dir string, id uint16) string {
	return dir + "-" + string(id)
}