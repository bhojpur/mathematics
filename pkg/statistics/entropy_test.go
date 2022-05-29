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

func ExampleEntropy() {
	d := []float64{1.1, 2.2, 3.3}
	e, _ := stats.Entropy(d)
	fmt.Println(e)
	// Output: 1.0114042647073518
}

func TestEntropy(t *testing.T) {
	for _, c := range []struct {
		in  stats.Float64Data
		out float64
	}{
		{stats.Float64Data{4, 8, 5, 1}, 1.2110440167801229},
		{stats.Float64Data{0.8, 0.01, 0.4}, 0.6791185708986585},
		{stats.Float64Data{0.8, 1.1, 0, 5}, 0.7759393943707658},
	} {
		got, err := stats.Entropy(c.in)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if !veryclose(got, c.out) {
			t.Errorf("Max(%.1f) => %.1f != %.1f", c.in, got, c.out)
		}
	}
	_, err := stats.Entropy([]float64{})
	if err == nil {
		t.Errorf("Empty slice didn't return an error")
	}
}

func BenchmarkEntropySmallFloatSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = stats.Entropy(makeFloatSlice(5))
	}
}

func BenchmarkEntropyLargeFloatSlice(b *testing.B) {
	lf := makeFloatSlice(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = stats.Entropy(lf)
	}
}
