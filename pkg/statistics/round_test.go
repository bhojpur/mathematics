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
	"math"
	"testing"

	stats "github.com/bhojpur/mathematics/pkg/statistics"
)

func ExampleRound() {
	rounded, _ := stats.Round(1.534424, 1)
	fmt.Println(rounded)
	// Output: 1.5
}

func TestRound(t *testing.T) {
	for _, c := range []struct {
		number   float64
		decimals int
		result   float64
	}{
		{0.1111, 1, 0.1},
		{-0.1111, 2, -0.11},
		{5.3253, 3, 5.325},
		{5.3258, 3, 5.326},
		{5.3253, 0, 5.0},
		{5.55, 1, 5.6},
	} {
		m, err := stats.Round(c.number, c.decimals)
		if err != nil {
			t.Errorf("Returned an error")
		}
		if m != c.result {
			t.Errorf("%.1f != %.1f", m, c.result)
		}

	}
	_, err := stats.Round(math.NaN(), 2)
	if err == nil {
		t.Errorf("Round should error on NaN")
	}
}

func BenchmarkRound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = stats.Round(0.1111, 1)
	}
}
