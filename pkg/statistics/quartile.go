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

import "math"

// Quartiles holds the three quartile points
type Quartiles struct {
	Q1 float64
	Q2 float64
	Q3 float64
}

// Quartile returns the three quartile points from a slice of data
func Quartile(input Float64Data) (Quartiles, error) {

	il := input.Len()
	if il == 0 {
		return Quartiles{}, EmptyInputErr
	}

	// Start by sorting a copy of the slice
	copy := sortedCopy(input)

	// Find the cutoff places depeding on if
	// the input slice length is even or odd
	var c1 int
	var c2 int
	if il%2 == 0 {
		c1 = il / 2
		c2 = il / 2
	} else {
		c1 = (il - 1) / 2
		c2 = c1 + 1
	}

	// Find the Medians with the cutoff points
	Q1, _ := Median(copy[:c1])
	Q2, _ := Median(copy)
	Q3, _ := Median(copy[c2:])

	return Quartiles{Q1, Q2, Q3}, nil

}

// InterQuartileRange finds the range between Q1 and Q3
func InterQuartileRange(input Float64Data) (float64, error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}
	qs, _ := Quartile(input)
	iqr := qs.Q3 - qs.Q1
	return iqr, nil
}

// Midhinge finds the average of the first and third quartiles
func Midhinge(input Float64Data) (float64, error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}
	qs, _ := Quartile(input)
	mh := (qs.Q1 + qs.Q3) / 2
	return mh, nil
}

// Trimean finds the average of the median and the midhinge
func Trimean(input Float64Data) (float64, error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	c := sortedCopy(input)
	q, _ := Quartile(c)

	return (q.Q1 + (q.Q2 * 2) + q.Q3) / 4, nil
}
