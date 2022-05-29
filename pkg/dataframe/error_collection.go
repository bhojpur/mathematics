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
	"sync"
)

// ErrorCollection is used to hold multiple errors.
type ErrorCollection struct {
	sync.Mutex
	errors []error
}

// NewErrorCollection returns a new ErrorCollection.
// ErrorCollection is compatible with errors.Is and
// errors.As functions.
func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{}
}

// AddError inserts a new error to the ErrorCollection.
// If err is nil, the function will panic.
func (ec *ErrorCollection) AddError(err error, lock ...bool) {
	if err == nil {
		panic("error must not be nil")
	}
	if len(lock) == 0 || lock[0] == true {
		// default
		ec.Lock()
		defer ec.Unlock()
	}
	ec.errors = append(ec.errors, err)
}

// IsNil returns whether the ErrorCollection contains any errors.
func (ec *ErrorCollection) IsNil(lock ...bool) bool {
	if len(lock) == 0 || lock[0] == true {
		// default
		ec.Lock()
		defer ec.Unlock()
	}
	return len(ec.errors) == 0
}

// Error implements the error interface.
func (ec *ErrorCollection) Error() string {
	ec.Lock()
	defer ec.Unlock()
	var out string
	for i, err := range ec.errors {
		if i != len(ec.errors)-1 {
			out = out + fmt.Sprintf("%v\n", err)
		} else {
			out = out + fmt.Sprintf("%v", err)
		}
	}
	return out
}

// Is returns true if ErrorCollection contains err.
func (ec *ErrorCollection) Is(err error) bool {
	ec.Lock()
	defer ec.Unlock()
	if err == nil && len(ec.errors) == 0 {
		return true
	}
	for _, v := range ec.errors {
		if errors.Is(v, err) {
			return true
		}
	}
	return false
}

// As returns true if ErrorCollection contains an error
// of the same type as target. If so, it will set target
// to that error.
func (ec *ErrorCollection) As(target interface{}) bool {
	ec.Lock()
	defer ec.Unlock()
	if target == nil {
		panic("errors: target cannot be nil")
	}
	for _, v := range ec.errors {
		if errors.As(v, target) {
			return true
		}
	}
	return false
}
