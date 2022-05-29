package utime

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
	"context"
	"errors"
	"time"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

var (

	// ErrCantValidate means that the validation algorithm could not determine whether the Series was in the reverse direction or not.
	ErrCantValidate = errors.New("can't validate")

	// ErrCantReplace means that the Validation algorithm was not able to determine what a missing value should be replaced as.
	ErrCantReplace = errors.New("can't replace")

	// ErrValidationFailed means that validation of the SeriesTime failed.
	ErrValidationFailed = errors.New("validation failed")
)

// MissingValueOption sets how the ValidateSeriesTime behaves when a missing (nil) value is encountered.
type MissingValueOption int

const (
	// Tolerate will ignore a missing value.
	Tolerate MissingValueOption = 0

	// Replace will replace the missing value with a valid time.Time.
	Replace MissingValueOption = 1

	// Error will generate an error if a missing value is encountered.
	Error MissingValueOption = 2
)

// ValidateSeriesTimeOptions configures how ValidateSeriesTime behaves.
type ValidateSeriesTimeOptions struct {
	// MissingValue configures what must happen when a nil Value is encountered.
	MissingValue MissingValueOption

	// DontLock can be set to true if the Series should not be locked.
	DontLock bool
}

// ValidateSeriesTime will return an error if the SeriesTime's intervals don't match timeFreq.
func ValidateSeriesTime(ctx context.Context, ts *dataframe.SeriesTime, timeFreq string, opts ValidateSeriesTimeOptions) error {

	if !opts.DontLock {
		ts.Lock()
		defer ts.Unlock()
	}

	reverse := false

	nRows := len(ts.Values)
	if nRows == 0 {
		return nil
	}

	// Determine reverse direction
	if ts.Values[0] == nil {
		if opts.MissingValue == Error {
			return &dataframe.RowError{Row: 0, Err: ErrValidationFailed}
		}
		return ErrCantValidate
	}

	if nRows == 1 {
		return nil
	}

	var nextNonNilVal *time.Time
	for i, v := range ts.Values {
		if err := ctx.Err(); err != nil {
			return err
		}

		if i == 0 {
			continue
		}
		if v != nil {
			nextNonNilVal = v
			break
		}
	}
	if nextNonNilVal == nil {
		// could not find a non-nil value
		switch opts.MissingValue {
		case Tolerate:
			return nil
		case Replace:
			return &dataframe.RowError{Row: 1, Err: ErrCantReplace}
		default:
			return &dataframe.RowError{Row: 1, Err: ErrValidationFailed}
		}
	}

	if (*ts.Values[0]).Equal(*nextNonNilVal) {
		return &dataframe.RowError{Row: 1, Err: ErrValidationFailed}
	} else if (*ts.Values[0]).After(*nextNonNilVal) {
		reverse = true
	}

	type rv struct {
		row    int
		repVal time.Time
	}

	rvs := []rv{}

	// Perform main validation
	gen, err := TimeIntervalGenerator(timeFreq)
	if err != nil {
		return err
	}
	ntg := gen(*ts.Values[0], reverse)

	for row, actualTime := range ts.Values {
		if err := ctx.Err(); err != nil {
			return err
		}

		expectedTime := ntg()
		if actualTime == nil {
			if opts.MissingValue == Error {
				return &dataframe.RowError{Row: row, Err: ErrValidationFailed}
			} else if opts.MissingValue == Replace {
				rvs = append(rvs, rv{row: row, repVal: expectedTime})
			}
		} else {
			if !expectedTime.Equal(*actualTime) {
				return &dataframe.RowError{Row: row, Err: ErrValidationFailed}
			}
		}
	}

	// Replace values
	for idx := len(rvs) - 1; idx >= 0; idx-- {
		if err := ctx.Err(); err != nil {
			return err
		}
		ts.Update(rvs[idx].row, &rvs[idx].repVal, dataframe.Options{DontLock: !opts.DontLock})
	}

	return nil
}
