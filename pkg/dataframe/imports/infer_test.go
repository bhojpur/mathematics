package imports

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
	"strings"
	"testing"
	"time"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

var ctx = context.Background()

func TestCSVImport(t *testing.T) {

	csvStr := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2015-05-07,NA,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United States",2012-02-01,32,321.31,54320
Spain,2012-02-01,66,555.42,00241
`

	opts := CSVLoadOptions{
		InferDataTypes: true,
		NilValue:       &[]string{"NA"}[0],
		DictateDataType: map[string]interface{}{
			"Id": float64(0),
		},
	}

	df, err := LoadFromCSV(ctx, strings.NewReader(csvStr), opts)
	if err != nil {
		t.Errorf("csv import error: %v", err)
		return
	}

	// Expected solution
	parseTime := func(s string) time.Time {
		t, _ := time.Parse("2006-01-02", s)
		return t
	}

	expDf := dataframe.NewDataFrame(
		dataframe.NewSeriesString("country", nil, "United States", "United States", "United Kingdom", "United States", "United Kingdom", "United States", "United States", "Spain"),
		dataframe.NewSeriesTime("date", nil, parseTime("2012-02-01"), parseTime("2012-02-01"), parseTime("2012-02-01"), parseTime("2012-02-01"), parseTime("2015-05-07"), parseTime("2012-02-01"), parseTime("2012-02-01"), parseTime("2012-02-01")),
		dataframe.NewSeriesInt64("age", nil, 50, 32, 17, 32, nil, 32, 32, 66),
		dataframe.NewSeriesFloat64("amount", nil, 112.1, 321.31, 18.2, 321.31, 18.2, 321.31, 321.31, 555.42),
		dataframe.NewSeriesFloat64("id", nil, 1234, 54320, 12345, 54320, 12345, 54320, 54320, 241),
	)

	if eq, _ := df.IsEqual(ctx, expDf); !eq {
		t.Errorf("csv import not equal")
	}
}
