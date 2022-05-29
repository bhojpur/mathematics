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

// It provides an interface for custom forecasting algorithms.

import (
	"context"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

// Forecast predicts the next n values of sdf using the forecasting algorithm alg.
// cfg is required to configure the parameters of the algorithm. r is used to select a subset of sdf to
// be the "training set". Values after r form the "validation set". evalFunc can be set to measure the
// quality of the predictions. sdf can be a SeriesFloat64 or a DataFrame. DataFrame input is not yet implemented.
//
// NOTE: You can find basic forecasting algorithms in forecast/algs subpackage.
func Forecast(ctx context.Context, sdf interface{}, r *dataframe.Range, alg ForecastingAlgorithm, cfg interface{}, n uint, evalFunc EvaluationFunc) (interface{}, []Confidence, float64, error) {

	switch sdf := sdf.(type) {
	case *dataframe.SeriesFloat64:

		err := alg.Configure(cfg)
		if err != nil {
			return nil, nil, 0, err
		}

		err = alg.Load(ctx, sdf, r)
		if err != nil {
			return nil, nil, 0, err
		}

		pred, cnfdnce, err := alg.Predict(ctx, n)
		if err != nil {
			return nil, nil, 0, err
		}

		var errVal float64
		if evalFunc != nil {
			errVal, err = alg.Evaluate(ctx, pred, evalFunc)
			if err != nil {
				return nil, nil, 0, err
			}
		}

		return pred, cnfdnce, errVal, nil

	case *dataframe.DataFrame:
		panic("sdf as a DataFrame is not yet implemented")
	default:
		panic("sdf must be a Series or DataFrame")
	}

	panic("no reach")
}
