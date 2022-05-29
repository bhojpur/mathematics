package dataframe

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
)

// ToSeriesInt64 is an interface used by the Dataframe to know if a particular
// Series can be converted to a SeriesInt64 Series.
type ToSeriesInt64 interface {

	// ToSeriesInt64 is used to convert a particular Series to a SeriesInt64.
	// If the returned Series is not nil but an error is still provided,
	// it means that some rows were not able to be converted. You can inspect
	// the error to determine which rows were unconverted.
	//
	// NOTE: The returned ErrorCollection should contain RowError objects.
	ToSeriesInt64(context.Context, bool, ...func(interface{}) (*int64, error)) (*SeriesInt64, error)
}

// ToSeriesString is an interface used by the Dataframe to know if a particular
// Series can be converted to a SeriesString Series.
type ToSeriesString interface {

	// ToSeriesString is used to convert a particular Series to a SeriesString.
	// If the returned Series is not nil but an error is still provided,
	// it means that some rows were not able to be converted. You can inspect
	// the error to determine which rows were unconverted.
	//
	// NOTE: The returned ErrorCollection should contain RowError objects.
	ToSeriesString(context.Context, bool, ...func(interface{}) (*string, error)) (*SeriesString, error)
}

// ToSeriesFloat64 is an interface used by the Dataframe to know if a particular
// Series can be converted to a SeriesFloat64 Series.
type ToSeriesFloat64 interface {

	// ToSeriesFloat64 is used to convert a particular Series to a SeriesFloat64.
	// If the returned Series is not nil but an error is still provided,
	// it means that some rows were not able to be converted. You can inspect
	// the error to determine which rows were unconverted.
	//
	// NOTE: The returned ErrorCollection should contain RowError objects.
	ToSeriesFloat64(context.Context, bool, ...func(interface{}) (float64, error)) (*SeriesFloat64, error)
}

// ToSeriesMixed is an interface used by the Dataframe to know if a particular
// Series can be converted to a ToSeriesMixed Series.
type ToSeriesMixed interface {

	// ToSeriesMixed is used to convert a particular Series to a ToSeriesMixed.
	// If the returned Series is not nil but an error is still provided,
	// it means that some rows were not able to be converted. You can inspect
	// the error to determine which rows were unconverted.
	//
	// NOTE: The returned ErrorCollection should contain RowError objects.
	ToSeriesMixed(context.Context, bool, ...func(interface{}) (interface{}, error)) (*SeriesMixed, error)
}
