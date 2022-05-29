package exports

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
	"encoding/csv"
	"io"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

// CSVExportOptions contains options for ExportToCSV function.
type CSVExportOptions struct {

	// NullString is used to set what nil values should be encoded to.
	// Common options are NULL, \N, NaN, NA.
	NullString *string

	// Range is used to export a subset of rows from the dataframe.
	Range dataframe.Range

	// Separator is the field delimiter. A common option is ',', which is
	// the default if CSVExportOptions is not provided.
	Separator rune

	// UseCRLF determines the line terminator.
	// When true, it is set to \r\n.
	UseCRLF bool
}

// ExportToCSV exports a Dataframe to a CSV file.
func ExportToCSV(ctx context.Context, w io.Writer, df *dataframe.DataFrame, options ...CSVExportOptions) error {

	df.Lock()
	defer df.Unlock()

	header := []string{}

	var r dataframe.Range

	nullString := "NaN" // Default will be "NaN"

	cw := csv.NewWriter(w)

	if len(options) > 0 {
		cw.Comma = options[0].Separator
		cw.UseCRLF = options[0].UseCRLF
		r = options[0].Range
		if options[0].NullString != nil {
			nullString = *options[0].NullString
		}
	}

	for _, aSeries := range df.Series {
		header = append(header, aSeries.Name())
	}
	if err := cw.Write(header); err != nil {
		return err
	}

	nRows := df.NRows(dataframe.DontLock)

	if nRows > 0 {

		s, e, err := r.Limits(nRows)
		if err != nil {
			return err
		}

		flushCount := 0
		for row := s; row <= e; row++ {

			if err := ctx.Err(); err != nil {
				return err
			}

			flushCount++
			// flush after every 100 writes
			if flushCount > 100 { // flush in the 101th count
				cw.Flush()
				if err := cw.Error(); err != nil {
					return err
				}
				flushCount = 1
			}

			sVals := []string{}
			for _, aSeries := range df.Series {
				val := aSeries.Value(row)
				if val == nil {
					sVals = append(sVals, nullString)
				} else {
					sVals = append(sVals, aSeries.ValueString(row, dataframe.DontLock))
				}
			}

			// Write every row
			if err := cw.Write(sVals); err != nil {
				return err
			}
		}

	}

	// flush before exit
	cw.Flush()
	if err := cw.Error(); err != nil {
		return err
	}

	return nil
}
