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

import "math"

// MedianAbsoluteDeviation finds the median of the absolute deviations from the dataset median
func MedianAbsoluteDeviation(input Float64Data) (mad float64, err error) {
	return MedianAbsoluteDeviationPopulation(input)
}

// MedianAbsoluteDeviationPopulation finds the median of the absolute deviations from the population median
func MedianAbsoluteDeviationPopulation(input Float64Data) (mad float64, err error) {
	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	i := copyslice(input)
	m, _ := Median(i)

	for key, value := range i {
		i[key] = math.Abs(value - m)
	}

	return Median(i)
}

// StandardDeviation the amount of variation in the dataset
func StandardDeviation(input Float64Data) (sdev float64, err error) {
	return StandardDeviationPopulation(input)
}

// StandardDeviationPopulation finds the amount of variation from the population
func StandardDeviationPopulation(input Float64Data) (sdev float64, err error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	// Get the population variance
	vp, _ := PopulationVariance(input)

	// Return the population standard deviation
	return math.Sqrt(vp), nil
}

// StandardDeviationSample finds the amount of variation from a sample
func StandardDeviationSample(input Float64Data) (sdev float64, err error) {

	if input.Len() == 0 {
		return math.NaN(), EmptyInputErr
	}

	// Get the sample variance
	vs, _ := SampleVariance(input)

	// Return the sample standard deviation
	return math.Sqrt(vs), nil
}
