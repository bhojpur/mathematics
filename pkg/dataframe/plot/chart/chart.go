package chart

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
	"time"

	"github.com/wcharczuk/go-chart"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

// S converts a SeriesFloat64 to a chart.Series for use with the "github.com/wcharczuk/go-chart" package.
// Currently x can be nil, a SeriesFloat64 or a SeriesTime. nil values in the x and y Series are ignored.
//
// NOTE: To "unjoin" the lines, you can adjust the style to chart.Style{StrokeWidth: chart.Disabled, DotWidth: 2}.
func S(ctx context.Context, y *dataframe.SeriesFloat64, x interface{}, r *dataframe.Range, style ...chart.Style) (chart.Series, error) {

	var out chart.Series

	if r == nil {
		r = &dataframe.Range{}
	}

	yNRows := y.NRows(dataframe.DontLock)

	start, end, err := r.Limits(yNRows)
	if err != nil {
		return nil, err
	}

	switch xx := x.(type) {
	case nil:
		cs := chart.ContinuousSeries{Name: y.Name(dataframe.DontLock)}

		if len(style) > 0 {
			cs.Style = style[0]
		}

		xVals := []float64{}
		yVals := []float64{}

		// Remove nil values
		for i, j := 0, start; j < end+1; i, j = i+1, j+1 {

			if err := ctx.Err(); err != nil {
				return nil, err
			}

			yval := y.Values[j]

			if dataframe.IsValidFloat64(yval) {
				yVals = append(yVals, yval)
				xVals = append(xVals, float64(i))
			}
		}

		cs.XValues = xVals
		cs.YValues = yVals

		out = cs
	case *dataframe.SeriesFloat64:

		cs := chart.ContinuousSeries{Name: y.Name(dataframe.DontLock)}

		if len(style) > 0 {
			cs.Style = style[0]
		}

		xVals := []float64{}
		yVals := []float64{}

		// Remove nil values
		for j := start; j < end+1; j++ {

			if err := ctx.Err(); err != nil {
				return nil, err
			}

			yval := y.Values[j]
			xval := xx.Values[j]

			if dataframe.IsValidFloat64(yval) {
				// Check x val is valid
				if dataframe.IsValidFloat64(xval) {
					yVals = append(yVals, yval)
					xVals = append(xVals, xval)
				}
			}
		}

		cs.XValues = xVals
		cs.YValues = yVals

		out = cs
	case *dataframe.SeriesTime:

		cs := chart.TimeSeries{Name: y.Name(dataframe.DontLock)}

		if len(style) > 0 {
			cs.Style = style[0]
		}

		xVals := []time.Time{}
		yVals := []float64{}

		// Remove nil values
		for j := start; j < end+1; j++ {

			if err := ctx.Err(); err != nil {
				return nil, err
			}

			yval := y.Values[j]
			xval := xx.Values[j]

			if dataframe.IsValidFloat64(yval) {
				// Check x val is valid
				if xval != nil {
					yVals = append(yVals, yval)
					xVals = append(xVals, *xval)
				}
			}
		}

		cs.XValues = xVals
		cs.YValues = yVals

		out = cs
	default:
		panic("unrecognized x")
	}

	return out, nil
}
