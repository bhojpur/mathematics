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

import "testing"

func TestGaussianExpectedValue(t *testing.T) {
	fails := 0
	for i := 0; i < 10; i++ {
		mean := 10.0
		x := NewGaussian(mean, 1.0)
		samplestats := ExpectedValueWithConfidence(x, ZScore95())
		t.Log(samplestats)
		if mean < samplestats.Mean-samplestats.CI {
			t.Logf("True mean below sample mean interval")
			fails += 1
		}
		if mean > samplestats.Mean+samplestats.CI {
			t.Logf("True mean above sample mean interval")
			fails += 1
		}
	}
	if fails > 2 {
		t.Errorf("Got more than one failure in a probabilistic test")
	}
}

func TestGaussianSample(t *testing.T) {
	x := NewNormal(5.0, 2.0)
	m := Materialize(x, 100)
	for _, s := range m.Samples {
		v := s.value
		if v < -3.0 || v > 13.0 {
			t.Error("Gaussian sample way out of range")
		}
	}
}

func TestGaussianMean(t *testing.T) {
	fails := 0
	for i := 0; i < 10; i++ {
		x := NewGaussian(5.0, 1.0)
		m := Materialize(x, 100)
		avg := m.Average()
		// If everything is working, this has about a 0.003% chance of a false positive
		// (99.9997% confidence interval with n=100, sigma=1.0 is +/- 0.4)
		t.Log(avg)
		if avg <= 4.6 || avg >= 5.4 {
			t.Log("Mean outside expected bounds (small chance of error)")
			fails += 1
		}
	}
	if fails > 1 {
		t.Error("Mean repeatedly outside expected bounds")
	}
}
