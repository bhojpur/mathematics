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
	"fmt"
	"reflect"
)

// checkConcreteType checks if ct is the zero-value of a concrete data type
func checkConcreteType(ct interface{}) error {

	if ct == nil {
		return fmt.Errorf("%v is not a valid concrete type", ct)
	}

	// Reject unacceptable types
	s := reflect.ValueOf(ct)
	switch s.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.Struct:

		// Make sure concrete type is zero value
		if !reflect.DeepEqual(ct, reflect.Zero(reflect.TypeOf(ct)).Interface()) {
			return fmt.Errorf("%v is not the zero value", ct)
		}
	default:
		return fmt.Errorf("%T is not a valid concrete type", ct)
	}

	return nil
}

// checkValue checks if an input value is valid.
// In order to be valid, it must be either nil or
// a non-pointer of the same type as the series's
// concrete type.
func (s *SeriesGeneric) checkValue(v interface{}) error {

	if v == nil {
		return nil
	}

	// Check if v is a pointer
	if reflect.TypeOf(v) != reflect.TypeOf(s.concreteType) {
		return fmt.Errorf("%v: value is invalid type", v)
	}

	return nil
}
