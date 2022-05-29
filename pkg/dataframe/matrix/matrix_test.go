package matrix

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
	"testing"

	dataframe "github.com/bhojpur/mathematics/pkg/dataframe"
)

func TestTranspose(t *testing.T) {

	s1 := dataframe.NewSeriesFloat64("0", nil, 1, 2)
	s2 := dataframe.NewSeriesFloat64("1", nil, 3, 4)
	s3 := dataframe.NewSeriesFloat64("2", nil, 5, 6)
	df := dataframe.NewDataFrame(s1, s2, s3)

	// Transpose df and transpose again to get the same matrix
	mw := MatrixWrap{df}
	nmw := mw.T().T()

	eq, err := mw.DataFrame.IsEqual(context.Background(), nmw.(MatrixWrap).DataFrame)
	if err != nil {
		t.Errorf("wrong err: expected: %v got: %v", nil, err)
	}

	if !eq {
		t.Errorf("matrix transpose error")
	}
}
