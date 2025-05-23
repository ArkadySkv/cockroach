// Copyright 2017 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package syncutil

import (
	"math"
	"sync/atomic"
)

// AtomicFloat64 mimics the atomic types in the sync/atomic standard library,
// but for the float64 type. If you'd like to implement additional methods,
// consider checking out the expvar Float type for guidance:
// https://golang.org/src/expvar/expvar.go?s=2188:2222#L69
type AtomicFloat64 struct {
	val atomic.Uint64
}

// Store atomically stores a float64 value.
func (f *AtomicFloat64) Store(val float64) {
	f.val.Store(math.Float64bits(val))
}

// Load atomically loads tha float64 value.
func (f *AtomicFloat64) Load() float64 {
	return math.Float64frombits(f.val.Load())
}

// Add atomically adds delta to the float64 value and returns the new value.
func (f *AtomicFloat64) Add(delta float64) (new float64) {
	for {
		oldInt := f.val.Load()
		oldFloat := math.Float64frombits(oldInt)
		newFloat := oldFloat + delta
		newInt := math.Float64bits(newFloat)
		if f.val.CompareAndSwap(oldInt, newInt) {
			return newFloat
		}
	}
}

// StoreIfHigher atomically stores the given value if it is higher than the
// current value. It returns the old value and whether a swap was carried out.
func (f *AtomicFloat64) StoreIfHigher(new float64) (old float64, stored bool) {
	newInt := math.Float64bits(new)
	for {
		oldInt := f.val.Load()
		oldFloat := math.Float64frombits(oldInt)
		if oldFloat >= new {
			return oldFloat, false
		}
		if f.val.CompareAndSwap(oldInt, newInt) {
			return oldFloat, true
		}
	}
}

// CompareAndSwap is atomically replaces the current value with 'new' if the
// existing value is 'old'. It returns whether this was the case.
func (f *AtomicFloat64) CompareAndSwap(old, new float64) (swapped bool) {
	oldInt := math.Float64bits(old)
	newInt := math.Float64bits(new)
	return f.val.CompareAndSwap(oldInt, newInt)
}

// AtomicString gives you atomic-style APIs for string.
type AtomicString struct {
	s atomic.Value
}

// Set atomically sets str as new value.
func (s *AtomicString) Set(val string) {
	s.s.Store(val)
}

// Get atomically returns the current value.
func (s *AtomicString) Get() string {
	val := s.s.Load()
	if val == nil {
		return ""
	}
	return val.(string)
}
