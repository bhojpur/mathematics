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

func TestQuartileOutliers(t *testing.T) {
	s1 := []float64{-1000, 1, 3, 4, 4, 6, 6, 6, 6, 7, 8, 15, 18, 100}
	o, _ := stats.QuartileOutliers(s1)

	if o.Mild[0] != 15 {
		t.Errorf("First Mild Outlier %v != 15", o.Mild[0])
	}

	if o.Mild[1] != 18 {
		t.Errorf("Second Mild Outlier %v != 18", o.Mild[1])
	}

	if o.Extreme[0] != -1000 {
		t.Errorf("First Extreme Outlier %v != -1000", o.Extreme[0])
	}

	if o.Extreme[1] != 100 {
		t.Errorf("Second Extreme Outlier %v != 100", o.Extreme[1])
	}

	_, err := stats.QuartileOutliers([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}
