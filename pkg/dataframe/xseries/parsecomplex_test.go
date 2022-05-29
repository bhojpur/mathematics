package xseries

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
	str       string
	expAnswer complex128
	expErr    error
}

func TestParseComplex(t *testing.T) {

	tests := []tcase{
		{
			str:       "99",
			expAnswer: complex(99, 0),
		},
		{
			str:       "99",
			expAnswer: complex(99, 0),
		},
		{
			str:       "-99",
			expAnswer: complex(-99, 0),
		},
		{
			str:       "1i",
			expAnswer: complex(0, 1),
		},
		{
			str:       "-1i",
			expAnswer: complex(0, -1),
		},
		{
			str:       "3-1i",
			expAnswer: complex(3, -1),
		},
		{
			str:       "3+1i",
			expAnswer: complex(3, 1),
		},
		{
			str:       "3-1i",
			expAnswer: complex(3, -1),
		},
		{
			str:       "3+1i",
			expAnswer: complex(3, 1),
		},
		{
			str:       "1i",
			expAnswer: complex(0, 1),
		},
		{
			str:       "-1i",
			expAnswer: complex(0, -1),
		},
		{
			str:       "3e3-1i",
			expAnswer: complex(3e3, -1),
		},
		{
			str:       "-3e3-1i",
			expAnswer: complex(-3e3, -1),
		},
		{
			str:       "3e3-1i",
			expAnswer: complex(3e3, -1),
		},
		{
			str:       "3e+3-1i",
			expAnswer: complex(3e+3, -1),
		},
		{
			str:       "-3e+3-1i",
			expAnswer: complex(-3e+3, -1),
		},
		{
			str:       "-3e+3-1i",
			expAnswer: complex(-3e+3, -1),
		},
		{
			str:       "3e+3-3e+3i",
			expAnswer: complex(3e+3, -3e+3),
		},
		{
			str:       "3e+3+3e+3i",
			expAnswer: complex(3e+3, 3e+3),
		},
	}

	for i, tc := range tests {

		got, gotErr := parseComplex(tc.str)
		if gotErr != nil {
			if tc.expErr == nil {
				t.Errorf("%d: |got: %v |expected: %v", i, gotErr, tc.expErr)
			}
		} else {
			if tc.expErr != nil {
				t.Errorf("%d: |got: %v |expected: %v", i, got, tc.expErr)
			} else {
				if got != tc.expAnswer {
					t.Errorf("%d: |got: %v |expected: %v", i, got, tc.expAnswer)
				}
			}
		}
	}

}
