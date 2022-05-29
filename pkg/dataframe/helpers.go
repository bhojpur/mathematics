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
	"errors"
	"fmt"
)

// ErrNoRows signifies that the Series, Dataframe or import data
// contains no rows of data.
var ErrNoRows = errors.New("contains no rows")

const (
	// FALSE is used convert a false (bool) to an int.
	FALSE = 0
	// TRUE is used convert a true (bool) to an int.
	TRUE = 1
)

// B converts a boolean to an int.
func B(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IsValidFloat64 returns true if f is neither Nan nor ±Inf.
// Otherwise it returns false.
func IsValidFloat64(f float64) bool {

	if isNaN(f) {
		return false
	}

	if isInf(f, 0) {
		return false
	}

	return true
}

// BoolValueFormatter is used by SetValueToStringFormatter
// to display an int as a bool. If the encountered value
// is not a 0 or 1, it will panic.
func BoolValueFormatter(v interface{}) string {
	if v == nil {
		return "NaN"
	}

	str := fmt.Sprintf("%v", v)
	switch str {
	case "0":
		return "false"
	case "1":
		return "true"
	default:
		_ = v.(bool) // Intentionally panic
		return ""
	}
}

// DontLock is short-hand for various functions that permit disabling locking.
var DontLock = dontLock
var dontLock = Options{DontLock: true}
