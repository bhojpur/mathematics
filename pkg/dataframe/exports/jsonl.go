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
	"encoding/json"
	"io"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

// JSONExportOptions contains options for ExportToJSON function.
type JSONExportOptions struct {

	// NullString is used to set what nil values should be encoded to.
	// Common options are strings: NULL, \N, NaN, NA.
	// If not set, then null (non-string) is used.
	NullString *string

	// Range is used to export a subset of rows from the Dataframe.
	Range dataframe.Range

	// SetEscapeHTML specifies whether problematic HTML characters should be escaped inside JSON quoted strings.
	SetEscapeHTML bool
}

// ExportToJSON exports a Dataframe in the jsonl format.
// Each line represents a row from the Dataframe.
//
// See: http://jsonlines.org/ for more information.
func ExportToJSON(ctx context.Context, w io.Writer, df *dataframe.DataFrame, options ...JSONExportOptions) error {

	df.Lock()
	defer df.Unlock()

	var r dataframe.Range
	var null *string // default is null

	enc := json.NewEncoder(w)

	if len(options) > 0 {

		r = options[0].Range

		enc.SetEscapeHTML(options[0].SetEscapeHTML)

		if options[0].NullString != nil {
			null = options[0].NullString
		}
	}

	nRows := df.NRows(dataframe.DontLock)

	if nRows > 0 {

		s, e, err := r.Limits(nRows)
		if err != nil {
			return err
		}

		for row := s; row <= e; row++ {

			if err := ctx.Err(); err != nil {
				return err
			}

			record := map[string]interface{}{}
			for _, aSeries := range df.Series {

				fieldName := aSeries.Name()

				val := aSeries.Value(row)

				if val == nil && null != nil {
					record[fieldName] = null
				} else {
					record[fieldName] = val
				}
			}

			if err := enc.Encode(record); err != nil {
				return err
			}

		}
	}

	return nil
}
