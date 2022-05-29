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
	"math"
	"testing"

	stats "github.com/bhojpur/mathematics/pkg/statistics"
)

func TestMedianAbsoluteDeviation(t *testing.T) {
	_, err := stats.MedianAbsoluteDeviation([]float64{1, 2, 3})
	if err != nil {
		t.Errorf("Returned an error")
	}
}

func TestMedianAbsoluteDeviationPopulation(t *testing.T) {
	s, _ := stats.MedianAbsoluteDeviation([]float64{1, 2, 3})
	m, err := stats.Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 1.00 {
		t.Errorf("%.10f != %.10f", m, 1.00)
	}

	s, _ = stats.MedianAbsoluteDeviation([]float64{-2, 0, 4, 5, 7})
	m, err = stats.Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 3.00 {
		t.Errorf("%.10f != %.10f", m, 3.00)
	}

	m, _ = stats.MedianAbsoluteDeviation([]float64{})
	if !math.IsNaN(m) {
		t.Errorf("%.1f != %.1f", m, math.NaN())
	}
}

func TestStandardDeviation(t *testing.T) {
	_, err := stats.StandardDeviation([]float64{1, 2, 3})
	if err != nil {
		t.Errorf("Returned an error")
	}
}

func TestStandardDeviationPopulation(t *testing.T) {
	s, _ := stats.StandardDeviationPopulation([]float64{1, 2, 3})
	m, err := stats.Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 0.82 {
		t.Errorf("%.10f != %.10f", m, 0.82)
	}
	s, _ = stats.StandardDeviationPopulation([]float64{-1, -2, -3.3})
	m, err = stats.Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 0.94 {
		t.Errorf("%.10f != %.10f", m, 0.94)
	}

	m, _ = stats.StandardDeviationPopulation([]float64{})
	if !math.IsNaN(m) {
		t.Errorf("%.1f != %.1f", m, math.NaN())
	}
}

func TestStandardDeviationSample(t *testing.T) {
	s, _ := stats.StandardDeviationSample([]float64{1, 2, 3})
	m, err := stats.Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 1.0 {
		t.Errorf("%.10f != %.10f", m, 1.0)
	}
	s, _ = stats.StandardDeviationSample([]float64{-1, -2, -3.3})
	m, err = stats.Round(s, 2)
	if err != nil {
		t.Errorf("Returned an error")
	}
	if m != 1.15 {
		t.Errorf("%.10f != %.10f", m, 1.15)
	}

	m, _ = stats.StandardDeviationSample([]float64{})
	if !math.IsNaN(m) {
		t.Errorf("%.1f != %.1f", m, math.NaN())
	}
}
