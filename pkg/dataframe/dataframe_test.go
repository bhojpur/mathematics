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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNRows(t *testing.T) {

	s1 := NewSeriesInt64("day", nil, 1, 2, 3)
	s2 := NewSeriesFloat64("sales", nil, 50.3, 23.4, 56.2)
	df := NewDataFrame(s1, s2)

	expected := 3

	if df.NRows() != expected {
		t.Errorf("wrong val: expected: %v actual: %v", expected, df.NRows())
	}

}

func TestInsertAndRemove(t *testing.T) {

	s1 := NewSeriesInt64("day", nil, 1, 2, 3)
	s2 := NewSeriesFloat64("sales", nil, 50.3, 23.4, 56.2)
	df := NewDataFrame(s1, s2)

	df.Append(&dontLock, 9, 123.6)

	df.Append(&dontLock, map[string]interface{}{
		"day":   10,
		"sales": nil,
	})

	df.Remove(0)

	df.Prepend(&dontLock, map[string]interface{}{
		"day":   99,
		"sales": 199.99,
	})

	df.Prepend(&dontLock, 1000, 10000)
	df.UpdateRow(0, &dontLock, 10000, 1000)
	df.Update(0, 1, 9000)

	expected := `+-----+-------+---------+
|     |  DAY  |  SALES  |
+-----+-------+---------+
| 0:  | 10000 |  9000   |
| 1:  |  99   | 199.99  |
| 2:  |   2   |  23.4   |
| 3:  |   3   |  56.2   |
| 4:  |   9   |  123.6  |
| 5:  |  10   |   NaN   |
+-----+-------+---------+
| 6X2 | INT64 | FLOAT64 |
+-----+-------+---------+`

	if strings.TrimSpace(df.Table()) != strings.TrimSpace(expected) {
		t.Errorf("wrong val: expected: %v actual: %v", expected, df.Table())
	}
}

func TestSwap(t *testing.T) {

	s1 := NewSeriesInt64("day", nil, 1, 2, 3)
	s2 := NewSeriesFloat64("sales", nil, 50.3, 23.4, 56.2)
	df := NewDataFrame(s1, s2)

	df.Swap(0, 2)

	expectedValues := [][]interface{}{
		{int64(3), int64(2), int64(1)},
		{56.2, 23.4, 50.3},
	}

	iterator := df.ValuesIterator(ValuesOptions{0, 1, true})
	df.Lock()
	for {
		row, vals, _ := iterator()
		if row == nil {
			break
		}

		for key, val := range vals {
			switch idx := key.(type) {
			case int:
				expected := expectedValues[idx][*row]
				actual := val //df.Series[idx].Value(*row)

				if !cmp.Equal(expected, actual, cmpopts.IgnoreUnexported(SeriesFloat64{}, SeriesInt64{}, SeriesString{}, SeriesTime{}, SeriesGeneric{})) {
					t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
				}
			}
		}
	}
	df.Unlock()
}

func TestNames(t *testing.T) {

	s1 := NewSeriesInt64("day", nil, 1, 2, 3)
	s2 := NewSeriesFloat64("sales", nil, 50.3, 23.4, 56.2)
	df := NewDataFrame(s1, s2)

	// Test names list
	expected := []string{"day", "sales"}

	actual := df.Names()

	if !cmp.Equal(expected, actual, cmpopts.IgnoreUnexported(SeriesFloat64{}, SeriesInt64{}, SeriesString{}, SeriesTime{}, SeriesGeneric{})) {
		t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
	}

	// Test name to column
	input := []string{
		"day",
		"sales",
	}

	actuals := []int{
		0,
		1,
	}

	for i, colName := range input {

		actual, err := df.NameToColumn(colName)
		if err != nil {
			t.Errorf("wrong val: %s err: %v", colName, err)
		} else {
			expected := actuals[i]
			if !cmp.Equal(expected, actual, cmpopts.IgnoreUnexported(SeriesFloat64{}, SeriesInt64{}, SeriesString{}, SeriesTime{}, SeriesGeneric{})) {
				t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
			}
		}
	}

	_, err := df.NameToColumn("unknown")
	if err == nil {
		t.Errorf("there should be an error when name is set to 'unknown'")
	}

}

func TestCopy(t *testing.T) {

	s1 := NewSeriesInt64("day", nil, 1, 2, 3)
	s2 := NewSeriesFloat64("sales", nil, 50.3, 23.4, 56.2)
	df := NewDataFrame(s1, s2)

	cp := df.Copy()

	if !cmp.Equal(df, cp, cmpopts.IgnoreUnexported(DataFrame{}, SeriesFloat64{}, SeriesInt64{}, SeriesString{}, SeriesTime{}, SeriesGeneric{})) {
		t.Errorf("wrong val: expected: %v actual: %v", df, cp)
	}
}

func TestSort(t *testing.T) {

	s1 := NewSeriesInt64("day", nil, nil, 1, 2, 4, 3, nil)
	s2 := NewSeriesFloat64("sales", nil, nil, 50.3, 23.4, 23.4, 56.2, nil)
	df := NewDataFrame(s1, s2)

	sks := []SortKey{
		{Key: "sales", Desc: true},
		{Key: "day", Desc: false},
	}

	df.Sort(context.Background(), sks)

	expectedValues := [][]interface{}{
		{int64(3), int64(1), int64(2), int64(4), nil, nil},
		{56.2, 50.3, 23.4, 23.4, nil, nil},
	}

	iterator := df.ValuesIterator(ValuesOptions{0, 1, true})
	df.Lock()
	for {
		row, vals, _ := iterator()
		if row == nil {
			break
		}

		for key, val := range vals {
			switch colName := key.(type) {
			case string:
				idx, _ := df.NameToColumn(colName, dontLock)

				expected := expectedValues[idx][*row]
				actual := val //df.Series[idx].Value(*row)

				if !cmp.Equal(expected, actual, cmpopts.IgnoreUnexported(SeriesFloat64{}, SeriesInt64{}, SeriesString{}, SeriesTime{}, SeriesGeneric{})) {
					t.Errorf("wrong val: expected: %T %v actual: %T %v", expected, expected, actual, actual)
				}
			}
		}
	}
	df.Unlock()

}

func TestDfIsEqual(t *testing.T) {
	ctx := context.Background()

	s1 := NewSeriesInt64("day", nil, nil, 1, 2, 4, 3, nil)
	s2 := NewSeriesFloat64("sales", nil, nil, 50.3, 23.4, 23.4, 56.2, nil)

	df1 := NewDataFrame(s1, s2)
	df2 := NewDataFrame(s1, s2)

	eq, err := df1.IsEqual(ctx, df2, IsEqualOptions{CheckName: true})
	if err != nil {
		t.Errorf("error encountered: %s\n", err)
	}

	if !eq {
		t.Errorf("Df1: [%T] %s is not equal to Df2: [%T] %s\n", df1, df1.String(), df2, df2.String())
	}
}
