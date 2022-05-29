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
	"math/rand"
	"sort"
)

// Sample returns sample from input with replacement or without
func Sample(input Float64Data, takenum int, replacement bool) ([]float64, error) {

	if input.Len() == 0 {
		return nil, EmptyInputErr
	}

	length := input.Len()
	if replacement {

		result := Float64Data{}
		rand.Seed(unixnano())

		// In every step, randomly take the num for
		for i := 0; i < takenum; i++ {
			idx := rand.Intn(length)
			result = append(result, input[idx])
		}

		return result, nil

	} else if !replacement && takenum <= length {

		rand.Seed(unixnano())

		// Get permutation of number of indexies
		perm := rand.Perm(length)
		result := Float64Data{}

		// Get element of input by permutated index
		for _, idx := range perm[0:takenum] {
			result = append(result, input[idx])
		}

		return result, nil

	}

	return nil, BoundsErr
}

// StableSample like stable sort, it returns samples from input while keeps the order of original data.
func StableSample(input Float64Data, takenum int) ([]float64, error) {
	if input.Len() == 0 {
		return nil, EmptyInputErr
	}

	length := input.Len()

	if takenum <= length {

		rand.Seed(unixnano())

		perm := rand.Perm(length)
		perm = perm[0:takenum]
		// Sort perm before applying
		sort.Ints(perm)
		result := Float64Data{}

		for _, idx := range perm {
			result = append(result, input[idx])
		}

		return result, nil

	}

	return nil, BoundsErr
}
