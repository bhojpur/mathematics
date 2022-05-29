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

// Entropy provides calculation of the entropy
func Entropy(input Float64Data) (float64, error) {
	input, err := normalize(input)
	if err != nil {
		return math.NaN(), err
	}
	var result float64
	for i := 0; i < input.Len(); i++ {
		v := input.Get(i)
		if v == 0 {
			continue
		}
		result += (v * math.Log(v))
	}
	return -result, nil
}

func normalize(input Float64Data) (Float64Data, error) {
	sum, err := input.Sum()
	if err != nil {
		return Float64Data{}, err
	}
	for i := 0; i < input.Len(); i++ {
		input[i] = input[i] / sum
	}
	return input, nil
}
