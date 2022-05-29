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

func TestSample(t *testing.T) {
	_, err := stats.Sample([]float64{}, 10, false)
	if err == nil {
		t.Errorf("should return an error")
	}

	_, err = stats.Sample([]float64{0.1, 0.2}, 10, false)
	if err == nil {
		t.Errorf("should return an error")
	}
}

func TestSampleWithoutReplacement(t *testing.T) {
	arr := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	result, _ := stats.Sample(arr, 5, false)
	checks := map[float64]bool{}
	for _, res := range result {
		_, ok := checks[res]
		if ok {
			t.Errorf("%v already seen", res)
		}
		checks[res] = true
	}
}

func TestSampleWithReplacement(t *testing.T) {
	arr := []float64{0.1, 0.2, 0.3, 0.4, 0.5}
	numsamples := 100
	result, _ := stats.Sample(arr, numsamples, true)
	if len(result) != numsamples {
		t.Errorf("%v != %v", len(result), numsamples)
	}
}

func TestStableSample(t *testing.T) {
	_, err := stats.StableSample(stats.Float64Data{}, 10)
	if err != stats.EmptyInputErr {
		t.Errorf("should return EmptyInputError when sampling an empty data")
	}
	_, err = stats.StableSample(stats.Float64Data{1.0, 2.0}, 10)
	if err != stats.BoundsErr {
		t.Errorf("should return BoundsErr when sampling size exceeds the maximum element size of data")
	}
	arr := []float64{1.0, 3.0, 2.0, -1.0, 5.0}
	locations := map[float64]int{
		1.0:  0,
		3.0:  1,
		2.0:  2,
		-1.0: 3,
		5.0:  4,
	}
	ret, _ := stats.StableSample(arr, 3)
	if len(ret) != 3 {
		t.Errorf("returned wrong sample size")
	}
	for i := 1; i < 3; i++ {
		if locations[ret[i]] < locations[ret[i-1]] {
			t.Errorf("doesn't keep order")
		}
	}
}
