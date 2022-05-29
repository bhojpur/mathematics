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

type statsError struct {
	err string
}

func (s statsError) Error() string {
	return s.err
}

func (s statsError) String() string {
	return s.err
}

// These are the package-wide error values.
// All error identification should use these values.
var (
	// ErrEmptyInput Input must not be empty
	ErrEmptyInput = statsError{"Input must not be empty."}
	// ErrNaN Not a number
	ErrNaN = statsError{"Not a number."}
	// ErrNegative Must not contain negative values
	ErrNegative = statsError{"Must not contain negative values."}
	// ErrZero Must not contain zero values
	ErrZero = statsError{"Must not contain zero values."}
	// ErrBounds Input is outside of range
	ErrBounds = statsError{"Input is outside of range."}
	// ErrSize Must be the same length
	ErrSize = statsError{"Must be the same length."}
	// ErrInfValue Value is infinite
	ErrInfValue = statsError{"Value is infinite."}
	// ErrYCoord Y Value must be greater than zero
	ErrYCoord = statsError{"Y Value must be greater than zero."}
)
