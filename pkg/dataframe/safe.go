//go:build js || appengine || safe
// +build js appengine safe

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

import "math"

// nan returns NaN.
// See: https://golang.org/pkg/math/#NaN
func nan() float64 {
	return math.NaN()
}

// isNaN returns whether f is NaN.
// See: https://golang.org/pkg/math/#IsNaN
func isNaN(f float64) bool {
	return f != f
}

// isInf returns whether f is +Inf or -Inf.
func isInf(f float64, sign int) bool {
	return sign >= 0 && f > 1.797693134862315708145274237317043567981e+308 || sign <= 0 && f < -1.797693134862315708145274237317043567981e+308
}
