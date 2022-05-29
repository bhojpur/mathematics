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

// Create working sample data to test if the legacy
// functions cause a runtime crash or return an error
func TestLegacy(t *testing.T) {

	// Slice of data
	s := []float64{-10, -10.001, 5, 1.1, 2, 3, 4.20, 5}

	// Slice of coordinates
	d := []stats.Coordinate{
		{1, 2.3},
		{2, 3.3},
		{3, 3.7},
		{4, 4.3},
		{5, 5.3},
	}

	// VarP rename compatibility
	_, err := stats.VarP(s)
	if err != nil {
		t.Errorf("VarP not successfully returning PopulationVariance.")
	}

	// VarS rename compatibility
	_, err = stats.VarS(s)
	if err != nil {
		t.Errorf("VarS not successfully returning SampleVariance.")
	}

	// StdDevP rename compatibility
	_, err = stats.StdDevP(s)
	if err != nil {
		t.Errorf("StdDevP not successfully returning StandardDeviationPopulation.")
	}

	// StdDevS rename compatibility
	_, err = stats.StdDevS(s)
	if err != nil {
		t.Errorf("StdDevS not successfully returning StandardDeviationSample.")
	}

	// LinReg rename compatibility
	_, err = stats.LinReg(d)
	if err != nil {
		t.Errorf("LinReg not successfully returning LinearRegression.")
	}

	// ExpReg rename compatibility
	_, err = stats.ExpReg(d)
	if err != nil {
		t.Errorf("ExpReg not successfully returning ExponentialRegression.")
	}

	// LogReg rename compatibility
	_, err = stats.LogReg(d)
	if err != nil {
		t.Errorf("LogReg not successfully returning LogarithmicRegression.")
	}
}
