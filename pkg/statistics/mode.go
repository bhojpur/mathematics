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

// Mode gets the mode [most frequent value(s)] of a slice of float64s
func Mode(input Float64Data) (mode []float64, err error) {
	// Return the input if there's only one number
	l := input.Len()
	if l == 1 {
		return input, nil
	} else if l == 0 {
		return nil, EmptyInputErr
	}

	c := sortedCopyDif(input)
	// Traverse sorted array,
	// tracking the longest repeating sequence
	mode = make([]float64, 5)
	cnt, maxCnt := 1, 1
	for i := 1; i < l; i++ {
		switch {
		case c[i] == c[i-1]:
			cnt++
		case cnt == maxCnt && maxCnt != 1:
			mode = append(mode, c[i-1])
			cnt = 1
		case cnt > maxCnt:
			mode = append(mode[:0], c[i-1])
			maxCnt, cnt = cnt, 1
		default:
			cnt = 1
		}
	}
	switch {
	case cnt == maxCnt:
		mode = append(mode, c[l-1])
	case cnt > maxCnt:
		mode = append(mode[:0], c[l-1])
		maxCnt = cnt
	}

	// Since length must be greater than 1,
	// check for slices of distinct values
	if maxCnt == 1 || len(mode)*maxCnt == l && maxCnt != l {
		return Float64Data{}, nil
	}

	return mode, nil
}
