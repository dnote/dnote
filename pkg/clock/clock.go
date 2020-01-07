/* Copyright (C) 2019, 2020 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package clock provides an abstract layer over the standard time package
package clock

import (
	"time"
)

//TODO: use mutex to avoid race

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
	currentTime time.Time
}

// SetNow sets the current time for the mock clock
func (c *Mock) SetNow(t time.Time) {
	c.currentTime = t
}

// Now returns the current time
func (c *Mock) Now() time.Time {
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
