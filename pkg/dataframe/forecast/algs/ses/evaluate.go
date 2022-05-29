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

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
	"github.com/bhojpur/mathematics/pkg/dataframe/forecast"
)

// Evaluate will measure the quality of the predicted values based on the evaluation calculation defined by evalFunc.
// It will compare the error between sf and the values from the end of the loaded data ("validation set").
// sf is usually the output of the Predict method.
//
// NOTE: You can use the functions directly from the validation subpackage if you need to do something
// other than that described above.
func (se *SimpleExpSmoothing) Evaluate(ctx context.Context, sf *dataframe.SeriesFloat64, evalFunc forecast.EvaluationFunc) (float64, error) {

	if evalFunc == nil {
		panic("evalFunc is nil")
	}

	// Determine outer range of loaded data
	loadedSeries := se.sf
	loadedRows := loadedSeries.NRows(dataframe.DontLock)

	_, te, err := se.tRange.Limits(loadedRows)
	if err != nil {
		return 0, err
	}

	s, e, err := (&dataframe.Range{Start: &[]int{te + 1}[0]}).Limits(loadedRows)
	if err != nil {
		// There is no data in validation set
		return 0, nil
	}

	lengthOfValidationSet := e - s + 1
	lengthOfPredictionSet := sf.NRows(dataframe.DontLock)

	// Pick the smallest range
	var minR int
	if lengthOfValidationSet < lengthOfPredictionSet {
		minR = lengthOfValidationSet
	} else {
		minR = lengthOfPredictionSet
	}

	errVal, _, err := evalFunc(ctx, loadedSeries.Values[s:s+minR], sf.Values[0:minR], nil)
	return errVal, err
}
