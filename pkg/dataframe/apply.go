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
	"context"
)

// ApplyDataFrameFn is used by the Apply function when used with DataFrames.
// vals contains the values for the current row. The keys contain ints (index of Series) and strings (name of Series).
// The returned map must only contain what values you intend to update. The key can be a string (name of Series) or int (index of Series).
// If nil is returned, the existing values for the row are unchanged.
type ApplyDataFrameFn func(vals map[interface{}]interface{}, row, nRows int) map[interface{}]interface{}

// ApplySeriesFn is used by the Apply function when used with Series.
// val contains the value of the current row. The returned value is the updated value.
type ApplySeriesFn func(val interface{}, row, nRows int) interface{}

// Apply will call fn for each row in the Series or DataFrame and replace the existing value with the new
// value returned by fn. When sdf is a DataFrame, fn must be of type ApplyDataFrameFn. When sdf is a Series, fn must be of type ApplySeriesFn.
func Apply(ctx context.Context, sdf interface{}, fn interface{}, opts ...FilterOptions) (interface{}, error) {

	switch typ := sdf.(type) {
	case Series:
		var x ApplySeriesFn

		switch v := fn.(type) {
		case func(val interface{}, row, nRows int) interface{}:
			x = ApplySeriesFn(v)
		default:
			x = fn.(ApplySeriesFn)
		}

		s, err := applySeries(ctx, typ, x, opts...)
		if s == nil {
			return nil, err
		}
		return s, err
	case *DataFrame:
		var x ApplyDataFrameFn

		switch v := fn.(type) {
		case func(vals map[interface{}]interface{}, row, nRows int) map[interface{}]interface{}:
			x = ApplyDataFrameFn(v)
		default:
			x = fn.(ApplyDataFrameFn)
		}

		df, err := applyDataFrame(ctx, typ, x, opts...)
		if df == nil {
			return nil, err
		}
		return df, err
	default:
		panic("sdf must be a Series or DataFrame")
	}

	return nil, nil
}

func applySeries(ctx context.Context, s Series, fn ApplySeriesFn, opts ...FilterOptions) (Series, error) {

	if fn == nil {
		panic("fn is required")
	}

	if len(opts) == 0 {
		opts = append(opts, FilterOptions{})
	}

	if !opts[0].DontLock {
		s.Lock()
		defer s.Unlock()
	}

	nRows := s.NRows(dontLock)

	var ns Series

	if !opts[0].InPlace {
		x, ok := s.(NewSerieser)
		if !ok {
			panic("s must implement NewSerieser interface if InPlace is false")
		}

		// Create a New Series
		ns = x.NewSeries(s.Name(dontLock), &SeriesInit{Capacity: nRows})
	}

	iterator := s.ValuesIterator(ValuesOptions{InitialRow: 0, Step: 1, DontReadLock: true})

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		row, val, nRows := iterator()
		if row == nil {
			break
		}

		newVal := fn(val, *row, nRows)

		if opts[0].InPlace {
			s.Update(*row, newVal, dontLock)
		} else {
			ns.Append(newVal, dontLock)
		}
	}

	if !opts[0].InPlace {
		return ns, nil
	}

	return nil, nil
}

func applyDataFrame(ctx context.Context, df *DataFrame, fn ApplyDataFrameFn, opts ...FilterOptions) (*DataFrame, error) {

	if fn == nil {
		panic("fn is required")
	}

	if len(opts) == 0 {
		opts = append(opts, FilterOptions{})
	}

	if !opts[0].DontLock {
		df.Lock()
		defer df.Unlock()
	}

	nRows := df.n

	var ndf *DataFrame

	if !opts[0].InPlace {

		// Create all series
		seriess := []Series{}
		for i := range df.Series {
			s := df.Series[i]

			x, ok := s.(NewSerieser)
			if !ok {
				panic("all Series in DataFrame must implement NewSerieser interface if InPlace is false")
			}

			seriess = append(seriess, x.NewSeries(df.Series[i].Name(dontLock), &SeriesInit{Capacity: nRows}))
		}

		// Create a new dataframe
		ndf = NewDataFrame(seriess...)
	}

	iterator := df.ValuesIterator(ValuesOptions{InitialRow: 0, Step: 1, DontReadLock: true})

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		row, vals, nRows := iterator()
		if row == nil {
			break
		}

		newVals := fn(vals, *row, nRows)

		if opts[0].InPlace {
			if newVals != nil {
				df.UpdateRow(*row, &dontLock, newVals)
			}
		} else {
			if newVals != nil {
				ndf.Append(&dontLock, newVals)
			} else {
				ndf.Append(&dontLock, vals)
			}
		}
	}

	if !opts[0].InPlace {
		return ndf, nil
	}

	return nil, nil
}
