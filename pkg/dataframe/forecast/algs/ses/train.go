package ses

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
	finalSmoothed *float64 // stores the smoothed value of the final observation point
	yOrigin       float64
	rmse          float64
	T             uint // how many observed values used in the forcasting process
}

func (se *SimpleExpSmoothing) trainSeries(ctx context.Context, start, end uint) error {

	var α float64 = se.cfg.Alpha

	var mse float64

	// Step 1: Calculate Smoothed values for existing observations
	for i, j := start, 1; i < end+1; i, j = i+1, j+1 {
		if err := ctx.Err(); err != nil {
			return err
		}

		if j == 1 {
			// not applicable
		} else if j == 2 {
			se.tstate.finalSmoothed = &se.sf.Values[start]
		} else {
			St := α*se.sf.Values[i-1] + (1-α)**se.tstate.finalSmoothed
			se.tstate.finalSmoothed = &St

			err := se.sf.Values[i] - St // actual value - smoothened value
			mse = mse + err*err
		}
	}
	se.tstate.T = end - start + 1

	// Step 2: Store the y origin
	se.tstate.yOrigin = se.sf.Values[end]

	// Step 3: Calculate rmse
	se.tstate.rmse = math.Sqrt(mse / float64(end-start-1))

	return nil
}
