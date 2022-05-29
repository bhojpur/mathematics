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
	"math"
)

type trainingState struct {
	initialSmooth        float64
	initialTrend         float64
	initialSeasonalComps []float64
	smoothingLevel       float64
	trendLevel           float64
	seasonalComps        []float64
	rmse                 float64
	T                    uint // how many observed values used in the forcasting process
}

func (hw *HoltWinters) trainSeries(ctx context.Context, start, end int) error {

	var (
		α, β, γ        float64 = hw.cfg.Alpha, hw.cfg.Beta, hw.cfg.Gamma
		period         int     = int(hw.cfg.Period)
		trnd, prevTrnd float64 // trend
		st, prevSt     float64 // smooth
	)

	y := hw.sf.Values[start : end+1]

	seasonals := initialSeasonalComponents(y, period, hw.cfg.SeasonalMethod)

	hw.tstate.initialSeasonalComps = initialSeasonalComponents(y, period, hw.cfg.SeasonalMethod)

	trnd = initialTrend(y, period)
	hw.tstate.initialTrend = trnd

	var mse float64 // mean squared error

	// Training smoothing Level
	for i := start; i < end+1; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		xt := y[i]

		if i == start { // Set initial smooth
			st = xt
			hw.tstate.initialSmooth = xt
		} else {
			if hw.cfg.SeasonalMethod == Multiplicative {
				// multiplicative method
				prevSt, st = st, α*(xt/seasonals[i%period])+(1-α)*(st+trnd)
				trnd = β*(st-prevSt) + (1-β)*trnd
				seasonals[i%period] = γ*(xt/st) + (1-γ)*seasonals[i%period]
			} else {
				// additive method
				prevSt, st = st, α*(xt-seasonals[i%period])+(1-α)*(st+trnd)
				prevTrnd, trnd = trnd, β*(st-prevSt)+(1-β)*trnd
				seasonals[i%period] = γ*(xt-prevSt-prevTrnd) + (1-γ)*seasonals[i%period]
			}

			err := (xt - seasonals[i%period]) // actual value - smoothened value
			mse = mse + err*err
		}

	}
	hw.tstate.T = uint(end - start + 1)
	hw.tstate.rmse = math.Sqrt(mse / float64(end-start))

	hw.tstate.smoothingLevel = st
	hw.tstate.trendLevel = trnd
	hw.tstate.seasonalComps = seasonals

	return nil
}
