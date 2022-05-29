package pandas

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
	"fmt"
	"math"
	"sync"

	"golang.org/x/sync/errgroup"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

func describeDataframe(ctx context.Context, df *dataframe.DataFrame, opts ...DescribeOptions) (DescribeOutput, error) {

	out := DescribeOutput{
		percentiles: opts[0].Percentiles,
	}

	// Compile whitelist and blacklist
	wl := map[int]struct{}{}
	bl := map[int]struct{}{}

	for _, v := range opts[0].Whitelist {
		switch _v := v.(type) {
		case int:
			wl[_v] = struct{}{}
		case string:
			idx, err := df.NameToColumn(_v, dataframe.DontLock)
			if err != nil {
				continue
			}
			wl[idx] = struct{}{}
		default:
			panic(fmt.Errorf("unknown whitelist item: %v", _v))
		}
	}

	for _, v := range opts[0].Blacklist {
		switch _v := v.(type) {
		case int:
			bl[_v] = struct{}{}
		case string:
			idx, err := df.NameToColumn(_v, dataframe.DontLock)
			if err != nil {
				continue
			}
			bl[idx] = struct{}{}
		default:
			panic(fmt.Errorf("unknown blacklist item: %v", _v))
		}
	}

	idxs := []int{}
	g, newCtx := errgroup.WithContext(ctx)
	var lock sync.Mutex
	los := map[int]DescribeOutput{}

	for idx, s := range df.Series {
		idx := idx

		// Check whitelist
		if _, exists := wl[idx]; exists || opts[0].Whitelist == nil {
			// Now check blacklist
			if _, exists := bl[idx]; !exists || opts[0].Blacklist == nil {

				idxs = append(idxs, idx)

				// Accept this Series
				out.headers = append(out.headers, s.Name())

				g.Go(func() error {

					lo, err := describeSeries(newCtx, df.Series[idx], opts[0])
					if err != nil {
						return err
					}

					lock.Lock()
					los[idx] = lo
					lock.Unlock()
					return nil
				})
			}
		}
	}

	err := g.Wait()
	if err != nil {
		return DescribeOutput{}, err
	}

	// Compile results together
	for _, idx := range idxs {
		ldo := los[idx]

		out.Count = append(out.Count, ldo.Count[0])
		out.NilCount = append(out.NilCount, ldo.NilCount[0])

		if len(ldo.Median) > 0 {
			out.Median = append(out.Median, ldo.Median[0])
		} else {
			out.Median = append(out.Median, math.NaN())
		}

		if len(ldo.Mean) > 0 {
			out.Mean = append(out.Mean, ldo.Mean[0])
		} else {
			out.Mean = append(out.Mean, math.NaN())
		}

		if len(ldo.StdDev) > 0 {
			out.StdDev = append(out.StdDev, ldo.StdDev[0])
		} else {
			out.StdDev = append(out.StdDev, math.NaN())
		}

		if len(ldo.Min) > 0 {
			out.Min = append(out.Min, ldo.Min[0])
		} else {
			out.Min = append(out.Min, math.NaN())
		}

		if len(ldo.Max) > 0 {
			out.Max = append(out.Max, ldo.Max[0])
		} else {
			out.Max = append(out.Max, math.NaN())
		}

		if len(ldo.Percentiles) > 0 {
			out.Percentiles = append(out.Percentiles, ldo.Percentiles[0])
		} else {
			out.Percentiles = append(out.Percentiles, []float64{})
		}
	}

	return out, nil
}
