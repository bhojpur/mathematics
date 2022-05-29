package hw

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

// See: http://www.itl.nist.gov/div898/handbook/pmc/section4/pmc435.htm
func initialTrend(y []float64, period int) float64 {

	var sum float64
	sum = 0.0

	for i := 0; i < period; i++ {
		sum += (y[period+i] - y[i]) / float64(period)
	}

	return sum / float64(period)
}

// See: http://www.itl.nist.gov/div898/handbook/pmc/section4/pmc435.htm
func initialSeasonalComponents(y []float64, period int, tsType Method) []float64 {

	nSeasons := len(y) / period

	seasonalAverage := make([]float64, nSeasons)
	seasonalIndices := make([]float64, period)

	// computing seasonal averages
	for i := 0; i < nSeasons; i++ {
		for j := 0; j < period; j++ {
			seasonalAverage[i] += y[(i*period)+j]
		}
		seasonalAverage[i] /= float64(period)
	}

	// Calculating initial Seasonal component values

	for i := 0; i < period; i++ {
		for j := 0; j < nSeasons; j++ {
			if tsType == Multiplicative {
				// Multiplcative seasonal component
				seasonalIndices[i] += y[(j*period)+i] / seasonalAverage[j]
			} else {
				// Additive seasonal component
				seasonalIndices[i] += y[(j*period)+i] - seasonalAverage[j]
			}

		}
		seasonalIndices[i] /= float64(nSeasons)
	}

	return seasonalIndices
}
