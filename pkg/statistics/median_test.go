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
	"reflect"
	"testing"

	stats "github.com/bhojpur/mathematics/pkg/statistics"
)

func ExampleMedian() {
	data := []float64{1.0, 2.1, 3.2, 4.823, 4.1, 5.8}
	median, _ := stats.Median(data)
	fmt.Println(median)
	// Output: 3.65
}

func TestMedian(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out float64
	}{
		{[]float64{5, 3, 4, 2, 1}, 3.0},
		{[]float64{6, 3, 2, 4, 5, 1}, 3.5},
		{[]float64{1}, 1.0},
	} {
		got, _ := stats.Median(c.in)
		if got != c.out {
			t.Errorf("Median(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}
	_, err := stats.Median([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}

func BenchmarkMedianSmallFloatSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = stats.Median(makeFloatSlice(5))
	}
}

func BenchmarkMedianLargeFloatSlice(b *testing.B) {
	lf := makeFloatSlice(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Median(lf)
	}
}

func TestMedianSortSideEffects(t *testing.T) {
	s := []float64{0.1, 0.3, 0.2, 0.4, 0.5}
	a := []float64{0.1, 0.3, 0.2, 0.4, 0.5}
	_, _ = stats.Median(s)
	if !reflect.DeepEqual(s, a) {
		t.Errorf("%.1f != %.1f", s, a)
	}
}
