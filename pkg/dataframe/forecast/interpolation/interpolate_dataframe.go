package interpolation

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
	"sync"

	"golang.org/x/sync/errgroup"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

func interpolateDataFrame(ctx context.Context, df *dataframe.DataFrame, opts InterpolateOptions) (map[interface{}]*dataframe.OrderedMapIntFloat64, error) {
	if !opts.DontLock {
		df.Lock()
		defer df.Unlock()
	}

	var lock sync.Mutex
	omaps := map[interface{}]*dataframe.OrderedMapIntFloat64{}

	if opts.HorizAxis != nil {
		switch s := opts.HorizAxis.(type) {
		case int:
			opts.HorizAxis = df.Series[s]
		case string:
			i, err := df.NameToColumn(s, dataframe.DontLock)
			if err != nil {
				return nil, err
			}
			opts.HorizAxis = df.Series[i]
		case dataframe.Series:

		default:
			panic("HorizAxis option must be a SeriesFloat64/SeriesTime or convertable to a SeriesFloat64")
		}
	}

	g, newCtx := errgroup.WithContext(ctx)

	for i := range df.Series {
		i := i
		if df.Series[i] == opts.HorizAxis {
			continue
		}

		fs, ok := df.Series[i].(*dataframe.SeriesFloat64)
		if !ok {
			continue
		}

		g.Go(func() error {
			omap, err := Interpolate(newCtx, fs, opts)
			if err != nil {
				return err
			}

			if !opts.InPlace {
				lock.Lock()
				omaps[i] = omap.(*dataframe.OrderedMapIntFloat64)
				omaps[df.Series[i].Name()] = omap.(*dataframe.OrderedMapIntFloat64)
				lock.Unlock()
			}

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	if opts.InPlace {
		return nil, nil
	}
	return omaps, nil
}
