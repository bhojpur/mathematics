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
	"bytes"
	"fmt"

	"github.com/olekukonko/tablewriter"
)

// TableOptions can be used to limit the number of rows and which Series
// are used when generating the table.
type TableOptions struct {

	// Series is used to display a given set of Series. When nil (default), all Series are displayed.
	// An index of the Series or the name of the Series can be provided.
	//
	// NOTE: This option only applies to DataFrames.
	//
	// Example:
	//
	//  opts :=  TableOptions{Series: []interface{}{1, "time"}}
	//
	Series []interface{}

	// R is used to limit the range of rows.
	R *Range

	// DontLock can be set to true if the DataFrame or Series should not be locked.
	DontLock bool
}

// Table will produce the DataFrame in a table.
func (df *DataFrame) Table(opts ...TableOptions) string {

	if len(opts) == 0 || !opts[0].DontLock {
		df.lock.RLock()
		defer df.lock.RUnlock()
	}

	if len(opts) == 0 {
		opts = append(opts, TableOptions{R: &Range{}})
	} else if opts[0].R == nil {
		opts[0].R = &Range{}
	}

	columns := map[interface{}]struct{}{}
	for _, v := range opts[0].Series {
		columns[v] = struct{}{}
	}

	data := [][]string{}

	headers := []string{""} // row header is blank
	footers := []string{fmt.Sprintf("%dx%d", df.n, len(df.Series))}
	for idx, aSeries := range df.Series {
		if len(columns) == 0 {
			headers = append(headers, aSeries.Name())
			footers = append(footers, aSeries.Type())
		} else {
			// Check idx
			_, exists := columns[idx]
			if exists {
				headers = append(headers, aSeries.Name())
				footers = append(footers, aSeries.Type())
				continue
			}

			// Check series name
			_, exists = columns[aSeries.Name()]
			if exists {
				headers = append(headers, aSeries.Name())
				footers = append(footers, aSeries.Type())
				continue
			}
		}
	}

	if df.n > 0 {
		s, e, err := opts[0].R.Limits(df.n)
		if err != nil {
			panic(err)
		}

		for row := s; row <= e; row++ {

			sVals := []string{fmt.Sprintf("%d:", row)}

			for idx, aSeries := range df.Series {
				if len(columns) == 0 {
					sVals = append(sVals, aSeries.ValueString(row))
				} else {
					// Check idx
					_, exists := columns[idx]
					if exists {
						sVals = append(sVals, aSeries.ValueString(row))
						continue
					}

					// Check series name
					_, exists = columns[aSeries.Name()]
					if exists {
						sVals = append(sVals, aSeries.ValueString(row))
						continue
					}
				}
			}

			data = append(data, sVals)
		}
	}

	var buf bytes.Buffer

	table := tablewriter.NewWriter(&buf)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.SetFooter(footers)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Render()

	return buf.String()
}

// String implements the fmt.Stringer interface. It does not lock the DataFrame.
func (df *DataFrame) String() string {

	if df.NRows() <= 6 {
		return df.Table(TableOptions{DontLock: true})
	}

	idx := []int{0, 1, 2, df.n - 3, df.n - 2, df.n - 1}

	data := [][]string{}

	headers := []string{""} // row header is blank
	footers := []string{fmt.Sprintf("%dx%d", df.n, len(df.Series))}
	for _, aSeries := range df.Series {
		headers = append(headers, aSeries.Name())
		footers = append(footers, aSeries.Type())
	}

	for j, row := range idx {

		if j == 3 {
			sVals := []string{"⋮"}

			for range df.Series {
				sVals = append(sVals, "⋮")
			}

			data = append(data, sVals)
		}

		sVals := []string{fmt.Sprintf("%d:", row)}

		for _, aSeries := range df.Series {
			sVals = append(sVals, aSeries.ValueString(row))
		}

		data = append(data, sVals)
	}

	var buf bytes.Buffer

	table := tablewriter.NewWriter(&buf)
	table.SetHeader(headers)
	for _, v := range data {
		table.Append(v)
	}
	table.SetFooter(footers)
	table.SetAlignment(tablewriter.ALIGN_CENTER)

	table.Render()

	return buf.String()
}
