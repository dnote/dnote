/* Copyright 2025 Dnote Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package clock provides an abstract layer over the standard time package
package clock

import (
	"sync"
	"time"
)

// Clock is an interface to the standard library time.
// It is used to implement a real or a mock clock. The latter is used in tests.
type Clock interface {
	Now() time.Time
}

type clock struct{}

func (c *clock) Now() time.Time {
	return time.Now()
}

// Mock is a mock instance of clock
type Mock struct {
	mu          sync.RWMutex
	currentTime time.Time
}

// SetNow sets the current time for the mock clock
func (c *Mock) SetNow(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentTime = t
}

// Now returns the current time
func (c *Mock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentTime
}

// New returns an instance of a real clock
func New() Clock {
	return &clock{}
}

// NewMock returns an instance of a mock clock
func NewMock() *Mock {
	return &Mock{
		currentTime: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
}
