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
	"math"
)

// Percentile finds the relative standing in a slice of floats
func Percentile(input Float64Data, percent float64) (percentile float64, err error) {
	length := input.Len()
	if length == 0 {
		return math.NaN(), EmptyInputErr
	}

	if length == 1 {
		return input[0], nil
	}

	if percent <= 0 || percent > 100 {
		return math.NaN(), BoundsErr
	}

	// Start by sorting a copy of the slice
	c := sortedCopy(input)

	// Multiply percent by length of input
	index := (percent / 100) * float64(len(c))

	// Check if the index is a whole number
	if index == float64(int64(index)) {

		// Convert float to int
		i := int(index)

		// Find the value at the index
		percentile = c[i-1]

	} else if index > 1 {

		// Convert float to int via truncation
		i := int(index)

		// Find the average of the index and following values
		percentile, _ = Mean(Float64Data{c[i-1], c[i]})

	} else {
		return math.NaN(), BoundsErr
	}

	return percentile, nil

}

// PercentileNearestRank finds the relative standing in a slice of floats using the Nearest Rank method
func PercentileNearestRank(input Float64Data, percent float64) (percentile float64, err error) {

	// Find the length of items in the slice
	il := input.Len()

	// Return an error for empty slices
	if il == 0 {
		return math.NaN(), EmptyInputErr
	}

	// Return error for less than 0 or greater than 100 percentages
	if percent < 0 || percent > 100 {
		return math.NaN(), BoundsErr
	}

	// Start by sorting a copy of the slice
	c := sortedCopy(input)

	// Return the last item
	if percent == 100.0 {
		return c[il-1], nil
	}

	// Find ordinal ranking
	or := int(math.Ceil(float64(il) * percent / 100))

	// Return the item that is in the place of the ordinal rank
	if or == 0 {
		return c[0], nil
	}
	return c[or-1], nil

}
