package hw

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

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
	"github.com/bhojpur/mathematics/pkg/dataframe/forecast"
)

// Predict forecasts the next n values for the loaded data.
func (hw *HoltWinters) Predict(ctx context.Context, n uint) (*dataframe.SeriesFloat64, []forecast.Confidence, error) {

	name := hw.sf.Name(dataframe.DontLock)
	nsf := dataframe.NewSeriesFloat64(name, &dataframe.SeriesInit{Capacity: int(n)})

	if n <= 0 {
		if len(hw.cfg.ConfidenceLevels) == 0 {
			return nsf, nil, nil
		}
		return nsf, []forecast.Confidence{}, nil
	}

	cnfdnce := []forecast.Confidence{}

	var (
		st        float64   = hw.tstate.smoothingLevel
		seasonals []float64 = hw.tstate.seasonalComps
		trnd      float64   = hw.tstate.trendLevel
		period    int       = int(hw.cfg.Period)
	)

	for i := uint(0); i < n; i++ {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}

		m := int(i + 1)

		var fval float64
		if hw.cfg.SeasonalMethod == Multiplicative {
			fval = (st + float64(m)*trnd) * seasonals[(m-1)%period]
		} else {
			fval = (st + float64(m)*trnd) + seasonals[(m-1)%period]
		}
		nsf.Append(fval, dataframe.DontLock)

		cis := map[float64]forecast.ConfidenceInterval{}
		for _, level := range hw.cfg.ConfidenceLevels {
			cis[level] = forecast.DriftConfidenceInterval(fval, level, hw.tstate.rmse, hw.tstate.T, n)
		}
		cnfdnce = append(cnfdnce, cis)
	}

	if len(hw.cfg.ConfidenceLevels) == 0 {
		return nsf, nil, nil
	}
	return nsf, cnfdnce, nil
}
