package dataframe

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"sort"
)

// IsEqualFunc is used to determine if a and b are considered equal.
type IsEqualFunc func(a, b interface{}) bool

// IsLessThanFunc returns true if a < b
type IsLessThanFunc func(a, b interface{}) bool

// SortKey is the key to sort a Dataframe
type SortKey struct {

	// Key can be an int (position of series) or string (name of series).
	Key interface{}

	// Desc can be set to sort in descending order.
	Desc bool

	seriesIndex int
}

type sorter struct {
	keys []SortKey
	df   *DataFrame
	ctx  context.Context
}

func (s *sorter) Len() int {
	return s.df.n
}

func (s *sorter) Less(i, j int) bool {

	if err := s.ctx.Err(); err != nil {
		panic(err)
	}

	for _, key := range s.keys {
		series := s.df.Series[key.seriesIndex]

		left := series.Value(i)
		right := series.Value(j)

		// Check if left and right are equal
		if series.IsEqualFunc(left, right) {
			continue
		} else {
			if key.Desc {
				// Sort in descending order
				return !series.IsLessThanFunc(left, right)
			}
			return series.IsLessThanFunc(left, right)
		}
	}

	return false
}

func (s *sorter) Swap(i, j int) {
	s.df.Swap(i, j, DontLock)
}

// SortOptions is used to configure the sort algorithm for a Dataframe or Series
type SortOptions struct {

	// Stable can be set if the original order of equal items must be maintained.
	//
	// See: https://golang.org/pkg/sort/#Stable
	Stable bool

	// Desc can be set to sort in descending order. This option is ignored when applied to a Dataframe.
	// Only use it with a Series.
	Desc bool

	// DontLock can be set to true if the Series should not be locked.
	DontLock bool
}

// Sort is used to sort the Dataframe according to different keys.
// It will return true if sorting was completed or false when the context is canceled.
func (df *DataFrame) Sort(ctx context.Context, keys []SortKey, opts ...SortOptions) (completed bool) {
	if len(keys) == 0 {
		return true
	}

	defer func() {
		if x := recover(); x != nil {
			if x == context.Canceled || x == context.DeadlineExceeded {
				completed = false
			} else {
				panic(x)
			}
		}
	}()

	if len(opts) == 0 || !opts[0].DontLock {
		// Default
		df.lock.Lock()
		defer df.lock.Unlock()
	}

	// Clear seriesIndex from keys
	defer func() {
		for i := range keys {
			key := &keys[i]
			key.seriesIndex = 0
		}
	}()

	// Convert keys to index
	for i := range keys {
		key := &keys[i]

		name, ok := key.Key.(string)
		if ok {
			col, err := df.NameToColumn(name, dontLock)
			if err != nil {
				panic(err)
			}
			key.seriesIndex = col
		} else {
			key.seriesIndex = key.Key.(int)
		}
	}

	s := &sorter{
		keys: keys,
		df:   df,
		ctx:  ctx,
	}

	if len(opts) == 0 || !opts[0].Stable {
		// Default
		sort.Sort(s)
	} else {
		sort.Stable(s)
	}

	return true
}
