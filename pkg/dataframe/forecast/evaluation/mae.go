package evaluation

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

	"github.com/bhojpur/mathematics/pkg/dataframe/forecast"
)

// MeanAbsoluteError represents the mean absolute error.
var MeanAbsoluteError = func(ctx context.Context, validationSet, forecastSet []float64, opts *forecast.EvaluationFuncOptions) (float64, int, error) {

	// Check if validationSet and forecastSet are the same size
	if len(validationSet) != len(forecastSet) {
		return 0, 0, forecast.ErrMismatchLen
	}

	var (
		n   int
		sum float64
	)

	for i := 0; i < len(validationSet); i++ {

		if err := ctx.Err(); err != nil {
			return 0.0, 0, err
		}

		actual := validationSet[i]
		predicted := forecastSet[i]

		if isInvalidFloat64(actual) || isInvalidFloat64(predicted) {
			if opts != nil && opts.SkipInvalids {
				continue
			} else {
				return 0.0, 0, forecast.ErrIndeterminate
			}
		}

		e := actual - predicted

		sum = sum + math.Abs(e)
		n = n + 1
	}

	if n == 0 {
		return 0.0, 0, forecast.ErrIndeterminate
	}

	return sum / float64(n), n, nil
}
