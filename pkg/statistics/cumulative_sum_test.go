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

func ExampleCumulativeSum() {
	data := []float64{1.0, 2.1, 3.2, 4.823, 4.1, 5.8}
	csum, _ := stats.CumulativeSum(data)
	fmt.Println(csum)
	// Output: [1 3.1 6.300000000000001 11.123000000000001 15.223 21.023]
}

func TestCumulativeSum(t *testing.T) {
	for _, c := range []struct {
		in  []float64
		out []float64
	}{
		{[]float64{1, 2, 3}, []float64{1, 3, 6}},
		{[]float64{1.0, 1.1, 1.2, 2.2}, []float64{1.0, 2.1, 3.3, 5.5}},
		{[]float64{-1, -1, 2, -3}, []float64{-1, -2, 0, -3}},
	} {
		got, err := stats.CumulativeSum(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if !reflect.DeepEqual(c.out, got) {
			t.Errorf("CumulativeSum(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}
	_, err := stats.CumulativeSum([]float64{})
	if err == nil {
		t.Errorf("Empty slice should have returned an error")
	}
}

func BenchmarkCumulativeSumSmallFloatSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = stats.CumulativeSum(makeFloatSlice(5))
	}
}

func BenchmarkCumulativeSumLargeFloatSlice(b *testing.B) {
	lf := makeFloatSlice(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.CumulativeSum(lf)
	}
}
