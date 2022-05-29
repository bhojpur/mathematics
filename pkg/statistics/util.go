package statistics

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
	"sort"
	"time"
)

// float64ToInt rounds a float64 to an int
func float64ToInt(input float64) (output int) {
	r, _ := Round(input, 0)
	return int(r)
}

// unixnano returns nanoseconds from UTC epoch
func unixnano() int64 {
	return time.Now().UTC().UnixNano()
}

// copyslice copies a slice of float64s
func copyslice(input Float64Data) Float64Data {
	s := make(Float64Data, input.Len())
	copy(s, input)
	return s
}

// sortedCopy returns a sorted copy of float64s
func sortedCopy(input Float64Data) (copy Float64Data) {
	copy = copyslice(input)
	sort.Float64s(copy)
	return
}

// sortedCopyDif returns a sorted copy of float64s
// only if the original data isn't sorted.
// Only use this if returned slice won't be manipulated!
func sortedCopyDif(input Float64Data) (copy Float64Data) {
	if sort.Float64sAreSorted(input) {
		return input
	}
	copy = copyslice(input)
	sort.Float64s(copy)
	return
}
