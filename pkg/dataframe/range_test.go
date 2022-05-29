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
	"testing"
)

type tcase struct {
	Start *int
	End   *int
	ExpN  int
	ExpS  int
	ExpE  int
}

func TestRange(t *testing.T) {

	vals := []int{0, 1, 2, 3}

	N := len(vals)

	i := func(i int) *int {
		return &i
	}

	tests := []tcase{
		{
			Start: nil,
			End:   nil,
			ExpN:  4,
			ExpS:  0,
			ExpE:  3,
		},
		{
			Start: i(1),
			End:   i(3),
			ExpN:  3,
			ExpS:  1,
			ExpE:  3,
		},
		{
			Start: nil,
			End:   i(-1),
			ExpN:  4,
			ExpS:  0,
			ExpE:  3,
		},
		{
			Start: nil,
			End:   i(-2),
			ExpN:  3,
			ExpS:  0,
			ExpE:  2,
		},
		{
			Start: i(-3),
			End:   i(-2),
			ExpN:  2,
			ExpS:  1,
			ExpE:  2,
		},
	}

	for i, tc := range tests {

		rng := &Range{Start: tc.Start, End: tc.End}

		nrows, err := rng.NRows(N)
		if err != nil {
			panic(err)
		}
		if nrows != tc.ExpN {
			t.Errorf("%d: |got: %v |expected: %v", i, nrows, tc.ExpN)
		}

		s, e, err := rng.Limits(N)
		if err != nil {
			panic(err)
		}
		if s != tc.ExpS || e != tc.ExpE {
			t.Errorf("%d: |got: %v,%v |expected: %v,%v", i, s, e, tc.ExpS, tc.ExpE)
		}
	}

}
