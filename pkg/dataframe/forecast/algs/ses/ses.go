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

// It implements the simple exponential smooting forecasting algorithm.

import (
	"context"
	"errors"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
	"github.com/bhojpur/mathematics/pkg/dataframe/forecast"
)

// ExponentialSmoothingConfig is used to configure the SES algorithm.
// SES models the error, trend and seasonal elements of the data with exponential smoothing.
//
// NOTE: SES algorithm does not tolerate nil values. You may need to use the interpolation subpackage.
type ExponentialSmoothingConfig struct {

	// Alpha must be between 0 and 1. The closer Alpha is to 1, the more the algorithm
	// prioritizes recent values over past values.
	Alpha float64

	// ConfidenceLevels are values between 0 and 1 (exclusive) that return the associated
	// confidence intervals for each forecasted value.
	ConfidenceLevels []float64
}

// Validate checks if the config is valid.
func (cfg *ExponentialSmoothingConfig) Validate() error {
	if (cfg.Alpha < 0.0) || (cfg.Alpha > 1.0) {
		return errors.New("Alpha must be between [0,1]")
	}

	for _, c := range cfg.ConfidenceLevels {
		if c <= 0.0 || c >= 1.0 {
			return errors.New("ConfidenceLevel value must be between (0,1)")
		}
	}

	return nil
}

// SimpleExpSmoothing represents the SES algorithm for time-series forecasting.
// It uses the bootstrapping method found here: https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc432.htm
type SimpleExpSmoothing struct {
	tstate trainingState
	cfg    ExponentialSmoothingConfig
	tRange dataframe.Range // training range
	sf     *dataframe.SeriesFloat64
}

// NewExponentialSmoothing creates a new SimpleExpSmoothing object.
func NewExponentialSmoothing() *SimpleExpSmoothing {
	return &SimpleExpSmoothing{}
}

// Configure sets the various parameters for the SES algorithm.
// config must be a ExponentialSmoothingConfig.
func (se *SimpleExpSmoothing) Configure(config interface{}) error {

	cfg := config.(ExponentialSmoothingConfig)
	if err := cfg.Validate(); err != nil {
		return err
	}

	se.cfg = cfg
	return nil
}

// Load loads historical data.
// r is used to limit which rows of sf are loaded. Prediction will always begin
// from the row after that defined by r. r can be thought of as defining a "training set".
//
// NOTE: SES algorithm does not tolerate nil values. You may need to use the interpolation subpackage.
func (se *SimpleExpSmoothing) Load(ctx context.Context, sf *dataframe.SeriesFloat64, r *dataframe.Range) error {

	if r == nil {
		r = &dataframe.Range{}
	}

	tLength := sf.NRows(dataframe.DontLock)

	nrows, _ := r.NRows(tLength)
	if nrows == 0 {
		return forecast.ErrInsufficientDataPoints
	}

	s, e, err := r.Limits(tLength)
	if err != nil {
		return err
	}

	// at least 3 observations required for SES.
	if e-s < 2 {
		return forecast.ErrInsufficientDataPoints
	}

	// Check if there are any nil values
	nils, err := sf.NilCount(dataframe.NilCountOptions{
		Ctx:          ctx,
		R:            r,
		StopAtOneNil: true,
		DontLock:     true,
	})
	if err != nil {
		return err
	}
	if nils > 0 {
		return forecast.ErrInsufficientDataPoints
	}

	se.tRange = *r
	se.sf = sf
	se.tstate = trainingState{}

	err = se.trainSeries(ctx, uint(s), uint(e))
	if err != nil {
		se.tRange = dataframe.Range{}
		se.sf = nil
		return err
	}

	return nil
}
