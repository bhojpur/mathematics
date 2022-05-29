package statistics

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
)

// Correlation describes the degree of relationship between two sets of data
func Correlation(data1, data2 Float64Data) (float64, error) {

	l1 := data1.Len()
	l2 := data2.Len()

	if l1 == 0 || l2 == 0 {
		return math.NaN(), EmptyInputErr
	}

	if l1 != l2 {
		return math.NaN(), SizeErr
	}

	sdev1, _ := StandardDeviationPopulation(data1)
	sdev2, _ := StandardDeviationPopulation(data2)

	if sdev1 == 0 || sdev2 == 0 {
		return 0, nil
	}

	covp, _ := CovariancePopulation(data1, data2)
	return covp / (sdev1 * sdev2), nil
}

// Pearson calculates the Pearson product-moment correlation coefficient between two variables
func Pearson(data1, data2 Float64Data) (float64, error) {
	return Correlation(data1, data2)
}

// AutoCorrelation is the correlation of a signal with a delayed copy of itself as a function of delay
func AutoCorrelation(data Float64Data, lags int) (float64, error) {
	if len(data) < 1 {
		return 0, EmptyInputErr
	}

	mean, _ := Mean(data)

	var result, q float64

	for i := 0; i < lags; i++ {
		v := (data[0] - mean) * (data[0] - mean)
		for i := 1; i < len(data); i++ {
			delta0 := data[i-1] - mean
			delta1 := data[i] - mean
			q += (delta0*delta1 - q) / float64(i+1)
			v += (delta1*delta1 - v) / float64(i+1)
		}

		result = q / v
	}

	return result, nil
}
