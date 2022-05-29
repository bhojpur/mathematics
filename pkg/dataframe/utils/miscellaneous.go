package utils

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

// Float64Seq will return a sequence of float64 values starting at start.
func Float64Seq(start, end, step float64, max ...int) []float64 {
	if len(max) > 0 && max[0] == 0 {
		return []float64{}
	}

	out := []float64{start}

	if step == 0 {
		return out
	}

	for {
		newVal := out[len(out)-1] + step

		if step > 0 {
			if newVal > end {
				break
			}
		} else {
			if newVal < end {
				break
			}
		}

		if len(max) > 0 && len(out) >= max[0] {
			break
		}

		out = append(out, newVal)
	}

	return out
}

// IntSeq will return a sequence of int values starting at start.
func IntSeq(start, end, step int, max ...int) []int {
	if len(max) > 0 && max[0] == 0 {
		return []int{}
	}

	out := []int{start}

	if step == 0 {
		return out
	}

	for {
		newVal := out[len(out)-1] + step

		if step > 0 {
			if newVal > end {
				break
			}
		} else {
			if newVal < end {
				break
			}
		}

		if len(max) > 0 && len(out) >= max[0] {
			break
		}

		out = append(out, newVal)
	}

	return out
}

// Int64Seq will return a sequence of int64 values starting at start.
func Int64Seq(start, end, step int64, max ...int) []int64 {
	if len(max) > 0 && max[0] == 0 {
		return []int64{}
	}

	out := []int64{start}

	if step == 0 {
		return out
	}

	for {
		newVal := out[len(out)-1] + step

		if step > 0 {
			if newVal > end {
				break
			}
		} else {
			if newVal < end {
				break
			}
		}

		if len(max) > 0 && len(out) >= max[0] {
			break
		}

		out = append(out, newVal)
	}

	return out
}
