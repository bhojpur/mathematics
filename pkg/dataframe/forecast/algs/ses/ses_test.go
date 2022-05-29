package ses_test

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
	"testing"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
	. "github.com/bhojpur/mathematics/pkg/dataframe/forecast/algs/ses"
)

var ctx = context.Background()

func TestSES(t *testing.T) {

	// Test data from https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc431.htm
	data12 := dataframe.NewSeriesFloat64("data", nil, 71, 70, 69, 68, 64, 65, 72, 78, 75, 75, 75, 70)

	alg := NewExponentialSmoothing()
	cfg := ExponentialSmoothingConfig{Alpha: 0.1}

	err := alg.Configure(cfg)
	if err != nil {
		t.Fatalf("configure error: %v", err)
	}

	err = alg.Load(ctx, data12, nil)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	pred, _, err := alg.Predict(ctx, 5)
	if err != nil {
		t.Fatalf("pred error: %v", err)
	}

	// Expected values from https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc432.htm
	expPred := dataframe.NewSeriesFloat64("expected", nil, 71.50, 71.35, 71.21, 71.09, 70.98)

	// compared expPred with pred
	iterator := pred.ValuesIterator(dataframe.ValuesOptions{Step: 1, DontReadLock: true}) // func() (*int, interface{}, int)

	for {
		row, val, _ := iterator()
		if row == nil {
			break
		}

		roundedPred := math.Round(val.(float64)*100) / 100
		if roundedPred != expPred.Values[*row] {
			t.Fatalf("forecasting error. expected = %v, actual = %v (rounded: %v)", expPred.Values[*row], val.(float64), roundedPred)
		}
	}
}
