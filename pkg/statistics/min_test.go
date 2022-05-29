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
	"fmt"
	"testing"

	stats "github.com/bhojpur/mathematics/pkg/statistics"
)

func ExampleMin() {
	d := stats.LoadRawData([]interface{}{1.1, "2", 3.0, 4, "5"})
	a, _ := stats.Min(d)
	fmt.Println(a)
	// Output: 1.1
}

func TestMin(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{1.1, 2, 3, 4, 5}, 1.1},
		{[]float64{10.534, 3, 5, 7, 9}, 3.0},
		{[]float64{-5, 1, 5}, -5.0},
		{[]float64{5}, 5},
	} {
		got, err := stats.Min(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if got != c.out {
			t.Errorf("Min(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}
	_, err := stats.Min([]float64{})
	if err == nil {
		t.Errorf("Empty slice didn't return an error")
	}
}

func BenchmarkMinSmallFloatSlice(b *testing.B) {
	testData := makeFloatSlice(5)
	for i := 0; i < b.N; i++ {
		_, _ = stats.Min(testData)
	}
}

func BenchmarkMinSmallRandFloatSlice(b *testing.B) {
	testData := makeRandFloatSlice(5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Min(testData)
	}
}

func BenchmarkMinLargeFloatSlice(b *testing.B) {
	testData := makeFloatSlice(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Min(testData)
	}
}

func BenchmarkMinLargeRandFloatSlice(b *testing.B) {
	testData := makeRandFloatSlice(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Min(testData)
	}
}
