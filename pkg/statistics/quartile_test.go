package statistics_test

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
	"testing"

	stats "github.com/bhojpur/mathematics/pkg/statistics"
)

func TestQuartile(t *testing.T) {
	s1 := []float64{6, 7, 15, 36, 39, 40, 41, 42, 43, 47, 49}
	s2 := []float64{7, 15, 36, 39, 40, 41}

	for _, c := range []struct {
		in []float64
		Q1 float64
		Q2 float64
		Q3 float64
	}{
		{s1, 15, 40, 43},
		{s2, 15, 37.5, 40},
	} {
		quartiles, err := stats.Quartile(c.in)
		if err != nil {
			t.Errorf("Should not have returned an error")
		}

		if quartiles.Q1 != c.Q1 {
			t.Errorf("Q1 %v != %v", quartiles.Q1, c.Q1)
		}
		if quartiles.Q2 != c.Q2 {
			t.Errorf("Q2 %v != %v", quartiles.Q2, c.Q2)
		}
		if quartiles.Q3 != c.Q3 {
			t.Errorf("Q3 %v != %v", quartiles.Q3, c.Q3)
		}
	}

	_, err := stats.Quartile([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}

func TestInterQuartileRange(t *testing.T) {
	s1 := []float64{102, 104, 105, 107, 108, 109, 110, 112, 115, 116, 118}
	iqr, _ := stats.InterQuartileRange(s1)

	if iqr != 10 {
		t.Errorf("IQR %v != 10", iqr)
	}

	_, err := stats.InterQuartileRange([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}

func TestMidhinge(t *testing.T) {
	s1 := []float64{1, 3, 4, 4, 6, 6, 6, 6, 7, 7, 7, 8, 8, 9, 9, 10, 11, 12, 13}
	mh, _ := stats.Midhinge(s1)

	if mh != 7.5 {
		t.Errorf("Midhinge %v != 7.5", mh)
	}

	_, err := stats.Midhinge([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}

func TestTrimean(t *testing.T) {
	s1 := []float64{1, 3, 4, 4, 6, 6, 6, 6, 7, 7, 7, 8, 8, 9, 9, 10, 11, 12, 13}
	tr, _ := stats.Trimean(s1)

	if tr != 7.25 {
		t.Errorf("Trimean %v != 7.25", tr)
	}

	_, err := stats.Trimean([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}
