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
	"time"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

// NewSeriesTimeOptions sets how NewSeriesTime decides at which row to stop.
type NewSeriesTimeOptions struct {

	// Size determines how many rows are generated.
	// This option can't be used with Until option.
	Size *int

	// Until is the maximum time in the generated Series.
	// This option can't be used with Size option.
	Until *time.Time
}

// NewSeriesTime will create a new SeriesTime with timeFreq prescribing the intervals between each row. Setting reverse will make the time series decrement per row.
func NewSeriesTime(ctx context.Context, name string, timeFreq string, startTime time.Time, reverse bool, opts NewSeriesTimeOptions) (*dataframe.SeriesTime, error) {

	if opts.Size != nil && opts.Until != nil {
		panic("Size and Until options can't be used together")
	}

	if opts.Size == nil && opts.Until == nil {
		panic("Either Size xor Until option required")
	}

	if opts.Until != nil {
		if reverse {
			if startTime.Before(*opts.Until) {
				panic("startTime must be after Until option")
			}
		} else {
			if !startTime.Before(*opts.Until) {
				panic("startTime must be before Until option")
			}
		}
	}

	// Generate time intervals.
	var times []*time.Time
	if opts.Size != nil {
		times = make([]*time.Time, 0, *opts.Size)
	} else {
		times = []*time.Time{}
	}

	gen, err := TimeIntervalGenerator(timeFreq)
	if err != nil {
		return nil, err
	}

	ntg := gen(startTime, reverse)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if opts.Size != nil && len(times) >= *opts.Size {
			break
		}

		nt := ntg()

		if opts.Until != nil {
			if reverse {
				if !nt.After(*opts.Until) {
					break
				}
			} else {
				if !nt.Before(*opts.Until) {
					break
				}
			}
		}

		times = append(times, &nt)
	}

	st := dataframe.NewSeriesTime(name, nil)
	st.Values = times

	return st, nil
}
