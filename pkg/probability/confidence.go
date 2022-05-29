package probability

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
	"fmt"
	"math"
)

type MeanAndConfidenceInterval struct {
	Mean float64
	CI   float64
}

func (mci MeanAndConfidenceInterval) String() string {
	return fmt.Sprintf("%0.4f +- %0.4f", mci.Mean, mci.CI)
}

func ExpectedValueWithConfidence(u Uncertain, opts ...Option) MeanAndConfidenceInterval {
	sampleSize := getSampleSize(opts, 1000)
	zScore := getZScore(opts, zScore95)

	m := Materialize(u, sampleSize)

	mean := m.Average()

	squaredError := 0.0
	for _, s := range m.Samples {
		squaredError += math.Pow(s.value-mean, 2.0)
	}
	sdev := math.Sqrt(squaredError / float64(sampleSize-1))

	// Since this is a sample standard deviation we're doing half of a t-test error
	// estimation using the square root of the number of samples to guide the rande
	// of the z scores.
	ci := zScore * sdev / math.Sqrt(float64(sampleSize))
	return MeanAndConfidenceInterval{
		Mean: mean,
		CI:   ci,
	}
}
