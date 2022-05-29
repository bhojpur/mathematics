package forecast

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
	"errors"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

// ForecastingAlgorithm defines the methods that all forecasting algorithms must implement.
type ForecastingAlgorithm interface {

	// Configure sets the various parameters for the algorithm.
	// config must be a struct that the particular algorithm recognizes.
	Configure(config interface{}) error

	// Load loads historical data.
	// Some forecasting algorithms do not tolerate nil values and will require interpolation.
	// r is used to limit which rows of sf are loaded. Prediction will always begin
	// from the row after that defined by r. r can be thought of as defining a "training set".
	Load(ctx context.Context, sf *dataframe.SeriesFloat64, r *dataframe.Range) error

	// Predict forecasts the next n values for the loaded data.
	//
	// NOTE: Not all forecasting algorithms return the confidence values.
	Predict(ctx context.Context, n uint) (*dataframe.SeriesFloat64, []Confidence, error)

	// Evaluate will measure the quality of the predicted values based on the evaluation calculation defined by evalFunc.
	// It will compare the error between sf and the values from the end of the loaded data ("validation set").
	// sf is usually the output of the Predict method.
	//
	// NOTE: You can use the functions directly from the validation subpackage if you need to do something
	// other than that described above.
	Evaluate(ctx context.Context, sf *dataframe.SeriesFloat64, evalFunc EvaluationFunc) (float64, error)
}

// EvaluationFuncOptions is used to modify the behavior of the EvaluationFunc.
type EvaluationFuncOptions struct {

	// SkipInvalids will skip Inf and NaN values.
	// If set to false (default) and an invalid value is encountered, then an ErrIndeterminate is returned.
	SkipInvalids bool
}

// EvaluationFunc compares the validationSet and forecastSet and calculates the error.
// See the validation subpackage for various approaches to calculating the error.
type EvaluationFunc func(ctx context.Context, validationSet, forecastSet []float64, opts *EvaluationFuncOptions) (float64, int, error)

// ErrInsufficientDataPoints signifies that a particular forecasting algorithm requires more
// data points to operate.
var ErrInsufficientDataPoints = errors.New("insufficient data points or nil values found")

// ErrMismatchLen signifies that there is a mismatch between the length of 2 Series or slices.
var ErrMismatchLen = errors.New("mismatch length")

// ErrIndeterminate indicates that the result of a calculation is indeterminate.
var ErrIndeterminate = errors.New("indeterminate")
